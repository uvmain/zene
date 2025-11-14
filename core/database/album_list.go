package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"zene/core/logger"
	"zene/core/types"

	"github.com/timematic/anytime"
)

func GetAlbumList(ctx context.Context, sortType string, limit int, offset int, fromYear int, toYear int, genre string, musicFolderId int, seed int) ([]types.AlbumId3, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return []types.AlbumId3{}, err
	}

	var args []interface{}
	var albums []types.AlbumId3

	query := `WITH album_stats AS (
    select	m.rowid,
		musicbrainz_album_id,
		COUNT(*) AS song_count,
		CAST(SUM(duration) AS INTEGER) AS duration
	from metadata m
	group by musicbrainz_album_id
	),
	album_plays AS (
    SELECT
			m.musicbrainz_album_id,
			SUM(pc.play_count) AS total_play_count,
			MAX(pc.last_played) AS last_played,
			pc.user_id
    FROM play_counts pc
	join metadata m ON m.musicbrainz_track_id = pc.musicbrainz_track_id
	where pc.user_id = ?
	GROUP BY m.musicbrainz_album_id
	),
	album_artist_map AS (
    SELECT
        artist,
        MIN(musicbrainz_artist_id) AS musicbrainz_artist_id
    FROM metadata
    WHERE musicbrainz_artist_id IS NOT NULL
	GROUP BY artist
	)
	SELECT
    m.musicbrainz_album_id AS id,
    m.album AS name,
    coalesce(maa.artist, m.artist),
    m.musicbrainz_album_id AS cover_art,
    a.song_count,
    a.duration,
    COALESCE(ap.total_play_count, 0) as play_count,
    m.date_added as created,
    m.musicbrainz_artist_id AS artist_id,
    s.created_at AS starred,
    CAST(REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') AS INTEGER) AS year,
    substr(m.genre,1,(instr(m.genre,';')-1)) AS genre,
    ap.last_played as played,
    COALESCE(ur.rating, 0) AS user_rating,
    m.label AS label_string,
    m.musicbrainz_album_id AS musicbrainz_id,
    m.genre AS genre_string,
    coalesce(maa.artist, m.artist) AS display_artist,
    LOWER(m.album) AS sort_name,
    m.release_date AS release_date_string,
    m.album_artist,
    maa.musicbrainz_artist_id as album_artist_id
	FROM album_stats a
	join metadata m on m.rowid = a.rowid and m.musicbrainz_album_id = a.musicbrainz_album_id
	join user_music_folders f on f.folder_id = m.music_folder_id
	join users u on u.id = f.user_id
	LEFT JOIN album_plays ap ON ap.musicbrainz_album_id = m.musicbrainz_album_id and ap.user_id = u.id
	LEFT JOIN user_ratings ur ON m.musicbrainz_album_id = ur.metadata_id AND ur.user_id = u.id
	LEFT JOIN album_artist_map maa ON maa.artist = m.album_artist
	LEFT JOIN user_stars s ON m.musicbrainz_album_id = s.metadata_id AND s.user_id = u.id
	`

	args = append(args, user.Id)

	if sortType == "bygenre" {
		query += ` join track_genres tg on tg.file_path = m.file_path and lower(tg.genre) = lower(?)`
		args = append(args, genre)
	}

	query += ` where u.id = ?`
	args = append(args, user.Id)

	if musicFolderId != 0 {
		query += ` and m.music_folder_id = ?`
		args = append(args, musicFolderId)
	}

	if sortType == "starred" {
		query += ` and starred is not null`
	}

	var yearSortOrder = 1
	if sortType == "byyear" {
		if fromYear > toYear {
			fromYear, toYear = toYear, fromYear
			yearSortOrder = -1
		}
		query += ` and year >= ? and year <= ?`
		args = append(args, fromYear, toYear)
	}

	switch sortType {
	case "random":
		if seed != 0 {
			query += fmt.Sprintf(" order BY (m.rowid * %d) %% 1000000", seed)
		} else {
			query += " order BY random()"
		}
	case "newest": // recently added albums
		query += " order BY m.date_added desc"
	case "byyear":
		if yearSortOrder == -1 {
			query += " order by year desc"
		} else {
			query += " order by year asc"
		}
	case "highest": // highest rated albums
		query += " order by ur.rating desc, m.musicbrainz_album_id desc"
	case "frequent": // most frequently played albums
		query += " order by play_count desc, m.musicbrainz_album_id desc"
	case "recent": // recently played albums
		query += " order by last_played desc, m.musicbrainz_album_id desc"
	case "alphabeticalbyname":
		query += " order by m.album asc"
	case "alphabeticalbyartist":
		query += " order by m.artist asc"
	case "release":
		query += " order by m.release_date desc"
	default:
		query += " order BY m.musicbrainz_album_id asc"
	}

	query += ` limit ? offset ?`
	args = append(args, limit, offset)

	rows, err := DB.QueryContext(ctx, query, args...)
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
		var albumArtistName string
		var albumArtistId sql.NullString

		if err := rows.Scan(&album.Id, &album.Name, &album.Artist, &album.CoverArt, &album.SongCount,
			&album.Duration, &album.PlayCount, &album.Created, &album.ArtistId, &starred,
			&album.Year, &album.Genre, &played, &album.UserRating, &labelString, &album.MusicBrainzId,
			&genresString, &album.DisplayArtist, &album.SortName, &releaseDateString,
			&albumArtistName, &albumArtistId); err != nil {
			logger.Printf("Failed to scan row in GetAlbumList: %v", err)
			return nil, err
		}

		if starred.Valid {
			album.Starred = starred.String
		}

		if played.Valid {
			album.Played = played.String
		}

		var albumArtistIdString string
		if albumArtistId.Valid {
			albumArtistIdString = albumArtistId.String
		} else {
			albumArtistIdString = album.ArtistId
		}

		album.Title = album.Name
		album.Album = album.Name

		album.RecordLabels = []types.ChildRecordLabel{}
		album.RecordLabels = append(album.RecordLabels, types.ChildRecordLabel{Name: labelString.String})

		album.Genres = []types.ItemGenre{}
		for _, genre := range strings.Split(genresString.String, ";") {
			album.Genres = append(album.Genres, types.ItemGenre{Name: genre})
		}

		releaseDateTime, err := anytime.Parse(releaseDateString.String)
		if err == nil {
			album.ReleaseDate = types.ItemDate{
				Year:  releaseDateTime.Year(),
				Month: int(releaseDateTime.Month()),
				Day:   releaseDateTime.Day(),
			}
		}

		album.Artists = []types.Artist{
			{Name: album.Artist, Id: album.ArtistId},
		}

		album.AlbumArtists = []types.Artist{
			{Name: albumArtistName, Id: albumArtistIdString},
		}

		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return albums, err
	}

	return albums, nil
}
