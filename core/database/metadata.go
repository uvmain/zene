package database

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
	"zene/core/logger"
	"zene/core/types"
)

func createMetadataTable(ctx context.Context) error {
	tableName := "metadata"
	schema := `CREATE TABLE IF NOT EXISTS metadata (
		file_path TEXT PRIMARY KEY,
		file_name TEXT NOT NULL,
		date_added TEXT NOT NULL,
		date_modified TEXT NOT NULL,
		format TEXT,
		duration TEXT,
		size TEXT,
		bitrate TEXT,
		title TEXT,
		artist TEXT,
		album TEXT,
		album_artist TEXT,
		genre TEXT,
		track_number TEXT,
		total_tracks TEXT,
		disc_number TEXT,
		total_discs TEXT,
		release_date TEXT,
		musicbrainz_artist_id TEXT NOT NULL,
		musicbrainz_album_id TEXT NOT NULL,
		musicbrainz_track_id TEXT NOT NULL,
		label TEXT
	);`
	err := createTable(ctx, tableName, schema)
	if err != nil {
		return err
	}
	createIndex(ctx, "idx_metadata_track_id", "metadata", "musicbrainz_track_id", false)
	createIndex(ctx, "idx_metadata_album_id", "metadata", "musicbrainz_album_id", false)
	createIndex(ctx, "idx_metadata_artist_id", "metadata", "musicbrainz_artist_id", false)
	return nil
}

func InsertMetadataRow(ctx context.Context, metadata types.Metadata) error {
	query := `
	INSERT INTO metadata (
		file_path, date_added, date_modified, file_name, format, duration, size, bitrate, title, artist, album,
		album_artist, genre, track_number, total_tracks, disc_number, total_discs, release_date,
		musicbrainz_artist_id, musicbrainz_album_id, musicbrainz_track_id, label
	) VALUES (
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?
	)
	ON CONFLICT(file_path) DO UPDATE SET
		date_modified = excluded.date_modified,
		file_name = excluded.file_name,
		format = excluded.format,
		duration = excluded.duration,
		size = excluded.size,
		bitrate = excluded.bitrate,
		title = excluded.title,
		artist = excluded.artist,
		album = excluded.album,
		album_artist = excluded.album_artist,
		genre = excluded.genre,
		track_number = excluded.track_number,
		total_tracks = excluded.total_tracks,
		disc_number = excluded.disc_number,
		total_discs = excluded.total_discs,
		release_date = excluded.release_date,
		musicbrainz_artist_id = excluded.musicbrainz_artist_id,
		musicbrainz_album_id = excluded.musicbrainz_album_id,
		musicbrainz_track_id = excluded.musicbrainz_track_id,
		label = excluded.label`

	_, err := DB.ExecContext(ctx, query,
		metadata.FilePath, metadata.DateAdded, metadata.DateModified, metadata.FileName,
		metadata.Format, metadata.Duration, metadata.Size, metadata.Bitrate,
		metadata.Title, metadata.Artist, metadata.Album, metadata.AlbumArtist,
		metadata.Genre, metadata.TrackNumber, metadata.TotalTracks, metadata.DiscNumber,
		metadata.TotalDiscs, metadata.ReleaseDate, metadata.MusicBrainzArtistID,
		metadata.MusicBrainzAlbumID, metadata.MusicBrainzTrackID, metadata.Label)

	if err != nil {
		return fmt.Errorf("inserting metadata row: %v", err)
	}

	return nil
}
}

func UpdateMetadataRow(ctx context.Context, metadata types.Metadata) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("taking a db conn from the pool in UpdateMetadataRow: %v", err)
	}
	defer DbPool.Put(conn)

	metadata.DateModified = time.Now().Format(time.RFC3339Nano)

	v := reflect.ValueOf(metadata)
	t := reflect.TypeOf(metadata)

	var queryParts []string
	var params []interface{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "file_path" || jsonTag == "" {
			continue
		}
		fieldValue := v.Field(i).Interface()
		queryParts = append(queryParts, fmt.Sprintf("%s = ?", jsonTag))
		params = append(params, fieldValue)
	}

	query := fmt.Sprintf("UPDATE metadata SET %s WHERE file_path = ?", strings.Join(queryParts, ", "))
	params = append(params, metadata.FilePath) // primary key goes in the where clause

	stmt, err := conn.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement failed: %w", err)
	}
	defer stmt.Finalize()

	for i, param := range params {
		switch v := param.(type) {
		case int:
			stmt.BindInt64(i+1, int64(v))
		case int64:
			stmt.BindInt64(i+1, v)
		case string:
			stmt.BindText(i+1, v)
		default:
			return fmt.Errorf("unsupported bind type %T at param %d", param, i+1)
		}
	}

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("updating metadata for %s: %w", metadata.FilePath, err)
	}

	logger.Printf("Updated metadata for %s", metadata.FilePath)
	return nil
}

func DeleteMetadataRow(ctx context.Context, filepath string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("taking a db conn from the pool in DeleteMetadataRow: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`delete FROM metadata WHERE file_path = $file_path;`)
	defer stmt.Finalize()
	stmt.SetText("$file_path", filepath)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("deleting metadata row %s: %v", filepath, err)
	}
	logger.Printf("Deleted metadata row %s", filepath)
	return nil
}

func SelectAllFilePathsAndModTimes(ctx context.Context) (map[string]string, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return nil, fmt.Errorf("taking a db conn from the pool in SelectAllFilePathsAndModTimes: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT file_path, date_modified FROM metadata;`)
	defer stmt.Finalize()

	fileModTimes := make(map[string]string)

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return nil, err
		} else if !hasRow {
			break
		} else {
			filePath := stmt.GetText("file_path")
			dateModified := stmt.GetText("date_modified")
			fileModTimes[filePath] = dateModified
		}
	}
	return fileModTimes, nil
}
