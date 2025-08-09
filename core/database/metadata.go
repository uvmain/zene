package database

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func createMetadataTable(ctx context.Context) {
	schema := `CREATE TABLE metadata (
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
		label TEXT,
		music_folder_id INTEGER DEFAULT 1,
		FOREIGN KEY (music_folder_id) REFERENCES music_folders(id) ON DELETE CASCADE
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_metadata_track_id", "metadata", "musicbrainz_track_id", false)
	createIndex(ctx, "idx_metadata_album_id", "metadata", "musicbrainz_album_id", false)
	createIndex(ctx, "idx_metadata_artist_id", "metadata", "musicbrainz_artist_id", false)
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

func UpdateMetadataRow(ctx context.Context, metadata types.Metadata) error {
	metadata.DateModified = logic.GetCurrentTimeFormatted()

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

	_, err := DB.ExecContext(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("updating metadata for %s: %v", metadata.FilePath, err)
	}

	logger.Printf("Updated metadata for %s", metadata.FilePath)
	return nil
}

func DeleteMetadataRow(ctx context.Context, filepath string) error {
	query := `DELETE FROM metadata WHERE file_path = ?`
	_, err := DB.ExecContext(ctx, query, filepath)
	if err != nil {
		return fmt.Errorf("deleting metadata row %s: %v", filepath, err)
	}
	logger.Printf("Deleted metadata row %s", filepath)
	return nil
}

func SelectAllFilePathsAndModTimes(ctx context.Context) (map[string]string, error) {
	query := `SELECT file_path, date_modified FROM metadata`
	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying file paths and mod times: %v", err)
	}
	defer rows.Close()

	fileModTimes := make(map[string]string)

	for rows.Next() {
		var filePath, dateModified string
		if err := rows.Scan(&filePath, &dateModified); err != nil {
			return nil, fmt.Errorf("scanning row: %v", err)
		}
		fileModTimes[filePath] = dateModified
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return fileModTimes, nil
}
