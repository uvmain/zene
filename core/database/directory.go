package database

import (
	"context"
	"database/sql"
	"strings"
	"zene/core/logger"
	"zene/core/types"
)

func GetAlbumDirectory(ctx context.Context, musicbrainzAlbumId string) (types.SubsonicDirectory, error) {
	directory := types.SubsonicDirectory{}
	directory.Id = musicbrainzAlbumId
	directory.CoverArt = musicbrainzAlbumId

	logger.Printf("Getting album directory for ID: %s", musicbrainzAlbumId)

	user, err := GetUserByContext(ctx)
	if err != nil {
		return types.SubsonicDirectory{}, err
	}

	query := `SELECT m.musicbrainz_artist_id AS parent,
		m.album AS name,
		s.created_at AS starred,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(AVG(gr.rating), 0.0) AS average_rating,
		COALESCE(SUM(pc.play_count), 0) AS play_count,
		COUNT(m.musicbrainz_track_id) AS song_count
	FROM metadata m
	JOIN user_music_folders f ON f.folder_id = m.music_folder_id
	LEFT JOIN user_stars s ON m.musicbrainz_album_id = s.metadata_id AND s.user_id = f.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_album_id = ur.metadata_id AND ur.user_id = f.user_id
	LEFT JOIN user_ratings gr ON m.musicbrainz_album_id = gr.metadata_id
	LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = f.user_id
	WHERE f.user_id = ? and m.musicbrainz_album_id = ?
	GROUP BY m.musicbrainz_album_id, m.album, s.created_at, ur.rating;`

	var starred sql.NullString

	err = DB.QueryRowContext(ctx, query, user.Id, musicbrainzAlbumId).Scan(
		&directory.Parent, &directory.Name, &starred, &directory.UserRating, &directory.AverageRating, &directory.PlayCount, &directory.SongCount,
	)

	if err == sql.ErrNoRows {
		return types.SubsonicDirectory{}, nil
	} else if err != nil {
		return types.SubsonicDirectory{}, err
	}

	if starred.Valid {
		directory.Starred = starred.String
	}

	directory.MediaType = "album"

	children, err := GetSongsForAlbum(ctx, musicbrainzAlbumId)
	if err != nil {
		return types.SubsonicDirectory{}, err
	}

	directory.Child = children

	return directory, nil
}

func GetArtistDirectory(ctx context.Context, musicbrainzArtistId string) (types.SubsonicDirectory, error) {
	directory := types.SubsonicDirectory{}
	directory.Id = musicbrainzArtistId
	directory.CoverArt = musicbrainzArtistId
	directory.MediaType = "artist"

	logger.Printf("Getting artist directory for ID: %s", musicbrainzArtistId)

	user, err := GetUserByContext(ctx)
	if err != nil {
		return types.SubsonicDirectory{}, err
	}

	query := `SELECT m.artist AS name,
		s.created_at AS starred,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(AVG(gr.rating), 0.0) AS average_rating,
		COALESCE(SUM(pc.play_count), 0) AS play_count,
		COUNT(m.musicbrainz_track_id) AS song_count,
		COUNT(distinct m.musicbrainz_album_id) AS album_count
	FROM metadata m
	JOIN user_music_folders f ON f.folder_id = m.music_folder_id AND f.user_id = 1
	LEFT JOIN user_stars s ON m.musicbrainz_album_id = s.metadata_id AND s.user_id = f.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_artist_id  = ur.metadata_id AND ur.user_id = f.user_id
	LEFT JOIN user_ratings gr ON m.musicbrainz_artist_id  = gr.metadata_id
	LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = f.user_id
	WHERE f.user_id = ? and m.musicbrainz_artist_id = ? 
	GROUP BY m.musicbrainz_artist_id, s.created_at, ur.rating;`

	var starred sql.NullString

	err = DB.QueryRowContext(ctx, query, user.Id, musicbrainzArtistId).Scan(
		&directory.Name, &starred, &directory.UserRating, &directory.AverageRating,
		&directory.PlayCount, &directory.SongCount, &directory.AlbumCount,
	)

	if err == sql.ErrNoRows {
		return types.SubsonicDirectory{}, nil
	} else if err != nil {
		return types.SubsonicDirectory{}, err
	}

	if starred.Valid {
		directory.Starred = starred.String
	}

	children, err := GetArtistChildren(ctx, musicbrainzArtistId)
	if err != nil {
		return types.SubsonicDirectory{}, err
	}

	directory.Child = children

	return directory, nil
}

func GetArtistChildren(ctx context.Context, musicbrainzArtistId string) ([]types.SubsonicChild, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return []types.SubsonicChild{}, err
	}

	children := []types.SubsonicChild{}

	query := `with album_plays AS (
		SELECT	m.musicbrainz_album_id,
			SUM(pc.play_count) AS play_count,
			MAX(pc.last_played) AS last_played,
			pc.user_id
		FROM play_counts pc
		join metadata m ON m.musicbrainz_track_id = pc.musicbrainz_track_id
		where pc.user_id = ?
		GROUP BY m.musicbrainz_album_id
	),
	album_artists as (
		select musicbrainz_album_id, musicbrainz_artist_id, album_artist
		from metadata
		where album_artist = artist
		group by musicbrainz_album_id
	),
	album_song_counts as (
		select count(musicbrainz_track_id) as track_count, musicbrainz_album_id
		from metadata
		group by musicbrainz_album_id
	)
	select m.musicbrainz_album_id as id, m.musicbrainz_artist_id as parent,
		m.album, coalesce(maa.album_artist, m.artist) as artist, REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year,
		substr(m.genre,1,(instr(m.genre,';')-1)) as genre, m.musicbrainz_album_id as cover_art,
		sum(m.duration) as duration, min(m.date_added) as created, m.label as label,
		m.album_artist, m.genre as genres, m.musicbrainz_artist_id as musicbrainz_artist,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(AVG(gr.rating), 0.0) AS average_rating,
		COALESCE(pc.play_count, 0) AS play_count,
		sc.track_count AS song_count,
		maa.musicbrainz_artist_id
	from metadata m
	join user_music_folders f on f.folder_id = m.music_folder_id
	LEFT JOIN user_stars s ON s.metadata_id = m.musicbrainz_album_id AND s.user_id = f.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_album_id = ur.metadata_id AND ur.user_id = f.user_id
	LEFT JOIN user_ratings gr ON m.musicbrainz_album_id = gr.metadata_id
	LEFT JOIN album_plays pc ON pc.musicbrainz_album_id = m.musicbrainz_album_id
	left join album_artists maa on maa.album_artist = m.album_artist
	join album_song_counts sc on sc.musicbrainz_album_id = m.musicbrainz_album_id
	where (m.musicbrainz_artist_id = ? or maa.musicbrainz_artist_id == ?) and f.user_id = ?
	group by m.musicbrainz_album_id
	order by m.release_date desc;`

	rows, err := DB.Query(query, user.Id, musicbrainzArtistId, musicbrainzArtistId, user.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var child types.SubsonicChild
		var albumArtistName sql.NullString
		var albumArtistId sql.NullString
		var genreString string
		var labelString string
		var durationFloat float64

		if err := rows.Scan(&child.Id, &child.Parent, &child.Album, &child.Artist, &child.Year,
			&child.Genre, &child.CoverArt, &durationFloat, &child.Created, &labelString, &albumArtistName, &genreString,
			&child.ArtistId, &child.UserRating, &child.AverageRating, &child.PlayCount, &child.SongCount, &albumArtistId); err != nil {
			return nil, err
		}
		child.Genres = []types.ChildGenre{}
		for _, genre := range strings.Split(genreString, ";") {
			child.Genres = append(child.Genres, types.ChildGenre{Name: genre})
		}

		child.Duration = int(durationFloat)
		child.Title = child.Album
		child.IsDir = true

		child.RecordLabels = []types.ChildRecordLabel{}
		child.RecordLabels = append(child.RecordLabels, types.ChildRecordLabel{Name: labelString})

		child.Artists = []types.ChildArtist{}
		child.Artists = append(child.Artists, types.ChildArtist{Id: child.ArtistId, Name: child.Artist})

		child.DisplayArtist = child.Artist

		child.AlbumArtists = []types.ChildArtist{}
		if albumArtistId.Valid && albumArtistName.Valid {
			child.AlbumArtists = append(child.AlbumArtists, types.ChildArtist{Id: albumArtistId.String, Name: albumArtistName.String})
		}

		child.DisplayAlbumArtist = albumArtistName.String
		children = append(children, child)
	}

	return children, nil
}
