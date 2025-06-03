package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"zene/core/types"

	"zombiezen.com/go/sqlite"
)

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

func SelectAlbumArtists(searchParam string, random string, limit string, recent string) ([]types.ArtistResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	var stmtText string
	var stmt *sqlite.Stmt

	stmtText = "select distinct m.album_artist, m.musicbrainz_artist_id FROM track_metadata m"
	if searchParam != "" {
		stmtText = fmt.Sprintf("%s JOIN artists_fts s ON m.file_id = s.file_id", stmtText)
	}
	stmtText = fmt.Sprintf("%s join files f on f.id = m.file_id where m.album_artist = m.artist", stmtText)

	if searchParam != "" {
		stmtText = fmt.Sprintf("%s and artists_fts MATCH $searchQuery", stmtText)
	}
	if recent == "true" {
		stmtText = fmt.Sprintf("%s ORDER BY f.date_added desc", stmtText)
	} else if random == "true" {
		stmtText = fmt.Sprintf("%s ORDER BY random()", stmtText)
	}
	if limit != "" {
		stmtText = fmt.Sprintf("%s limit $limit", stmtText)
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
			return []types.ArtistResponse{}, fmt.Errorf("failed to convert limit to int: %v", err)
		}
		stmt.SetInt64("$limit", int64(limitInt))
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
