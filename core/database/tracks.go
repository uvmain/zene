package database

import (
	"context"
	"fmt"
	"strconv"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func SelectAllTracks(ctx context.Context, random string, limit string, offset string, recent string, chronological string) ([]types.MetadataWithPlaycounts, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.MetadataWithPlaycounts{}, fmt.Errorf("taking a db conn from the pool in SelectAllTracks: %v", err)
	}
	defer DbPool.Put(conn)

	userId, _ := logic.GetUserIdFromContext(ctx)
	stmtText := getUnendedMetadataWithPlaycountsSql(userId)

	if recent == "true" {
		stmtText = fmt.Sprintf("%s ORDER BY m.date_added desc", stmtText)
	} else if random != "" {
		randomInteger, err := strconv.Atoi(random)
		if err == nil {
			stmtText += fmt.Sprintf(" ORDER BY ((m.rowid * %d) %% 1000000)", randomInteger)
		} else {
			logger.Printf("Error setting randomness: %v", err)
		}
	} else if chronological == "true" {
		stmtText = fmt.Sprintf("%s ORDER BY m.release_date desc", stmtText)
	}

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return []types.MetadataWithPlaycounts{}, fmt.Errorf("invalid limit value: %v", err)
		}
		stmtText = fmt.Sprintf("%s limit %d", stmtText, limitInt)
	}

	if offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return []types.MetadataWithPlaycounts{}, fmt.Errorf("invalid offset value: %v", err)
		}
		stmtText = fmt.Sprintf("%s offset %d", stmtText, offsetInt)
	}

	stmtText = fmt.Sprintf("%s;", stmtText)
	stmt := conn.Prep(stmtText)

	defer stmt.Finalize()

	var rows []types.MetadataWithPlaycounts
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return []types.MetadataWithPlaycounts{}, err
		} else if !hasRow {
			break
		}

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

	if rows == nil {
		rows = []types.MetadataWithPlaycounts{}
	}
	return rows, nil
}

func SelectTrack(ctx context.Context, musicBrainzTrackId string) (types.MetadataWithPlaycounts, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return types.MetadataWithPlaycounts{}, fmt.Errorf("taking a db conn from the pool in t: %v", err)
	}
	defer DbPool.Put(conn)

	userId, _ := logic.GetUserIdFromContext(ctx)

	stmtText := getUnendedMetadataWithPlaycountsSql(userId)
	stmtText = fmt.Sprintf("%s where m.musicbrainz_track_id = $musicbrainz_track_id limit 1;", stmtText)
	stmt := conn.Prep(stmtText)
	defer stmt.Finalize()
	stmt.SetText("$musicbrainz_track_id", musicBrainzTrackId)

	var row types.MetadataWithPlaycounts

	if hasRow, err := stmt.Step(); err != nil {
		return types.MetadataWithPlaycounts{}, err
	} else if !hasRow {
		return types.MetadataWithPlaycounts{}, nil
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
		row.UserPlayCount = stmt.GetInt64("user_play_count")
		row.GlobalPlayCount = stmt.GetInt64("global_play_count")
	}

	return row, nil
}

func SelectTrackFilesForScanner(ctx context.Context) ([]types.File, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.File{}, fmt.Errorf("taking a db conn from the pool in SelectTrackFilesForScanner: %v", err)
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
