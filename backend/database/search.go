package database

import (
	"fmt"
	"zene/types"
)

func SearchMetadata(searchQuery string) ([]types.TrackMetadata, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtSelectFullTextSearchFromMetadata
	stmt.Reset()
	stmt.ClearBindings()

	if searchQuery != "" {
		stmt.SetText("$searchQuery", searchQuery)
	} else {
		return []types.TrackMetadata{}, fmt.Errorf("FTS Query cannot be empty")
	}

	var rows []types.TrackMetadata
	for {
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
