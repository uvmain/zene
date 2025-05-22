package database

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"zene/types"

	"zombiezen.com/go/sqlite"
)

func createMetadataTable() {
	tableName := "track_metadata"
	schema := `CREATE TABLE IF NOT EXISTS track_metadata (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_id INTEGER,
		filename TEXT,
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
	createTable(tableName, schema)
}

func createMetadataTriggers() {
	createTriggerIfNotExists("track_metadata_after_delete_album_art", `CREATE TRIGGER track_metadata_after_delete_album_art AFTER DELETE ON track_metadata
	BEGIN
			DELETE FROM album_art WHERE musicbrainz_album_id = old.musicbrainz_album_id;
	END;`)
}

func InsertTrackMetadataRow(fileRowId int, metadata types.TrackMetadata) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtInsertTrackMetadataRow
	stmt.Reset()
	stmt.ClearBindings()
	stmt.SetInt64("$file_id", int64(fileRowId))
	stmt.SetText("$filename", metadata.Filename)
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

func DeleteMetadataByFileId(file_id int) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtDeleteMetadataByFileId
	stmt.Reset()
	stmt.ClearBindings()
	stmt.SetInt64("$file_id", int64(file_id))

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to delete metadata row for file_id %d: %v", file_id, err)
	}
	log.Printf("Deleted metadata row for file_id %d", file_id)
	return nil
}

func SelectAllArtists() ([]types.ArtistResponse, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtSelectAllArtists
	stmt.Reset()
	stmt.ClearBindings()

	var rows []types.ArtistResponse

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return []types.ArtistResponse{}, err
		} else if !hasRow {
			break
		} else {
			row := types.ArtistResponse{
				Artist:              stmt.GetText("artist"),
				MusicBrainzArtistID: stmt.GetText("musicbrainz_artist_id"),
			}
			rows = append(rows, row)
		}
	}
	if rows == nil {
		rows = []types.ArtistResponse{}
	}
	return rows, nil
}

func SelectAllAlbums(random string, limit string, recent string) ([]types.AlbumsResponse, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	var stmt *sqlite.Stmt

	if recent == "true" {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = stmtSelectAlbumsRecentlyAddedWithLimit
			stmt.Reset()
			stmt.ClearBindings()
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = stmtSelectAlbumsRecentlyAdded
			stmt.Reset()
			stmt.ClearBindings()
		}
	} else if random == "true" {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = stmtSelectRandomizedAlbumsWithLimit
			stmt.Reset()
			stmt.ClearBindings()
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = stmtSelectRandomizedAlbums
			stmt.Reset()
			stmt.ClearBindings()
		}
	} else {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = stmtSelectAlbumsWithLimit
			stmt.Reset()
			stmt.ClearBindings()
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = stmtSelectAllAlbums
			stmt.Reset()
			stmt.ClearBindings()
		}
	}

	var rows []types.AlbumsResponse
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return []types.AlbumsResponse{}, err
		} else if !hasRow {
			break
		} else {
			row := types.AlbumsResponse{
				Album:               stmt.GetText("album"),
				Artist:              stmt.GetText("album_artist"),
				MusicBrainzAlbumID:  stmt.GetText("musicbrainz_album_id"),
				MusicBrainzArtistID: stmt.GetText("musicbrainz_artist_id"),
				Genres:              stmt.GetText("genre"),
				ReleaseDate:         stmt.GetText("release_date"),
			}
			rows = append(rows, row)
		}
	}
	if rows == nil {
		rows = []types.AlbumsResponse{}
	}
	return rows, nil
}

func SelectAllMetadata(random string, limit string, recent string) ([]types.TrackMetadata, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	var stmt *sqlite.Stmt

	if recent == "true" {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = stmtSelectMetadataRecentlyAddedWithLimit
			stmt.Reset()
			stmt.ClearBindings()
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = stmtSelectMetadataRecentlyAdded
			stmt.Reset()
			stmt.ClearBindings()
		}
	} else if random == "true" {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = stmtSelectMetadataRandomisedWithLimit
			stmt.Reset()
			stmt.ClearBindings()
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = stmtSelectMetadataRandomised
			stmt.Reset()
			stmt.ClearBindings()
		}
	} else {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = stmtSelectMetadataWithLimit
			stmt.Reset()
			stmt.ClearBindings()
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = stmtSelectMetadata
			stmt.Reset()
			stmt.ClearBindings()
		}
	}

	var rows []types.TrackMetadata
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return []types.TrackMetadata{}, err
		} else if !hasRow {
			break
		}

		row := types.TrackMetadata{
			Id:                  int(stmt.GetInt64("id")),
			FileId:              int(stmt.GetInt64("file_id")),
			Filename:            stmt.GetText("filename"),
			Format:              stmt.GetText("format"),
			Duration:            stmt.GetText("duration"),
			Size:                stmt.GetText("size"),
			Bitrate:             stmt.GetText("bitrate"),
			Title:               stmt.GetText("title"),
			Artist:              stmt.GetText("artist"),
			Album:               stmt.GetText("album"),
			AlbumArtist:         stmt.GetText("album_artist"),
			Genre:               stmt.GetText("genre"),
			TrackNumber:         stmt.GetText("track_number"),
			TotalTracks:         stmt.GetText("total_tracks"),
			DiscNumber:          stmt.GetText("disc_number"),
			TotalDiscs:          stmt.GetText("total_discs"),
			ReleaseDate:         stmt.GetText("release_date"),
			MusicBrainzArtistID: stmt.GetText("musicbrainz_artist_id"),
			MusicBrainzAlbumID:  stmt.GetText("musicbrainz_album_id"),
			MusicBrainzTrackID:  stmt.GetText("musicbrainz_track_id"),
			Label:               stmt.GetText("label"),
		}
		rows = append(rows, row)
	}

	if rows == nil {
		rows = []types.TrackMetadata{}
	}
	return rows, nil
}

func SelectMetadataByAlbumID(musicbrainz_album_id string) ([]types.TrackMetadata, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtSelectMetadataByAlbumID
	stmt.Reset()
	stmt.ClearBindings()
	stmt.SetText("$musicbrainz_album_id", musicbrainz_album_id)

	var rows []types.TrackMetadata

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return []types.TrackMetadata{}, err
		} else if !hasRow {
			break
		} else {

			row := types.TrackMetadata{
				Id:                  int(stmt.GetInt64("id")),
				FileId:              int(stmt.GetInt64("file_id")),
				Filename:            stmt.GetText("filename"),
				Format:              stmt.GetText("format"),
				Duration:            stmt.GetText("duration"),
				Size:                stmt.GetText("size"),
				Bitrate:             stmt.GetText("bitrate"),
				Title:               stmt.GetText("title"),
				Artist:              stmt.GetText("artist"),
				Album:               stmt.GetText("album"),
				AlbumArtist:         stmt.GetText("album_artist"),
				Genre:               stmt.GetText("genre"),
				TrackNumber:         stmt.GetText("track_number"),
				TotalTracks:         stmt.GetText("total_tracks"),
				DiscNumber:          stmt.GetText("disc_number"),
				TotalDiscs:          stmt.GetText("total_discs"),
				ReleaseDate:         stmt.GetText("release_date"),
				MusicBrainzArtistID: stmt.GetText("musicbrainz_artist_id"),
				MusicBrainzAlbumID:  stmt.GetText("musicbrainz_album_id"),
				MusicBrainzTrackID:  stmt.GetText("musicbrainz_track_id"),
				Label:               stmt.GetText("label"),
			}
			rows = append(rows, row)
		}
	}
	if rows == nil {
		rows = []types.TrackMetadata{}
	}
	return rows, nil
}

func SelectDistinctGenres(searchParam string) ([]types.GenreResponse, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtSelectDistinctGenres
	stmt.Reset()
	stmt.ClearBindings()

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
