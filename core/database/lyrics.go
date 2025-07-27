package database

import (
	"context"
	"fmt"
	"zene/core/logger"
	"zene/core/types"
)

func createLyricsTable(ctx context.Context) error {
	tableName := "track_lyrics"
	schema := `CREATE TABLE IF NOT EXISTS track_lyrics (
		musicbrainz_track_id TEXT PRIMARY KEY,
    plain_lyrics TEXT,
    synced_lyrics TEXT
	);`

	err := createTable(ctx, tableName, schema)
	if err != nil {
		return err
	}

	return nil
}

func UpsertTrackLyrics(ctx context.Context, musicbrainzTrackId string, lyrics types.Lyrics) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("taking a db conn from the pool in UpsertTrackLyrics: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`INSERT INTO track_lyrics (musicbrainz_track_id, plain_lyrics, synced_lyrics)
		VALUES ($musicbrainz_track_id, $plain_lyrics, $synced_lyrics)
		ON CONFLICT(musicbrainz_track_id) DO UPDATE SET plain_lyrics=excluded.plain_lyrics, synced_lyrics=excluded.synced_lyrics
	`)
	defer stmt.Finalize()
	stmt.SetText("$musicbrainz_track_id", musicbrainzTrackId)
	stmt.SetText("$plain_lyrics", lyrics.PlainLyrics)
	stmt.SetText("$synced_lyrics", lyrics.SyncedLyrics)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("upserting track lyrics row: %v", err)
	}
	return nil
}

func GetLyricsForMusicBrainzTrackId(ctx context.Context, musicbrainzTrackId string) (types.Lyrics, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return types.Lyrics{}, err
	}
	defer DbPool.Put(conn)
	stmt := conn.Prep("SELECT plain_lyrics, synced_lyrics FROM track_lyrics WHERE musicbrainz_track_id = $musicbrainz_track_id")
	defer stmt.Finalize()
	stmt.SetText("$musicbrainz_track_id", musicbrainzTrackId)
	hasRow, err := stmt.Step()
	if err != nil {
		logger.Printf("Error querying lyrics for %s: %v", musicbrainzTrackId, err)
		return types.Lyrics{}, err
	}
	if !hasRow {
		logger.Printf("No lyrics found for %s", musicbrainzTrackId)
		return types.Lyrics{}, nil
	} else {
		plainLyrics := stmt.GetText("plain_lyrics")
		syncedLyrics := stmt.GetText("synced_lyrics")
		lyrics := types.Lyrics{
			PlainLyrics:  plainLyrics,
			SyncedLyrics: syncedLyrics,
		}
		logger.Printf("Retrieved lyrics for %s", musicbrainzTrackId)
		return lyrics, nil
	}
}
