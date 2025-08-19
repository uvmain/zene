package database

import (
	"context"
	"database/sql"
	"zene/core/types"
)

func GetAlbumDirectory(ctx context.Context, musicbrainzAlbumId string) (types.Directory, error) {
	directory := types.Directory{}
	directory.Id = musicbrainzAlbumId
	directory.CoverArt = musicbrainzAlbumId

	user, err := GetUserByContext(ctx)
	if err != nil {
		return types.Directory{}, err
	}

	query := `SELECT 
		m.musicbrainz_artist_id AS parent,
		m.album AS name,
		s.created_at AS starred,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(AVG(gr.rating), 0.0) AS average_rating,
		COALESCE(SUM(pc.play_count), 0) AS play_count,
		COUNT(m.musicbrainz_track_id) AS song_count
	FROM metadata m
	JOIN user_music_folders f ON f.folder_id = m.music_folder_id AND f.user_id = 1
	LEFT JOIN user_stars s ON m.musicbrainz_artist_id = s.metadata_id AND s.user_id = f.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_artist_id = ur.metadata_id AND ur.user_id = f.user_id
	LEFT JOIN user_ratings gr ON m.musicbrainz_artist_id = gr.metadata_id
	LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = f.user_id
	WHERE f.user_id = ? and m.musicbrainz_album_id = ?
	GROUP BY m.musicbrainz_artist_id, m.album, s.created_at, ur.rating;;`

	var starred sql.NullString

	err = DB.QueryRowContext(ctx, query, user.Id, musicbrainzAlbumId).Scan(
		&directory.Parent, &directory.Name, &starred, &directory.UserRating, &directory.AverageRating, &directory.PlayCount, &directory.SongCount,
	)

	if err == sql.ErrNoRows {
		return types.Directory{}, nil
	} else if err != nil {
		return types.Directory{}, err
	}

	if starred.Valid {
		directory.Starred = starred.String
	}

	children, err := GetSongsForAlbum(ctx, musicbrainzAlbumId)
	if err != nil {
		return types.Directory{}, err
	}

	directory.Child = children

	return directory, nil
}
