package database

import (
	"context"
	"fmt"
	"log"
	"zene/core/logic"
	"zene/core/types"
)

func SearchMetadata(ctx context.Context, searchQuery string) ([]types.TrackMetadata, error) {
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

	stmt := conn.Prep(`select distinct m.file_id, m.filename, m.format, m.duration, m.size, m.bitrate, m.title, m.artist, m.album,
		m.album_artist, m.genre, m.track_number, m.total_tracks, m.disc_number, m.total_discs, m.release_date,
		m.musicbrainz_artist_id, m.musicbrainz_album_id, m.musicbrainz_track_id, m.label
		FROM track_metadata m JOIN track_metadata_fts f ON m.file_id = f.file_id
		WHERE track_metadata_fts MATCH $searchQuery
		ORDER BY m.file_id DESC;`)
	defer stmt.Finalize()

	if searchQuery != "" {
		stmt.SetText("$searchQuery", searchQuery)
	} else {
		return []types.TrackMetadata{}, fmt.Errorf("FTS Query cannot be empty")
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

	if err := logic.CheckContext(ctx); err != nil {
		return []types.TrackMetadata{}, err
	}

	if rows == nil {
		rows = []types.TrackMetadata{}
	}
	return rows, nil
}
