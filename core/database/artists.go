package database

import (
	"context"
	"database/sql"
	"strings"
	"zene/core/logic"
	"zene/core/types"
)

func SelectArtistByMusicBrainzArtistId(ctx context.Context, musicbrainzArtistId string) (types.Artist, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return types.Artist{}, err
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
	and m.musicbrainz_artist_id = ?
	group by m.artist`

	var result types.Artist

	var starred sql.NullString

	err = DB.QueryRowContext(ctx, query, user.Id, musicbrainzArtistId).Scan(
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

func GetArtistNameByMusicBrainzArtistId(ctx context.Context, musicBrainzArtistId string) string {
	var artistName string
	query := `SELECT artist FROM metadata WHERE musicbrainz_artist_id = ? limit 1`
	err := DB.QueryRowContext(ctx, query, musicBrainzArtistId).Scan(&artistName)
	if err != nil {
		return ""
	}
	return artistName
}

func GetArtistIdByName(ctx context.Context, artistName string) (string, error) {
	query := `SELECT musicbrainz_artist_id
		from metadata
		where lower(artist) = lower(?)
		limit 1;`

	var result string

	err = DB.QueryRowContext(ctx, query, artistName).Scan(&result)

	return result, err
}

func GetArtistNameById(ctx context.Context, musicBrainzArtistId string) (string, error) {
	query := `SELECT artist FROM metadata WHERE musicbrainz_artist_id = ? limit 1;`

	var result string

	err := DB.QueryRowContext(ctx, query, musicBrainzArtistId).Scan(&result)
	if err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return result, nil
}
