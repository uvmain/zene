package database

import (
	"context"
	"database/sql"
	"fmt"
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
		coalesce(count(distinct m.musicbrainz_album_id), 0) as album_count,
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
	result.ArtistImageUrl = logic.GetUnauthenticatedImageUrl(result.Id, 600)
	result.MusicBrainzId = result.Id
	result.SortName = strings.ToLower(result.Name)

	if starred.Valid {
		result.Starred = starred.String
	}

	albums, err := GetArtistChildren(ctx, musicbrainzArtistId)
	if err != nil {
		return result, fmt.Errorf("getting artist children: %v", err)
	}

	result.Album = []types.SubsonicChild{}
	result.Album = append(result.Album, albums...)

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

func GetArtistList(ctx context.Context, musicFolderIds []int, limit int, offset int, sortType string, seed int) ([]types.Artist, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return []types.Artist{}, err
	}

	var args []interface{}
	query := `with artist_plays AS (
    SELECT m.musicbrainz_artist_id,
			SUM(pc.play_count) AS play_count,
			MAX(pc.last_played) AS last_played
    FROM play_counts pc
		join metadata m ON m.musicbrainz_track_id = pc.musicbrainz_track_id
		and pc.user_id = ?
		`
	args = append(args, user.Id)

	query += ` GROUP BY m.musicbrainz_artist_id
	),
	stars AS (
		select metadata_id,
			MAX(created_at) AS date_starred
		from user_stars
		where user_id = ?
		group by metadata_id
	),
	album_count AS (
		select rowid,
			musicbrainz_artist_id,
			count(distinct musicbrainz_album_id) as album_count
		from metadata
		where artist = album_artist
		group by musicbrainz_artist_id
	),
	gr AS (
			SELECT metadata_id, AVG(rating) AS avg_rating
			FROM user_ratings
			GROUP BY metadata_id
		)
	SELECT
		m.musicbrainz_artist_id,
		m.artist,
		s.date_starred,
		COALESCE(max(ur.rating), 0) AS user_rating,
		COALESCE(gr.avg_rating, 0.0) AS average_rating,
		coalesce(ac.album_count, 0) AS album_count
	FROM metadata m
	left join album_count ac on ac.musicbrainz_artist_id = m.musicbrainz_artist_id
	JOIN user_music_folders mf ON mf.folder_id = m.music_folder_id
	JOIN users u ON u.id = mf.user_id AND u.id = ?
	LEFT JOIN artist_plays ap ON m.musicbrainz_artist_id = ap.musicbrainz_artist_id
	LEFT JOIN stars s ON m.musicbrainz_artist_id = s.metadata_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_artist_id = ur.metadata_id AND ur.user_id = u.id
	LEFT JOIN gr ON m.musicbrainz_artist_id = gr.metadata_id
	`

	args = append(args, user.Id, user.Id)

	if len(musicFolderIds) > 0 {
		query += ` WHERE m.music_folder_id IN (`
		for i := 0; i < len(musicFolderIds); i++ {
			if i > 0 {
				query += `,`
			}
			query += `?`
			args = append(args, musicFolderIds[i])
		}
		query += `)`
	}

	query += ` GROUP BY m.artist`

	switch sortType {
	case "random":
		if seed != 0 {
			query += fmt.Sprintf(" order BY (m.rowid * %d) %% 1000000", seed)
		} else {
			query += " order BY random()"
		}
	case "newest": // recently added artists
		query += " order by m.date_added desc"
	case "highest": // highest rated artists
		query += " order by ur.rating desc, m.artist desc"
	case "frequent": // most frequently played artists
		query += " order by ap.play_count desc, m.date_added desc"
	case "recent": // recently played artists
		query += " order by last_played desc, m.date_added desc"
	case "alphabetical":
		query += " order by m.artist asc"
	case "starred":
		query += ` having date_starred is not null
			order by date_starred desc, m.artist asc`
	default:
		query += " order BY m.artist asc"
	}

	if limit > 0 {
		query += ` limit ? offset ?`
		args = append(args, limit, offset)
	}

	rows, err := DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying GetArtistList: %v", err)
	}
	defer rows.Close()

	var artists []types.Artist

	for rows.Next() {
		var artist types.Artist
		var starred sql.NullString
		if err := rows.Scan(&artist.Id, &artist.Name, &starred, &artist.UserRating, &artist.AverageRating, &artist.AlbumCount); err != nil {
			return nil, fmt.Errorf("scanning artist row: %v", err)
		}
		artist.CoverArt = artist.Id
		artist.ArtistImageUrl = logic.GetUnauthenticatedImageUrl(artist.Id, 600)
		artist.MusicBrainzId = artist.Id
		artist.SortName = strings.ToLower(artist.Name)

		if starred.Valid {
			artist.Starred = starred.String
		}
		artists = append(artists, artist)
	}

	return artists, nil
}
