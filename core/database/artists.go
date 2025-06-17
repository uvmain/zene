package database

import (
	"context"
	"fmt"
	"strconv"
	"zene/core/types"

	"zombiezen.com/go/sqlite"
)

func SelectArtistByMusicBrainzArtistId(ctx context.Context, musicbrainzArtistId string) (types.ArtistResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return types.ArtistResponse{}, fmt.Errorf("Failed to take a db conn from the pool in SelectArtistByMusicBrainzArtistId: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT DISTINCT artist, musicbrainz_artist_id FROM metadata where musicbrainz_artist_id = $musicbrainz_artist_id limit 1;`)
	defer stmt.Finalize()
	stmt.SetText("$musicbrainz_artist_id", musicbrainzArtistId)

	if hasRow, err := stmt.Step(); err != nil {
		return types.ArtistResponse{}, err
	} else if !hasRow {
		return types.ArtistResponse{}, nil
	} else {
		var row types.ArtistResponse
		row.Artist = stmt.GetText("artist")
		row.MusicBrainzArtistID = stmt.GetText("musicbrainz_artist_id")
		row.ImageURL = fmt.Sprintf("/api/artists/%s/art", stmt.GetText("musicbrainz_artist_id"))
		return row, nil
	}
}

func SelectAlbumsByArtistId(ctx context.Context, musicbrainz_artist_id string, random string, recent string, chronological string, limit string, offset string) ([]types.AlbumsResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.AlbumsResponse{}, fmt.Errorf("Failed to take a db conn from the pool in SelectTracksByArtistId: %w", err)
	}
	defer DbPool.Put(conn)

	var stmtText string
	var stmt *sqlite.Stmt

	stmtText = "SELECT DISTINCT musicbrainz_album_id, album, musicbrainz_artist_id, artist, genre, release_date FROM metadata WHERE musicbrainz_artist_id = $musicbrainz_artist_id GROUP BY musicbrainz_album_id"

	if recent == "true" {
		stmtText = fmt.Sprintf("%s ORDER BY date_added desc", stmtText)
	} else if random == "true" {
		stmtText = fmt.Sprintf("%s ORDER BY random()", stmtText)
	} else if chronological == "true" {
		stmtText = fmt.Sprintf("%s ORDER BY release_date desc", stmtText)
	}

	if limit != "" {
		stmtText = fmt.Sprintf("%s limit $limit", stmtText)
	}
	if offset != "" {
		stmtText = fmt.Sprintf("%s offset $offset", stmtText)
	}

	stmtText = fmt.Sprintf("%s;", stmtText)

	stmt = conn.Prep(stmtText)
	defer stmt.Finalize()

	stmt.SetText("$musicbrainz_artist_id", musicbrainz_artist_id)

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return []types.AlbumsResponse{}, fmt.Errorf("Failed to convert limit to int: %v", err)
		}
		stmt.SetInt64("$limit", int64(limitInt))
	}
	if offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return []types.AlbumsResponse{}, fmt.Errorf("Failed to convert limit to int: %v", err)
		}
		stmt.SetInt64("$offset", int64(offsetInt))
	}

	var rows []types.AlbumsResponse
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return []types.AlbumsResponse{}, err
		} else if !hasRow {
			break
		}

		row := types.AlbumsResponse{
			MusicBrainzAlbumID:  stmt.GetText("musicbrainz_album_id"),
			Album:               stmt.GetText("album"),
			MusicBrainzArtistID: stmt.GetText("musicbrainz_artist_id"),
			Artist:              stmt.GetText("artist"),
			Genres:              stmt.GetText("genre"),
			ReleaseDate:         stmt.GetText("release_date"),
		}
		rows = append(rows, row)
	}

	if rows == nil {
		rows = []types.AlbumsResponse{}
	}
	return rows, nil
}

func SelectTracksByArtistId(ctx context.Context, musicbrainz_artist_id string, random string, limit string, offset string, recent string) ([]types.Metadata, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.Metadata{}, fmt.Errorf("Failed to take a db conn from the pool in SelectTracksByArtistId: %v", err)
	}
	defer DbPool.Put(conn)

	var stmtText string
	var stmt *sqlite.Stmt

	stmtText = "SELECT * FROM metadata where musicbrainz_artist_id = $musicbrainz_artist_id"

	if recent == "true" {
		stmtText = fmt.Sprintf("%s ORDER BY date_added desc", stmtText)
	} else if random == "true" {
		stmtText = fmt.Sprintf("%s ORDER BY random()", stmtText)
	}
	if limit != "" {
		stmtText = fmt.Sprintf("%s limit $limit", stmtText)
	}
	if offset != "" {
		stmtText = fmt.Sprintf("%s offset $offset", stmtText)
	}

	stmtText = fmt.Sprintf("%s;", stmtText)

	stmt = conn.Prep(stmtText)
	defer stmt.Finalize()

	stmt.SetText("$musicbrainz_artist_id", musicbrainz_artist_id)

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return []types.Metadata{}, fmt.Errorf("Failed to convert limit to int: %v", err)
		}
		stmt.SetInt64("$limit", int64(limitInt))
	}
	if offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return []types.Metadata{}, fmt.Errorf("Failed to convert limit to int: %v", err)
		}
		stmt.SetInt64("$offset", int64(offsetInt))
	}

	var rows []types.Metadata
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return []types.Metadata{}, err
		} else if !hasRow {
			break
		}

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

	if rows == nil {
		rows = []types.Metadata{}
	}
	return rows, nil
}

func SelectAlbumArtists(ctx context.Context, searchParam string, random string, recent string, chronological string, limit string, offset string) ([]types.ArtistResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.ArtistResponse{}, fmt.Errorf("Failed to take a db conn from the pool in SelectAlbumArtists: %v", err)
	}
	defer DbPool.Put(conn)

	var stmtText string
	var stmt *sqlite.Stmt

	stmtText = "select distinct m.album_artist, m.musicbrainz_artist_id FROM metadata m"
	if searchParam != "" {
		stmtText = fmt.Sprintf("%s JOIN artists_fts f ON m.file_path = f.file_path", stmtText)
	}
	stmtText = fmt.Sprintf("%s where m.album_artist = m.artist", stmtText)

	if searchParam != "" {
		stmtText = fmt.Sprintf("%s and artists_fts MATCH $searchQuery", stmtText)
	}

	if recent == "true" {
		stmtText = fmt.Sprintf("%s ORDER BY m.date_added desc", stmtText)
	} else if random == "true" {
		stmtText = fmt.Sprintf("%s ORDER BY random()", stmtText)
	} else if chronological == "true" {
		stmtText = fmt.Sprintf("%s ORDER BY m.release_date desc", stmtText)
	}

	if limit != "" {
		stmtText = fmt.Sprintf("%s limit $limit", stmtText)
	}
	if offset != "" {
		stmtText = fmt.Sprintf("%s offset $offset", stmtText)
	}

	stmtText = fmt.Sprintf("%s;", stmtText)

	stmt = conn.Prep(stmtText)
	defer stmt.Finalize()

	if searchParam != "" {
		stmt.SetText("$searchQuery", searchParam)
	}
	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return []types.ArtistResponse{}, fmt.Errorf("Failed to convert limit to int: %v", err)
		}
		stmt.SetInt64("$limit", int64(limitInt))
	}
	if offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return []types.ArtistResponse{}, fmt.Errorf("Failed to convert limit to int: %v", err)
		}
		stmt.SetInt64("$offset", int64(offsetInt))
	}

	var artists []types.ArtistResponse

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return []types.ArtistResponse{}, err
		} else if !hasRow {
			break
		} else {
			var artist types.ArtistResponse
			artist.Artist = stmt.GetText("album_artist")
			artist.MusicBrainzArtistID = stmt.GetText("musicbrainz_artist_id")
			artist.ImageURL = fmt.Sprintf("/api/artists/%s/art", stmt.GetText("musicbrainz_artist_id"))
			artists = append(artists, artist)
		}
	}
	return artists, nil
}
