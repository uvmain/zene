package database

import (
	"context"
	"database/sql"
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
	query := `INSERT INTO track_lyrics (musicbrainz_track_id, plain_lyrics, synced_lyrics)
		VALUES (?, ?, ?)
		ON CONFLICT(musicbrainz_track_id) DO UPDATE SET plain_lyrics=excluded.plain_lyrics, synced_lyrics=excluded.synced_lyrics`

	_, err := DB.ExecContext(ctx, query, musicbrainzTrackId, lyrics.PlainLyrics, lyrics.SyncedLyrics)
	if err != nil {
		return fmt.Errorf("upserting track lyrics row: %v", err)
	}
	return nil
}

func GetLyricsForMusicBrainzTrackId(ctx context.Context, musicbrainzTrackId string) (types.Lyrics, error) {
	query := "SELECT plain_lyrics, synced_lyrics FROM track_lyrics WHERE musicbrainz_track_id = ?"
	var plainLyrics, syncedLyrics string
	err := DB.QueryRowContext(ctx, query, musicbrainzTrackId).Scan(&plainLyrics, &syncedLyrics)
	if err == sql.ErrNoRows {
		logger.Printf("No lyrics found for %s", musicbrainzTrackId)
		return types.Lyrics{}, nil
	} else if err != nil {
		logger.Printf("Error querying lyrics for %s: %v", musicbrainzTrackId, err)
		return types.Lyrics{}, err
	}

	lyrics := types.Lyrics{
		PlainLyrics:  plainLyrics,
		SyncedLyrics: syncedLyrics,
	}
	logger.Printf("Retrieved lyrics for %s", musicbrainzTrackId)
	return lyrics, nil
}
