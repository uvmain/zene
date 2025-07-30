package database

import (
	"database/sql"
	"context"
	"fmt"
	"zene/core/types"
)

func SearchMetadata(ctx context.Context, searchQuery string) ([]types.Metadata, error) {
	if searchQuery == "" {
		return []types.Metadata{}, fmt.Errorf("FTS Query cannot be empty")
	}

	query := `SELECT DISTINCT m.*
		FROM metadata m JOIN metadata_fts f ON f.file_path = m.file_path
		WHERE metadata_fts MATCH ?
		ORDER BY m.file_path DESC`
		
	rows, err := DB.QueryContext(ctx, query, searchQuery)
	if err != nil {
		return []types.Metadata{}, fmt.Errorf("searching metadata: %v", err)
	}
	defer rows.Close()

	var result []types.Metadata
	for rows.Next() {
		var row types.Metadata
		err := rows.Scan(
			&row.FilePath, &row.FileName, &row.DateAdded, &row.DateModified,
			&row.Format, &row.Duration, &row.Size, &row.Bitrate,
			&row.Title, &row.Artist, &row.Album, &row.AlbumArtist,
			&row.Genre, &row.TrackNumber, &row.TotalTracks, &row.DiscNumber,
			&row.TotalDiscs, &row.ReleaseDate, &row.MusicBrainzArtistID,
			&row.MusicBrainzAlbumID, &row.MusicBrainzTrackID, &row.Label,
		)
		if err != nil {
			return []types.Metadata{}, fmt.Errorf("scanning metadata row: %v", err)
		}
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return []types.Metadata{}, fmt.Errorf("rows error: %v", err)
	}

	if result == nil {
		result = []types.Metadata{}
	}
	return result, nil
}
