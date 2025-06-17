package database

import (
	"context"
	"fmt"
	"strconv"
	"zene/core/types"

	"zombiezen.com/go/sqlite"
)

func SelectTracksByAlbumID(ctx context.Context, musicbrainz_album_id string) ([]types.Metadata, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.Metadata{}, fmt.Errorf("Failed to take a db conn from the pool in SelectTracksByAlbumID: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT * FROM metadata where musicbrainz_album_id = $musicbrainz_album_id;`)
	defer stmt.Finalize()
	stmt.SetText("$musicbrainz_album_id", musicbrainz_album_id)

	var rows []types.Metadata

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return []types.Metadata{}, err
		} else if !hasRow {
			break
		} else {

			row := types.Metadata{
				FilePath:            stmt.GetText("file_path"),
				DateAdded:           stmt.GetText("date_added"),
				DateModified:        stmt.GetText("date_modified"),
				FileName:            stmt.GetText("file_name"),
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
		rows = []types.Metadata{}
	}
	return rows, nil
}

func SelectAllAlbums(ctx context.Context, random string, limit string, recent string) ([]types.AlbumsResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.AlbumsResponse{}, fmt.Errorf("Failed to take a db conn from the pool in SelectAllAlbums: %v", err)
	}
	defer DbPool.Put(conn)

	var stmt *sqlite.Stmt

	if recent == "true" {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = conn.Prep(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date, date_added FROM metadata group by album ORDER BY date_added desc limit $limit;`)
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = conn.Prep(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date, date_added FROM metadata group by album ORDER BY date_added desc;`)
		}
	} else if random == "true" {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = conn.Prep(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date FROM metadata group by album ORDER BY random() limit $limit;`)
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = conn.Prep(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date FROM metadata group by album ORDER BY random();`)
		}
	} else {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = conn.Prep(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date FROM metadata group by album ORDER BY album limit $limit;`)
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = conn.Prep(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date FROM metadata group by album ORDER BY album;`)
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

func SelectAlbum(ctx context.Context, musicbrainzAlbumId string) (types.AlbumsResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return types.AlbumsResponse{}, fmt.Errorf("Failed to take a db conn from the pool in SelectAlbum: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT album, album_artist, musicbrainz_album_id, musicbrainz_artist_id, genre, release_date FROM metadata where musicbrainz_album_id = $musicbrainz_album_id limit 1;`)
	defer stmt.Finalize()
	stmt.SetText("$musicbrainz_album_id", musicbrainzAlbumId)

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
