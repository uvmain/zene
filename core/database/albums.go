package database

import (
	"context"
	"log"
	"strconv"
	"zene/core/logic"
	"zene/core/types"

	"zombiezen.com/go/sqlite"
)

func SelectTracksByAlbumID(ctx context.Context, musicbrainz_album_id string) ([]types.TrackMetadata, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	if err := logic.CheckContext(ctx); err != nil {
		return []types.TrackMetadata{}, err
	}

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
		if err := logic.CheckContext(ctx); err != nil {
			return []types.TrackMetadata{}, err
		}

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

func SelectAllAlbums(ctx context.Context, random string, limit string, recent string) ([]types.AlbumsResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	if err := logic.CheckContext(ctx); err != nil {
		return []types.AlbumsResponse{}, err
	}

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

	if err := logic.CheckContext(ctx); err != nil {
		return []types.AlbumsResponse{}, err
	}

	defer stmt.Finalize()

	var rows []types.AlbumsResponse
	for {
		if err := logic.CheckContext(ctx); err != nil {
			return []types.AlbumsResponse{}, err
		}
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

func SelectAlbum(ctx context.Context, musicbrainzAlbumId string) (types.AlbumsResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	if err := logic.CheckContext(ctx); err != nil {
		return types.AlbumsResponse{}, err
	}

	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT album, album_artist, musicbrainz_album_id, musicbrainz_artist_id, genre, release_date FROM track_metadata where musicbrainz_album_id = $musicbrainzAlbumId limit 1;`)
	defer stmt.Finalize()
	stmt.SetText("$musicbrainzAlbumId", musicbrainzAlbumId)

	if err := logic.CheckContext(ctx); err != nil {
		return types.AlbumsResponse{}, err
	}

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
		return row, nil
	}
}
