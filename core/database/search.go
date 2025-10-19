package database

import (
	"context"
	"database/sql"
	"strings"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"

	"github.com/timematic/anytime"
)

func SearchArtists(ctx context.Context, searchQuery string, limit int, offset int, musicFolderId int) ([]types.Artist, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return []types.Artist{}, err
	}

	var args []interface{}
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
	where u.user_id = ?`

	args = append(args, user.Id)

	if searchQuery != "" {
		query += ` and lower(m.artist) like lower(?)`
		args = append(args, "%"+searchQuery+"%")
	}

	if musicFolderId != 0 {
		query += ` and m.music_folder_id = ?`
		args = append(args, musicFolderId)
	}

	query += ` group by m.musicbrainz_artist_id
	order by m.musicbrainz_artist_id asc
	limit ?
	offset ?`
	args = append(args, limit, offset)

	var rows *sql.Rows

	rows, err = DB.QueryContext(ctx, query, args...)
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
			logger.Printf("Failed to scan row in SearchArtists: %v", err)
			return nil, err
		}

		result.CoverArt = result.Id
		result.ArtistImageUrl = logic.GetUnauthenticatedImageUrl(result.Id, 600)
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

func SearchAlbums(ctx context.Context, searchQuery string, limit int, offset int, musicFolderId int) ([]types.AlbumId3, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return []types.AlbumId3{}, err
	}

	var albums []types.AlbumId3
	var args []interface{}
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
		)
		select m.musicbrainz_album_id as id,
			m.album as name,
			coalesce(maa.album_artist, m.album_artist) as artist,
			m.musicbrainz_album_id as cover_art,
			count(m.musicbrainz_track_id) as song_count,
			cast(sum(m.duration) as integer) as duration,
			COALESCE(ap.play_count, 0) as play_count,
			min(m.date_added) as created,
			m.musicbrainz_artist_id as artist_id,
			s.created_at as starred,
			REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year,
			substr(m.genre,1,(instr(m.genre,';')-1)) as genre,
			max(ap.last_played) as played,
			COALESCE(ur.rating, 0) AS user_rating,
			m.label as label_string,
			m.musicbrainz_album_id as musicbrainz_id,
			m.genre as genre_string,
			m.artist as display_artist,
			lower(m.album) as sort_name,
			m.release_date as release_date_string,
			maa.musicbrainz_artist_id as album_artist_id,
			coalesce(maa.album_artist, m.album_artist) as album_artist_name
		from metadata m
		join user_music_folders f on m.music_folder_id = f.folder_id
		LEFT JOIN album_plays ap ON ap.musicbrainz_album_id = m.musicbrainz_album_id
		LEFT JOIN user_stars s ON m.musicbrainz_album_id = s.metadata_id AND s.user_id = f.user_id
		LEFT JOIN user_ratings ur ON m.musicbrainz_album_id = ur.metadata_id AND ur.user_id = f.user_id
		left join album_artists maa on maa.musicbrainz_album_id = m.musicbrainz_album_id
		where f.user_id = ?`

	args = append(args, user.Id, user.Id)

	if searchQuery != "" {
		query += ` and lower(m.album) like lower(?)`
		args = append(args, "%"+searchQuery+"%")
	}

	if musicFolderId != 0 {
		query += ` and m.music_folder_id = ?`
		args = append(args, musicFolderId)
	}

	query += ` group by m.musicbrainz_album_id
	order by m.musicbrainz_album_id asc
	limit ? offset ?`
	args = append(args, limit, offset)

	var rows *sql.Rows

	rows, err = DB.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.AlbumId3{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var album types.AlbumId3

		var starred sql.NullString
		var labelString sql.NullString
		var genresString sql.NullString
		var releaseDateString sql.NullString
		var played sql.NullString
		var albumArtistId sql.NullString
		var albumArtistName sql.NullString

		if err := rows.Scan(&album.Id, &album.Name, &album.Artist, &album.CoverArt, &album.SongCount,
			&album.Duration, &album.PlayCount, &album.Created, &album.ArtistId, &starred,
			&album.Year, &album.Genre, &played, &album.UserRating,
			&labelString, &album.MusicBrainzId, &genresString,
			&album.DisplayArtist, &album.SortName, &releaseDateString,
			&albumArtistId, &albumArtistName); err != nil {
			logger.Printf("Failed to scan row in SearchAlbums: %v", err)
			return nil, err
		}

		if starred.Valid {
			album.Starred = starred.String
		}

		if played.Valid {
			album.Played = played.String
		}

		album.Title = album.Name
		album.Album = album.Name

		album.RecordLabels = []types.ChildRecordLabel{}
		album.RecordLabels = append(album.RecordLabels, types.ChildRecordLabel{Name: labelString.String})

		album.Genres = []types.ItemGenre{}
		for _, genre := range strings.Split(genresString.String, ";") {
			album.Genres = append(album.Genres, types.ItemGenre{Name: genre})
		}

		album.Artists = []types.Artist{
			{Id: album.ArtistId, Name: album.Artist},
		}

		album.AlbumArtists = []types.Artist{}
		if albumArtistId.Valid && albumArtistName.Valid {
			album.AlbumArtists = append(album.AlbumArtists, types.Artist{Id: albumArtistId.String, Name: albumArtistName.String})
		}

		releaseDateTime, err := anytime.Parse(releaseDateString.String)
		if err == nil {
			album.ReleaseDate = types.ItemDate{
				Year:  releaseDateTime.Year(),
				Month: int(releaseDateTime.Month()),
				Day:   releaseDateTime.Day(),
			}
		}

		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return albums, err
	}

	return albums, nil
}

func SearchSongs(ctx context.Context, searchQuery string, limit int, offset int, musicFolderId int) ([]types.SubsonicChild, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return []types.SubsonicChild{}, err
	}

	var args []interface{}
	query := `select m.musicbrainz_track_id as id,
		m.musicbrainz_album_id as parent,
		m.title,
		m.album,
		m.artist,
		COALESCE(m.track_number, 0) as track,
		REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year,
		substr(m.genre,1,(instr(m.genre,';')-1)) as genre,
		m.musicbrainz_track_id as cover_art,
		m.size,
		m.duration,
		m.bitrate,
		m.file_path as path,
		m.date_added as created,
		m.disc_number,
		m.musicbrainz_artist_id as artist_id,
		m.genre, m.album_artist, m.bit_depth, m.sample_rate, m.channels,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(AVG(gr.rating), 0.0) AS average_rating,
		COALESCE(SUM(pc.play_count), 0) AS play_count,
		max(pc.last_played) as played,
		us.created_at AS starred,
		maa.musicbrainz_artist_id as album_artist_id
	from user_music_folders u
	join metadata m on m.music_folder_id = u.folder_id
	LEFT JOIN user_stars us ON m.musicbrainz_track_id = us.metadata_id AND us.user_id = u.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_track_id = ur.metadata_id AND ur.user_id = u.user_id
	LEFT JOIN user_ratings gr ON m.musicbrainz_track_id = gr.metadata_id
	LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = u.user_id
	JOIN metadata maa ON maa.musicbrainz_album_id = m.musicbrainz_album_id
	where u.user_id = ?`
	args = append(args, user.Id)

	if searchQuery != "" {
		query += ` and lower(m.title) like lower(?)`
		args = append(args, "%"+searchQuery+"%")
	}

	if musicFolderId != 0 {
		query += ` and m.music_folder_id = ?`
		args = append(args, musicFolderId)
	}

	query += ` group by m.musicbrainz_track_id
	order by m.musicbrainz_track_id asc
	limit ? offset ?`
	args = append(args, limit, offset)

	var rows *sql.Rows
	var results []types.SubsonicChild

	rows, err = DB.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.SubsonicChild{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var result types.SubsonicChild

		var albumArtistName sql.NullString
		var albumArtistId sql.NullString
		var genreString string
		var durationFloat float64
		var played sql.NullString
		var starred sql.NullString

		if err := rows.Scan(&result.Id, &result.Parent, &result.Title, &result.Album, &result.Artist, &result.Track,
			&result.Year, &result.Genre, &result.CoverArt,
			&result.Size, &durationFloat, &result.BitRate, &result.Path, &result.Created, &result.DiscNumber, &result.ArtistId,
			&genreString, &albumArtistName, &result.BitDepth, &result.SamplingRate, &result.ChannelCount,
			&result.UserRating, &result.AverageRating, &result.PlayCount, &played, &starred, &albumArtistId); err != nil {
			return nil, err
		}
		result.Genres = []types.ChildGenre{}
		for _, genre := range strings.Split(genreString, ";") {
			result.Genres = append(result.Genres, types.ChildGenre{Name: genre})
		}

		if played.Valid {
			result.Played = played.String
		}
		if starred.Valid {
			result.Starred = starred.String
		}

		result.Duration = int(durationFloat)
		result.IsDir = false
		result.MusicBrainzId = result.Id
		result.AlbumId = result.Parent

		result.Artists = []types.ChildArtist{}
		result.Artists = append(result.Artists, types.ChildArtist{Id: result.ArtistId, Name: result.Artist})

		result.DisplayArtist = result.Artist

		result.AlbumArtists = []types.ChildArtist{}
		if albumArtistId.Valid && albumArtistName.Valid {
			result.AlbumArtists = append(result.AlbumArtists, types.ChildArtist{Id: albumArtistId.String, Name: albumArtistName.String})
		}

		result.DisplayAlbumArtist = albumArtistName.String

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return results, err
	}

	return results, nil
}
