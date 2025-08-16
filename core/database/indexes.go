package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"zene/core/logic"
	"zene/core/types"
)

func GetIndexes(ctx context.Context, userId int64, musicFolderIds []int64, ifModifiedSince int64) (types.SubsonicIndexes, error) {
	latestScan, err := GetLatestCompletedScan(ctx)
	if err != nil {
		return types.SubsonicIndexes{}, err
	}

	latestScanTime := logic.GetStringTimeFormatted(latestScan.CompletedDate)
	latestScanTimeUnix := latestScanTime.UnixMilli()

	response := types.SubsonicIndexes{}
	response.IgnoredArticles = ""
	response.LastModified = latestScanTimeUnix

	if ifModifiedSince != 0 && latestScanTimeUnix <= ifModifiedSince {
		return response, nil
	}

	indexes, err := getArtistIndexes(ctx, userId, musicFolderIds)
	if err != nil {
		return types.SubsonicIndexes{}, err
	}
	response.Indexes = &indexes

	return response, nil
}

func getArtistIndexes(ctx context.Context, userId int64, musicFolderIds []int64) ([]types.Index, error) {

	var rows *sql.Rows
	var err error

	query := `SELECT case when substr(m.artist,1,1) GLOB '*[0-9]*' then '#'
		else upper(substr(m.artist,1,1)) end as artist_index,
		m.musicbrainz_artist_id, m.artist, s.created_at, ur.rating AS user_rating, AVG(gr.rating) AS average_rating
		FROM metadata m
		LEFT JOIN user_stars s ON m.musicbrainz_artist_id = s.metadata_id AND s.user_id = ?
		LEFT JOIN user_ratings ur ON m.musicbrainz_artist_id = ur.metadata_id AND ur.user_id = ?
		LEFT JOIN user_ratings gr ON m.musicbrainz_artist_id = gr.metadata_id`

	args := []any{userId, userId}
	if len(musicFolderIds) > 0 {
		placeholders := make([]string, len(musicFolderIds))
		for i, id := range musicFolderIds {
			placeholders[i] = "?"
			args = append(args, id)
			_ = id
		}
		query += ` WHERE m.music_folder_id IN (` + strings.Join(placeholders, ",") + `)`
	}

	query += ` GROUP BY m.musicbrainz_artist_id, m.artist, s.created_at, ur.rating
		ORDER BY artist_index asc, artist asc;`

	rows, err = DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying artistIndexes: %v", err)
	}
	defer rows.Close()

	var returnIndex []types.Index

	for rows.Next() {
		var artist types.Artist
		var artistIndex string
		var starred sql.NullString
		var userRating sql.NullInt64
		var avgRating sql.NullFloat64
		if err := rows.Scan(&artistIndex, &artist.Id, &artist.Name, &starred, &userRating, &avgRating); err != nil {
			return nil, fmt.Errorf("scanning artist row: %v", err)
		}
		if starred.Valid {
			artist.Starred = starred.String
		}
		if userRating.Valid {
			artist.UserRating = userRating.Int64
		}
		if avgRating.Valid {
			artist.AverageRating = avgRating.Float64
		}

		artist.ArtistImageUrl = logic.GetArtistImageUrl(artist.Id)

		artistEntry := types.Artist{
			Id:             artist.Id,
			Name:           artist.Name,
			ArtistImageUrl: artist.ArtistImageUrl,
			CoverArt:       artist.Id,
			Starred:        artist.Starred,
			UserRating:     artist.UserRating,
			AverageRating:  artist.AverageRating,
		}
		if len(returnIndex) > 0 && returnIndex[len(returnIndex)-1].Name == artistIndex {
			// append to existing index group
			returnIndex[len(returnIndex)-1].Artist = append(returnIndex[len(returnIndex)-1].Artist, artistEntry)
		} else {
			// start new index group
			returnIndex = append(returnIndex, types.Index{
				Name:   artistIndex,
				Artist: []types.Artist{artistEntry},
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating artist rows: %v", err)
	}

	return returnIndex, nil
}
