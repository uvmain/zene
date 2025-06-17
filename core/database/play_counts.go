package database

import (
	"context"
	"fmt"
	"time"
	"zene/core/types"
)

func createPlayCountsTable(ctx context.Context) {
	tableName := "play_counts"
	schema := `CREATE TABLE IF NOT EXISTS play_counts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		musicbrainz_track_id TEXT NOT NULL,
		play_count INTEGER NOT NULL DEFAULT 0,
		last_played TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (musicbrainz_track_id) REFERENCES metadata(musicbrainz_track_id) ON DELETE CASCADE,
		UNIQUE (user_id, musicbrainz_track_id)
	);`
	createTable(ctx, tableName, schema)
	createIndex(ctx, "idx_play_counts_track_id", "play_counts", "musicbrainz_track_id", false)
	createIndex(ctx, "idx_play_counts_user_id", "play_counts", "user_id", false)
}

func UpsertPlaycount(ctx context.Context, userId int64, musicbrainzTrackId string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("Failed to take a db conn from the pool in UpsertPlaycount: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`INSERT INTO play_counts (user_id, musicbrainz_track_id, play_count, last_played)
		VALUES ($user_id, $musicbrainz_track_id, 1, $last_played)
		ON CONFLICT(user_id, musicbrainz_track_id)
		DO UPDATE SET play_count = play_count + 1, last_played = excluded.last_played;`)
	defer stmt.Finalize()

	stmt.SetInt64("$user_id", userId)
	stmt.SetText("$musicbrainz_track_id", musicbrainzTrackId)
	stmt.SetText("$last_played", time.Now().Format(time.RFC3339Nano))

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("Failed to upsert playcount: %v", err)
	}
	return nil
}

func GetPlaycounts(ctx context.Context, musicbrainzTrackId string, userId int64) ([]types.Playcount, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.Playcount{}, fmt.Errorf("Failed to take a db conn from the pool in GetPlaycountsForUserId: %v", err)
	}
	defer DbPool.Put(conn)

	var stmtText string

	stmtText = "SELECT id, user_id, musicbrainz_track_id, play_count, last_played FROM play_counts"
	if musicbrainzTrackId != "" {
		stmtText = fmt.Sprintf("%s where musicbrainz_track_id = $musicbrainz_track_id", stmtText)
		if userId != 0 {
			stmtText = fmt.Sprintf("%s and user_id = $user_id", stmtText)
		}
	} else if userId != 0 {
		stmtText = fmt.Sprintf("%s where user_id = $user_id", stmtText)
	}
	stmtText = fmt.Sprintf("%s order by play_count desc, last_played desc", stmtText)

	stmt := conn.Prep(stmtText)
	defer stmt.Finalize()

	if musicbrainzTrackId != "" {
		stmt.SetText("$musicbrainz_track_id", musicbrainzTrackId)
	}
	if userId != 0 {
		stmt.SetInt64("$user_id", userId)
	}

	var rows []types.Playcount
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return []types.Playcount{}, err
		} else if !hasRow {
			break
		}

		row := types.Playcount{
			Id:                 stmt.GetInt64("id"),
			UserId:             stmt.GetInt64("user_id"),
			MusicBrainzTrackID: stmt.GetText("musicbrainz_track_id"),
			PlayCount:          stmt.GetInt64("play_count"),
			LastPlayed:         stmt.GetText("last_played"),
		}
		rows = append(rows, row)
	}

	if rows == nil {
		rows = []types.Playcount{}
	}
	return rows, nil
}
