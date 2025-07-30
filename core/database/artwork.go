package database

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"time"
	"zene/core/types"
)

func createAlbumArtTable(ctx context.Context) error {
	tableName := "album_art"
	schema := `CREATE TABLE IF NOT EXISTS album_art (
		musicbrainz_album_id TEXT PRIMARY KEY,
		date_modified TEXT NOT NULL
	);`
	err := createTable(ctx, tableName, schema)
	return err
}

func createArtistArtTable(ctx context.Context) error {
	tableName := "artist_art"
	schema := `CREATE TABLE IF NOT EXISTS artist_art (
		musicbrainz_artist_id TEXT PRIMARY KEY,
		date_modified TEXT NOT NULL
	);`
	err := createTable(ctx, tableName, schema)
	return err
}

func SelectAlbumArtByMusicBrainzAlbumId(ctx context.Context, musicbrainzAlbumId string) (types.AlbumArtRow, error) {
	query := `SELECT musicbrainz_album_id, date_modified FROM album_art WHERE musicbrainz_album_id = ?`
	var row types.AlbumArtRow
	err := DB.QueryRowContext(ctx, query, musicbrainzAlbumId).Scan(&row.MusicbrainzAlbumId, &row.DateModified)
	if err == sql.ErrNoRows {
		return types.AlbumArtRow{}, nil
	} else if err != nil {
		return types.AlbumArtRow{}, err
	}
	return row, nil
}

func InsertAlbumArtRow(ctx context.Context, musicbrainzAlbumId string, dateModified string) error {
	query := `INSERT INTO album_art (musicbrainz_album_id, date_modified)
		VALUES (?, ?)
		ON CONFLICT(musicbrainz_album_id) DO UPDATE SET date_modified=excluded.date_modified
		WHERE excluded.date_modified>album_art.date_modified`

	_, err := DB.ExecContext(ctx, query, musicbrainzAlbumId, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("inserting album art row: %v", err)
	}
	return nil
}

func SelectArtistSubDirectories(ctx context.Context, musicbrainzArtistId string) ([]string, error) {
	query := `SELECT DISTINCT file_path FROM metadata WHERE musicbrainz_artist_id = ?`
	rows, err := DB.QueryContext(ctx, query, musicbrainzArtistId)
	if err != nil {
		return nil, fmt.Errorf("querying artist subdirectories: %v", err)
	}
	defer rows.Close()

	uniqueDirectories := make(map[string]struct{})

	for rows.Next() {
		var filePath string
		if err := rows.Scan(&filePath); err != nil {
			return nil, fmt.Errorf("scanning file path: %v", err)
		}
		directory := filepath.Dir(filePath)
		uniqueDirectories[directory] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	var result []string
	for directory := range uniqueDirectories {
		result = append(result, directory)
	}

	return result, nil
}

func SelectArtistArtByMusicBrainzArtistId(ctx context.Context, musicbrainzArtistId string) (types.ArtistArtRow, error) {
	query := `SELECT musicbrainz_artist_id, date_modified FROM artist_art WHERE musicbrainz_artist_id = ?`
	var row types.ArtistArtRow
	err := DB.QueryRowContext(ctx, query, musicbrainzArtistId).Scan(&row.MusicbrainzArtistId, &row.DateModified)
	if err == sql.ErrNoRows {
		return types.ArtistArtRow{}, nil
	} else if err != nil {
		return types.ArtistArtRow{}, err
	}
	return row, nil
}

func InsertArtistArtRow(ctx context.Context, musicbrainzArtistId string, dateModified string) error {
	query := `INSERT INTO artist_art (musicbrainz_artist_id, date_modified)
	VALUES (?, ?)
	ON CONFLICT(musicbrainz_artist_id) DO UPDATE SET date_modified=excluded.date_modified
	WHERE excluded.date_modified>artist_art.date_modified`

	_, err := DB.ExecContext(ctx, query, musicbrainzArtistId, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("inserting artist art row: %v", err)
	}
	return nil
}
