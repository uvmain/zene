package database

import (
	"context"
	"fmt"
	"log"
	"zene/core/types"
)

func SearchMetadata(ctx context.Context, searchQuery string) ([]types.Metadata, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in SearchMetadata: %v", err)
		return []types.Metadata{}, err
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`select distinct m.*
		FROM metadata m JOIN metadata_fts f ON f.file_path = m.file_path
		WHERE metadata_fts MATCH $searchQuery
		ORDER BY m.file_path DESC;`)
	defer stmt.Finalize()

	if searchQuery != "" {
		stmt.SetText("$searchQuery", searchQuery)
	} else {
		return []types.Metadata{}, fmt.Errorf("FTS Query cannot be empty")
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
