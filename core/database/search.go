package database

import (
	"context"
	"database/sql"
	"strings"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func SearchArtists(ctx context.Context, searchQuery string, limit int, offset int, musicFolderId int) ([]types.Artist, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return []types.Artist{}, err
	}

	query := `SELECT m.musicbrainz_artist_id as id,
		m.artist as name,
		count(distinct m.musicbrainz_album_id) as album_count,
		s.created_at as starred,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(AVG(gr.rating),0.0) AS average_rating
	FROM user_music_folders u
	join metadata m on m.music_folder_id = u.folder_id
	LEFT JOIN user_stars s ON m.musicbrainz_artist_id = s.metadata_id AND s.user_id = u.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_artist_id = ur.metadata_id AND ur.user_id = u.user_id
	LEFT JOIN user_ratings gr ON m.musicbrainz_artist_id = gr.metadata_id
	where u.user_id = ?
	and lower(artist) like lower(?)
	group by m.musicbrainz_artist_id
	order by m.musicbrainz_artist_id asc
	limit ?
	offset ?`

	var rows *sql.Rows

	rows, err = DB.QueryContext(ctx, query, user.Id, "%"+searchQuery+"%", limit, offset)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.Artist{}, err
	}
	defer rows.Close()

	var results []types.Artist

	for rows.Next() {
		var result types.Artist
		var starred sql.NullString
		if err := rows.Scan(&result.Id, &result.Name, &result.AlbumCount, &starred, &result.UserRating, &result.AverageRating); err != nil {
			logger.Printf("Failed to scan row in SelectSimilarArtists: %v", err)
			return nil, err
		}

		result.CoverArt = result.Id
		result.ArtistImageUrl = logic.GetUnauthenticatedImageUrl(result.Id)
		if starred.Valid {
			result.Starred = starred.String
		}
		result.SortName = strings.ToLower(result.Name)

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return results, err
	}

	return results, nil
}
