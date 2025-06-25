package database

import (
	"context"
	"fmt"
	"strconv"
	"zene/core/logic"
	"zene/core/types"
)

func SelectTracksByAlbumId(ctx context.Context, musicbrainz_album_id string) ([]types.MetadataWithPlaycounts, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.MetadataWithPlaycounts{}, fmt.Errorf("Failed to take a db conn from the pool in SelectTracksByAlbumId: %v", err)
	}
	defer DbPool.Put(conn)

	userId, _ := logic.GetUserIdFromContext(ctx)
	stmtText := getUnendedMetadataWithPlaycountsSql(userId)

	stmtText = fmt.Sprintf("%s where musicbrainz_album_id = $musicbrainz_album_id order by cast(disc_number AS INTEGER), cast(track_number AS INTEGER);", stmtText)

	stmt := conn.Prep(stmtText)
	defer stmt.Finalize()
	stmt.SetText("$musicbrainz_album_id", musicbrainz_album_id)

	var rows []types.MetadataWithPlaycounts

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return []types.MetadataWithPlaycounts{}, err
		} else if !hasRow {
			break
		} else {
			row := types.MetadataWithPlaycounts{
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
				UserPlayCount:       stmt.GetInt64("user_play_count"),
				GlobalPlayCount:     stmt.GetInt64("global_play_count"),
			}
			rows = append(rows, row)
		}
	}
	if rows == nil {
		rows = []types.MetadataWithPlaycounts{}
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

	stmtText := "SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date, date_added FROM metadata group by album"

	if recent == "true" {
		stmtText += " ORDER BY date_added desc"
	} else if random != "" {
		integer, err := strconv.Atoi(random)
		if err == nil {
			stmtText += fmt.Sprintf(" ORDER BY ((rowid * %d) %% 1000000)", integer)
		}
	} else {
		stmtText += " ORDER BY album"
	}

	if limit != "" {
		integer, err := strconv.Atoi(random)
		if err == nil {
			stmtText += fmt.Sprintf(" limit %d", integer)
		}
	}

	stmtText += ";"
	stmt := conn.Prep(stmtText)

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
