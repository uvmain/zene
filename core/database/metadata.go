package database

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"zene/core/types"
)

func createMetadataTable(ctx context.Context) {
	tableName := "track_metadata"
	schema := `CREATE TABLE IF NOT EXISTS track_metadata (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		dir_path TEXT NOT NULL,
		file_path TEXT NOT NULL UNIQUE,
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
	createIndex(ctx, "idx_metadata_track_id", "track_metadata", "musicbrainz_track_id", false)
	createIndex(ctx, "idx_metadata_album_id", "track_metadata", "musicbrainz_album_id", false)
	createIndex(ctx, "idx_metadata_artist_id", "track_metadata", "musicbrainz_artist_id", false)
}

func InsertTrackMetadataRow(ctx context.Context, fileRowId int, metadata types.TrackMetadata) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in InsertTrackMetadataRow: %v", err)
		return err
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`INSERT INTO track_metadata (
		id, dir_path, file_path, date_added, date_modified, file_name, format, duration, size, bitrate, title, artist, album,
		album_artist, genre, track_number, total_tracks, disc_number, total_discs, release_date,
		musicbrainz_artist_id, musicbrainz_album_id, musicbrainz_track_id, label
	) VALUES (
	  $id, $dir_path, $file_path, $date_added, $date_modified, $file_name, $format, $duration, $size, $bitrate, $title, $artist, $album,
		$album_artist, $genre, $track_number, $total_tracks, $disc_number, $total_discs, $release_date,
		$musicbrainz_artist_id, $musicbrainz_album_id, $musicbrainz_track_id, $label
	 )`)
	defer stmt.Finalize()
	stmt.SetInt64("$id", int64(metadata.Id))
	stmt.SetText("$dir_path", metadata.DirPath)
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
		return fmt.Errorf("failed to insert metadata row: %v", err)
	}

	return nil
}

func DeleteMetadataByFileId(ctx context.Context, file_id int) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in DeleteMetadataByFileId: %v", err)
		return err
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`delete FROM track_metadata WHERE file_id = $file_id;`)
	defer stmt.Finalize()
	stmt.SetInt64("$file_id", int64(file_id))

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to delete metadata row for file_id %d: %v", file_id, err)
	}
	log.Printf("Deleted metadata row for file_id %d", file_id)
	return nil
}

func SelectDistinctGenres(ctx context.Context, searchParam string) ([]types.GenreResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in SelectDistinctGenres: %v", err)
		return []types.GenreResponse{}, err
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT DISTINCT genre FROM track_metadata;`)
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
