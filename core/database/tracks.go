package database

import (
	"context"
	"log"
	"strconv"
	"zene/core/logic"
	"zene/core/types"

	"zombiezen.com/go/sqlite"
)

func SelectAllTracks(ctx context.Context, random string, limit string, recent string) ([]types.TrackMetadata, error) {
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

	if err := logic.CheckContext(ctx); err != nil {
		return []types.TrackMetadata{}, err
	}

	var rows []types.TrackMetadata
	for {
		if err := logic.CheckContext(ctx); err != nil {
			return []types.TrackMetadata{}, err
		}
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

func SelectTrack(ctx context.Context, musicBrainzTrackId string) (types.TrackMetadata, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	if err := logic.CheckContext(ctx); err != nil {
		return types.TrackMetadata{}, err
	}
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

	if err := logic.CheckContext(ctx); err != nil {
		return types.TrackMetadata{}, err
	}

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
