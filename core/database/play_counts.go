package database

import (
	"context"
	"fmt"
	"zene/core/logic"
	"zene/core/types"
)

func createPlayCountsTable(ctx context.Context) {
	schema := `CREATE TABLE play_counts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		musicbrainz_track_id TEXT NOT NULL,
		play_count INTEGER NOT NULL DEFAULT 0,
		last_played TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE (user_id, musicbrainz_track_id)
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_playcounts_user_track ", "play_counts", []string{"user_id", "musicbrainz_track_id"}, true)
	createIndex(ctx, "idx_playcounts_track ", "play_counts", []string{"musicbrainz_track_id"}, false)
	createIndex(ctx, "idx_play_counts_user", "play_counts", []string{"user_id"}, false)
}

func UpsertPlayCount(ctx context.Context, userId int, musicbrainzTrackId string) error {
	query := `INSERT INTO play_counts (user_id, musicbrainz_track_id, play_count, last_played)
		VALUES (?, ?, 1, ?)
		ON CONFLICT(user_id, musicbrainz_track_id)
		DO UPDATE SET play_count = play_count + 1, last_played = excluded.last_played`

	_, err := DB.ExecContext(ctx, query, userId, musicbrainzTrackId, logic.GetCurrentTimeFormatted())
	if err != nil {
		return fmt.Errorf("upserting playcount: %v", err)
	}
	return nil
}

func GetPlayCounts(ctx context.Context, musicbrainzTrackId string, userId int) ([]types.Playcount, error) {
	var query string
	var args []interface{}

	query = "SELECT id, user_id, musicbrainz_track_id, play_count, last_played FROM play_counts"

	var conditions []string
	if musicbrainzTrackId != "" {
		conditions = append(conditions, "musicbrainz_track_id = ?")
		args = append(args, musicbrainzTrackId)
	}
	if userId != 0 {
		conditions = append(conditions, "user_id = ?")
		args = append(args, userId)
	}

	if len(conditions) > 0 {
		query += " WHERE " + fmt.Sprintf("%s", conditions[0])
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	query += " ORDER BY play_count DESC, last_played DESC"

	rows, err := DB.QueryContext(ctx, query, args...)
	if err != nil {
		return []types.Playcount{}, fmt.Errorf("querying play counts: %v", err)
	}
	defer rows.Close()

	var result []types.Playcount
	for rows.Next() {
		var row types.Playcount
		err := rows.Scan(&row.Id, &row.UserId, &row.MusicBrainzTrackID, &row.PlayCount, &row.LastPlayed)
		if err != nil {
			return []types.Playcount{}, fmt.Errorf("scanning play count row: %v", err)
		}
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return []types.Playcount{}, fmt.Errorf("rows error: %v", err)
	}

	if result == nil {
		result = []types.Playcount{}
	}
	return result, nil
}
