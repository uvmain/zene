package database

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"zene/core/logic"
	"zene/core/types"
)

func migrateArt(ctx context.Context) {
	createTable(ctx, `CREATE TABLE album_art (musicbrainz_album_id TEXT PRIMARY KEY, date_modified TEXT NOT NULL);`)
	createTable(ctx, `CREATE TABLE artist_art (musicbrainz_artist_id TEXT PRIMARY KEY, date_modified TEXT NOT NULL);`)
}

func SelectAlbumArtByMusicBrainzAlbumId(ctx context.Context, musicbrainzAlbumId string) (types.AlbumArtRow, error) {
	query := `SELECT aa.date_modified, m.file_path
		FROM metadata m
		left join album_art aa on aa.musicbrainz_album_id = m.musicbrainz_album_id
		WHERE m.musicbrainz_album_id = ?
		limit 1;`

	var row types.AlbumArtRow
	var dateModified sql.NullString
	err := DbRead.QueryRowContext(ctx, query, musicbrainzAlbumId).Scan(&dateModified, &row.FilePath)
	if err != nil {
		return types.AlbumArtRow{}, err
	}
	if dateModified.Valid {
		row.DateModified = dateModified.String
	}
	return row, nil
}

func UpsertAlbumArtRow(ctx context.Context, musicbrainzAlbumId string) error {
	query := `INSERT INTO album_art (musicbrainz_album_id, date_modified)
		VALUES (?, ?)
		ON CONFLICT(musicbrainz_album_id) DO UPDATE SET date_modified=excluded.date_modified
		WHERE excluded.date_modified>album_art.date_modified`

	_, err := DbWrite.ExecContext(ctx, query, musicbrainzAlbumId, logic.GetCurrentTimeFormatted())
	if err != nil {
		return fmt.Errorf("inserting album art row: %v", err)
	}
	return nil
}

func DeleteAlbumArtRow(ctx context.Context, musicbrainzAlbumId string) error {
	query := `DELETE FROM album_art WHERE musicbrainz_album_id = ?`
	_, err := DbWrite.ExecContext(ctx, query, musicbrainzAlbumId)
	if err != nil {
		return fmt.Errorf("deleting album art row: %v", err)
	}
	return nil
}

func SelectAlbumArtIds(ctx context.Context) ([]string, error) {
	query := `SELECT distinct musicbrainz_album_id FROM album_art`
	rows, err := DbRead.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying album art IDs: %v", err)
	}
	defer rows.Close()

	var albumArtIds []string
	for rows.Next() {
		var albumArtId string
		if err := rows.Scan(&albumArtId); err != nil {
			return nil, fmt.Errorf("scanning album art ID: %v", err)
		}
		albumArtIds = append(albumArtIds, albumArtId)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return albumArtIds, nil
}

func SelectArtistSubDirectories(ctx context.Context, musicbrainzArtistId string) ([]string, error) {
	query := `SELECT DISTINCT file_path FROM metadata WHERE musicbrainz_artist_id = ? and album_artist = artist`
	rows, err := DbRead.QueryContext(ctx, query, musicbrainzArtistId)
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
	err := DbRead.QueryRowContext(ctx, query, musicbrainzArtistId).Scan(&row.MusicbrainzArtistId, &row.DateModified)
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

	_, err := DbWrite.ExecContext(ctx, query, musicbrainzArtistId, logic.GetCurrentTimeFormatted())
	if err != nil {
		return fmt.Errorf("inserting artist art row: %v", err)
	}
	return nil
}
