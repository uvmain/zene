package database

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
	"zene/core/logger"
	"zene/core/types"
)

func createMetadataTable(ctx context.Context) {
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
		musicbrainz_artist_id TEXT,
		musicbrainz_album_id TEXT,
		musicbrainz_track_id TEXT,
		label TEXT
	);`
	createTable(ctx, tableName, schema)
	createIndex(ctx, "idx_metadata_track_id", "metadata", "musicbrainz_track_id", false)
	createIndex(ctx, "idx_metadata_album_id", "metadata", "musicbrainz_album_id", false)
	createIndex(ctx, "idx_metadata_artist_id", "metadata", "musicbrainz_artist_id", false)
}

func InsertMetadataRow(ctx context.Context, metadata types.Metadata) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("Failed to take a db conn from the pool in InsertMetadataRow: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`
	INSERT INTO metadata (
		file_path, date_added, date_modified, file_name, format, duration, size, bitrate, title, artist, album,
		album_artist, genre, track_number, total_tracks, disc_number, total_discs, release_date,
		musicbrainz_artist_id, musicbrainz_album_id, musicbrainz_track_id, label
	) VALUES (
		$file_path, $date_added, $date_modified, $file_name, $format, $duration, $size, $bitrate, $title, $artist, $album,
		$album_artist, $genre, $track_number, $total_tracks, $disc_number, $total_discs, $release_date,
		$musicbrainz_artist_id, $musicbrainz_album_id, $musicbrainz_track_id, $label
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
		label = excluded.label;`)

	defer stmt.Finalize()
	stmt.SetText("$file_path", metadata.FilePath)
	stmt.SetText("$date_added", metadata.DateAdded)
	stmt.SetText("$date_modified", metadata.DateModified)
	stmt.SetText("$file_name", metadata.FileName)
	stmt.SetText("$format", metadata.Format)
	stmt.SetText("$duration", metadata.Duration)
	stmt.SetText("$size", metadata.Size)
	stmt.SetText("$bitrate", metadata.Bitrate)
	stmt.SetText("$title", metadata.Title)
	stmt.SetText("$artist", metadata.Artist)
	stmt.SetText("$album", metadata.Album)
	stmt.SetText("$album_artist", metadata.AlbumArtist)
	stmt.SetText("$genre", metadata.Genre)
	stmt.SetText("$track_number", metadata.TrackNumber)
	stmt.SetText("$total_tracks", metadata.TotalTracks)
	stmt.SetText("$disc_number", metadata.DiscNumber)
	stmt.SetText("$total_discs", metadata.TotalDiscs)
	stmt.SetText("$release_date", metadata.ReleaseDate)
	stmt.SetText("$musicbrainz_artist_id", metadata.MusicBrainzArtistID)
	stmt.SetText("$musicbrainz_album_id", metadata.MusicBrainzAlbumID)
	stmt.SetText("$musicbrainz_track_id", metadata.MusicBrainzTrackID)
	stmt.SetText("$label", metadata.Label)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("Failed to insert metadata row: %v", err)
	}

	return nil
}

func UpdateMetadataRow(ctx context.Context, metadata types.Metadata) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("Failed to take a db conn from the pool in UpdateMetadataRow: %v", err)
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
		return fmt.Errorf("Failed to update metadata for %s: %w", metadata.FilePath, err)
	}

	logger.Printf("Updated metadata for %s", metadata.FilePath)
	return nil
}

func DeleteMetadataRow(ctx context.Context, filepath string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("Failed to take a db conn from the pool in DeleteMetadataRow: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`delete FROM metadata WHERE file_path = $file_path;`)
	defer stmt.Finalize()
	stmt.SetText("$file_path", filepath)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("Failed to delete metadata row %s: %v", filepath, err)
	}
	logger.Printf("Deleted metadata row %s", filepath)
	return nil
}

func SelectDistinctGenres(ctx context.Context, searchParam string) ([]types.GenreResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.GenreResponse{}, fmt.Errorf("Failed to take a db conn from the pool in SelectDistinctGenres: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT DISTINCT genre FROM metadata;`)
	defer stmt.Finalize()

	var genres []string

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return []types.GenreResponse{}, err
		} else if !hasRow {
			break
		} else {
			row := stmt.GetText("genre")
			splits := strings.Split(row, ";")
			for _, split := range splits {
				trimmed := strings.TrimSpace(split)
				if trimmed != "" {
					if searchParam != "" {
						if strings.Contains(strings.ToLower(trimmed), strings.ToLower(searchParam)) {
							genres = append(genres, trimmed)
						}
					} else {
						genres = append(genres, trimmed)
					}
				}
			}
		}
	}

	dict := map[string]int{}
	for _, num := range genres {
		dict[num]++
	}

	var ss []types.GenreResponse
	for k, v := range dict {
		ss = append(ss, types.GenreResponse{
			Genre: k,
			Count: v,
		})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Count > ss[j].Count
	})

	return ss, nil
}

func SelectAllFilePathsAndModTimes(ctx context.Context) (map[string]string, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to take a db conn from the pool in SelectAllFilePathsAndModTimes: %v", err)
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
