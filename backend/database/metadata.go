package database

import (
	"context"
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

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`INSERT INTO track_metadata (
		file_id, filename, format, duration, size, bitrate, title, artist, album,
		album_artist, genre, track_number, total_tracks, disc_number, total_discs, release_date,
		musicbrainz_artist_id, musicbrainz_album_id, musicbrainz_track_id, label
	) VALUES (
	  $file_id, $filename, $format, $duration, $size, $bitrate, $title, $artist, $album,
		$album_artist, $genre, $track_number, $total_tracks, $disc_number, $total_discs, $release_date,
		$musicbrainz_artist_id, $musicbrainz_album_id, $musicbrainz_track_id, $label
	 )`)
	defer stmt.Finalize()
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

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
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

func SelectArtistByMusicBrainzArtistId(musicbrainzArtistId string) (types.ArtistResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT DISTINCT artist, musicbrainz_artist_id FROM track_metadata	where musicbrainz_artist_id = $musicbrainz_artist_id limit 1;`)
	defer stmt.Finalize()
	stmt.SetText("$musicbrainz_artist_id", musicbrainzArtistId)

	var row types.ArtistResponse

	if hasRow, err := stmt.Step(); err != nil {
		return types.ArtistResponse{}, err
	} else if !hasRow {
		return types.ArtistResponse{}, nil
	} else {
		var row types.ArtistResponse
		row.Artist = stmt.GetText("artist")
		row.MusicBrainzArtistID = stmt.GetText("musicbrainz_artist_id")
		row.ImageURL = fmt.Sprintf("/api/artists/%s/art", stmt.GetText("musicbrainz_artist_id"))
	}

	return row, nil
}

func SelectAllArtists() ([]types.ArtistResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT DISTINCT artist, musicbrainz_artist_id FROM track_metadata	ORDER BY artist;`)
	defer stmt.Finalize()

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
				ImageURL:            fmt.Sprintf("/api/artists/%s/art", stmt.GetText("musicbrainz_artist_id")),
			}
			rows = append(rows, row)
		}
	}
	if rows == nil {
		rows = []types.ArtistResponse{}
	}
	return rows, nil
}

func SelectAllAlbumArtists() ([]types.ArtistResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT DISTINCT album_artist, musicbrainz_artist_id FROM track_metadata where artist = album_artist ORDER BY artist;`)
	defer stmt.Finalize()

	var rows []types.ArtistResponse

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return []types.ArtistResponse{}, err
		} else if !hasRow {
			break
		} else {
			row := types.ArtistResponse{
				Artist:              stmt.GetText("album_artist"),
				MusicBrainzArtistID: stmt.GetText("musicbrainz_artist_id"),
				ImageURL:            fmt.Sprintf("/api/artists/%s/art", stmt.GetText("musicbrainz_artist_id")),
			}
			rows = append(rows, row)
		}
	}
	if rows == nil {
		rows = []types.ArtistResponse{}
	}
	return rows, nil
}

func SelectAlbum(musicbrainzAlbumId string) (types.AlbumsResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT DISTINCT album, album_artist, musicbrainz_album_id, musicbrainz_artist_id, genre, release_date FROM track_metadata limit 1;`)
	defer stmt.Finalize()

	var row types.AlbumsResponse

	if hasRow, err := stmt.Step(); err != nil {
		return types.AlbumsResponse{}, err
	} else if !hasRow {
		return types.AlbumsResponse{}, nil
	} else {
		var row types.AlbumsResponse
		row.Album = stmt.GetText("album")
		row.Artist = stmt.GetText("album_artist")
		row.MusicBrainzAlbumID = stmt.GetText("musicbrainz_album_id")
		row.MusicBrainzArtistID = stmt.GetText("musicbrainz_artist_id")
		row.Genres = stmt.GetText("genre")
		row.ReleaseDate = stmt.GetText("release_date")
	}

	return row, nil
}

func SelectAllAlbums(random string, limit string, recent string) ([]types.AlbumsResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	var stmt *sqlite.Stmt

	if recent == "true" {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = conn.Prep(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date, date_added FROM track_metadata m join files f on m.file_id = f.id group by album ORDER BY f.date_added desc limit $limit;`)
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = conn.Prep(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date, date_added FROM track_metadata m join files f on m.file_id = f.id group by album ORDER BY f.date_added desc;`)
		}
	} else if random == "true" {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = conn.Prep(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date FROM track_metadata group by album ORDER BY random() limit $limit;`)
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = conn.Prep(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date FROM track_metadata group by album ORDER BY random();`)
		}
	} else {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = conn.Prep(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date FROM track_metadata group by album ORDER BY album limit $limit;`)
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = conn.Prep(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date FROM track_metadata group by album ORDER BY album;`)
		}
	}

	defer stmt.Finalize()

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

func SelectAllTracks(random string, limit string, recent string) ([]types.TrackMetadata, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	var stmt *sqlite.Stmt

	if recent == "true" {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = conn.Prep(`SELECT * FROM track_metadata m join files f on m.file_id = f.id ORDER BY f.date_added desc limit $limit;`)
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = conn.Prep(`SELECT * FROM track_metadata m join files f on m.file_id = f.id ORDER BY f.date_added desc;`)
		}
	} else if random == "true" {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = conn.Prep(`SELECT * FROM track_metadata order by random() limit $limit;`)
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = conn.Prep(`SELECT * FROM track_metadata order by random();`)
		}
	} else {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = conn.Prep(`SELECT * FROM track_metadata ORDER BY id limit $limit;`)
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = conn.Prep(`SELECT * FROM track_metadata ORDER BY id;`)
		}
	}

	defer stmt.Finalize()

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

func SelectTrack(musicBrainzTrackId string) (types.TrackMetadata, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	var stmt *sqlite.Stmt

	stmt = conn.Prep(`SELECT * FROM track_metadata where musicbrainz_track_id = $musicbrainz_track_id limit 1;`)
	defer stmt.Finalize()
	stmt.SetText("$musicbrainz_track_id", musicBrainzTrackId)

	var row types.TrackMetadata

	if hasRow, err := stmt.Step(); err != nil {
		return types.TrackMetadata{}, err
	} else if !hasRow {
		return types.TrackMetadata{}, nil
	} else {
		row.Id = int(stmt.GetInt64("id"))
		row.FileId = int(stmt.GetInt64("file_id"))
		row.Filename = stmt.GetText("filename")
		row.Format = stmt.GetText("format")
		row.Duration = stmt.GetText("duration")
		row.Size = stmt.GetText("size")
		row.Bitrate = stmt.GetText("bitrate")
		row.Title = stmt.GetText("title")
		row.Artist = stmt.GetText("artist")
		row.Album = stmt.GetText("album")
		row.AlbumArtist = stmt.GetText("album_artist")
		row.Genre = stmt.GetText("genre")
		row.TrackNumber = stmt.GetText("track_number")
		row.TotalTracks = stmt.GetText("total_tracks")
		row.DiscNumber = stmt.GetText("disc_number")
		row.TotalDiscs = stmt.GetText("total_discs")
		row.ReleaseDate = stmt.GetText("release_date")
		row.MusicBrainzArtistID = stmt.GetText("musicbrainz_artist_id")
		row.MusicBrainzAlbumID = stmt.GetText("musicbrainz_album_id")
		row.MusicBrainzTrackID = stmt.GetText("musicbrainz_track_id")
		row.Label = stmt.GetText("label")
	}

	return row, nil
}

func SelectMetadataByAlbumID(musicbrainz_album_id string) ([]types.TrackMetadata, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT * FROM track_metadata where musicbrainz_album_id = $musicbrainz_album_id ORDER BY id;`)
	defer stmt.Finalize()
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
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
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

func SearchForArtists(searchParam string) ([]types.ArtistResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`select distinct m.artist, m.musicbrainz_artist_id
		FROM track_metadata m JOIN artists_fts f ON m.file_id = f.file_id
		WHERE artists_fts MATCH $searchQuery
		ORDER BY m.file_id DESC;`)
	defer stmt.Finalize()
	stmt.SetText("$searchQuery", searchParam)

	var artists []types.ArtistResponse

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return []types.ArtistResponse{}, err
		} else if !hasRow {
			break
		} else {
			var artist types.ArtistResponse
			artist.Artist = stmt.GetText("artist")
			artist.MusicBrainzArtistID = stmt.GetText("musicbrainz_artist_id")
			artist.ImageURL = fmt.Sprintf("/api/artists/%s/art", stmt.GetText("musicbrainz_artist_id"))
			artists = append(artists, artist)
		}
	}
	return artists, nil
}
