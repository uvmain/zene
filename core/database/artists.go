package database

import (
	"context"
	"database/sql"
	"strings"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func SelectArtistByMusicBrainzArtistId(ctx context.Context, userId int, musicbrainzArtistId string) (types.Artist, error) {
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
	and m.musicbrainz_artist_id = ?
	group by m.artist`

	var result types.Artist

	var starred sql.NullString

	err := DB.QueryRowContext(ctx, query, userId, musicbrainzArtistId).Scan(
		&result.Id, &result.Name, &result.AlbumCount, &starred, &result.UserRating, &result.AverageRating,
	)
	if err == sql.ErrNoRows {
		return types.Artist{}, nil
	} else if err != nil {
		return types.Artist{}, err
	}

	result.CoverArt = result.Id
	result.ArtistImageUrl = logic.GetUnauthenticatedImageUrl(result.Id)
	result.MusicBrainzId = result.Id
	result.SortName = strings.ToLower(result.Name)

	logger.Printf("user rating: %v", result.UserRating)
	logger.Printf("average rating: %v", result.AverageRating)

	if starred.Valid {
		result.Starred = starred.String
	}

	albums, err := GetArtistChildren(ctx, musicbrainzArtistId)

	result.Album = []types.SubsonicChild{}
	for _, album := range albums {
		result.Album = append(result.Album, album)
	}

	return result, nil
}
