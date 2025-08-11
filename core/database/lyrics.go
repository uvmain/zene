package database

import (
	"context"
	"database/sql"
	"fmt"
	"zene/core/logger"
	"zene/core/types"
)

func createLyricsTable(ctx context.Context) {
	schema := `CREATE TABLE track_lyrics (
		musicbrainz_track_id TEXT PRIMARY KEY,
    plain_lyrics TEXT,
    synced_lyrics TEXT
	);`
	createTable(ctx, schema)
}

func UpsertTrackLyrics(ctx context.Context, musicbrainzTrackId string, lyrics types.LyricsDatabaseRow) error {
	query := `INSERT INTO track_lyrics (musicbrainz_track_id, plain_lyrics, synced_lyrics)
		VALUES (?, ?, ?)
		ON CONFLICT(musicbrainz_track_id) DO UPDATE SET plain_lyrics=excluded.plain_lyrics, synced_lyrics=excluded.synced_lyrics`

	_, err := DB.ExecContext(ctx, query, musicbrainzTrackId, lyrics.PlainLyrics, lyrics.SyncedLyrics)
	if err != nil {
		return fmt.Errorf("upserting track lyrics row: %v", err)
	}
	return nil
}

func GetLyricsForMusicBrainzTrackId(ctx context.Context, musicbrainzTrackId string) (types.LyricsDatabaseRow, error) {
	query := "SELECT musicbrainz_track_id, plain_lyrics, synced_lyrics FROM track_lyrics WHERE musicbrainz_track_id = ?"
	var musicbrainzTrackID, plainLyrics, syncedLyrics string
	err := DB.QueryRowContext(ctx, query, musicbrainzTrackId).Scan(&musicbrainzTrackID, &plainLyrics, &syncedLyrics)
	if err == sql.ErrNoRows {
		logger.Printf("No lyrics found for %s", musicbrainzTrackId)
		return types.LyricsDatabaseRow{}, nil
	} else if err != nil {
		logger.Printf("Error querying lyrics for %s: %v", musicbrainzTrackId, err)
		return types.LyricsDatabaseRow{}, err
	}

	lyrics := types.LyricsDatabaseRow{
		MusicBrainzTrackID: musicbrainzTrackID,
		PlainLyrics:        plainLyrics,
		SyncedLyrics:       syncedLyrics,
	}
	logger.Printf("Retrieved lyrics for %s", musicbrainzTrackId)
	return lyrics, nil
}
