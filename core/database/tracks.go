package database

import (
	"context"
	"fmt"
	"strconv"
	"zene/core/types"

	"zombiezen.com/go/sqlite"
)

func SelectAllTracks(ctx context.Context, random string, limit string, recent string) ([]types.Metadata, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.Metadata{}, fmt.Errorf("Failed to take a db conn from the pool in SelectAllTracks: %v", err)
	}
	defer DbPool.Put(conn)

	var stmt *sqlite.Stmt

	if recent == "true" {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = conn.Prep(`SELECT * FROM metadata ORDER BY date_added desc limit $limit;`)
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = conn.Prep(`SELECT * FROM metadata ORDER BY date_added desc;`)
		}
	} else if random == "true" {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = conn.Prep(`SELECT * FROM metadata order by random() limit $limit;`)
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = conn.Prep(`SELECT * FROM metadata order by random();`)
		}
	} else {
		if limit != "" {
			limitInt, _ := strconv.Atoi(limit)
			stmt = conn.Prep(`SELECT * FROM metadata limit $limit;`)
			stmt.SetInt64("$limit", int64(limitInt))
		} else {
			stmt = conn.Prep(`SELECT * FROM metadata;`)
		}
	}

	defer stmt.Finalize()

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

func SelectTrack(ctx context.Context, musicBrainzTrackId string) (types.Metadata, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return types.Metadata{}, fmt.Errorf("Failed to take a db conn from the pool in t: %v", err)
	}
	defer DbPool.Put(conn)

	var stmt *sqlite.Stmt

	stmt = conn.Prep(`SELECT * FROM metadata where musicbrainz_track_id = $musicbrainz_track_id limit 1;`)
	defer stmt.Finalize()
	stmt.SetText("$musicbrainz_track_id", musicBrainzTrackId)

	var row types.Metadata

	if hasRow, err := stmt.Step(); err != nil {
		return types.Metadata{}, err
	} else if !hasRow {
		return types.Metadata{}, nil
	} else {
		row.FilePath = stmt.GetText("file_path")
		row.DateAdded = stmt.GetText("date_added")
		row.DateModified = stmt.GetText("date_modified")
		row.FileName = stmt.GetText("file_name")
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

func SelectTrackFilesForScanner(ctx context.Context) ([]types.File, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.File{}, fmt.Errorf("Failed to take a db conn from the pool in SelectTrackFilesForScanner: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT file_path, date_modified FROM metadata;`)
	defer stmt.Finalize()

	var rows []types.File
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return []types.File{}, err
		} else if !hasRow {
			break
		}

		row := types.File{
			FilePathAbs:  stmt.GetText("file_path"),
			DateModified: stmt.GetText("date_modified"),
		}
		rows = append(rows, row)
	}

	if rows == nil {
		rows = []types.File{}
	}
	return rows, nil
}
