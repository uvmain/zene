package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"

	"github.com/timematic/anytime"
)

func GetAlbumList(ctx context.Context, sortType string, limit int, offset int, fromYear string, toYear string, genre string, musicFolderId int) ([]types.AlbumId3, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return []types.AlbumId3{}, err
	}

	var args []interface{}
	var albums []types.AlbumId3

	query := `WITH album_stats AS (
  SELECT
    m.musicbrainz_album_id,
    COUNT(m.musicbrainz_track_id) AS song_count,
    CAST(SUM(m.duration) AS INTEGER) AS duration,
    COALESCE(SUM(pc.play_count), 0) AS play_count,
    MIN(m.date_added) AS created,
    MAX(pc.last_played) AS played
  FROM metadata m
  JOIN user_music_folders f ON m.music_folder_id = f.folder_id
  LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = f.user_id
  WHERE f.user_id = ?`
	args = append(args, user.Id)

	if musicFolderId != 0 {
		query += ` and m.music_folder_id = ?`
		args = append(args, musicFolderId)
	}

	query += `GROUP BY m.musicbrainz_album_id
	)
	SELECT
		m.musicbrainz_album_id AS id,
		m.album AS name,
		m.artist,
		m.musicbrainz_album_id AS cover_art,
		a.song_count,
		a.duration,
		a.play_count,
		a.created,
		m.musicbrainz_artist_id AS artist_id,
		s.created_at AS starred,
		REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') AS year,
		substr(m.genre,1,(instr(m.genre,';')-1)) AS genre,
		a.played,
		COALESCE(ur.rating, 0) AS user_rating,
		m.label AS label_string,
		m.musicbrainz_album_id AS musicbrainz_id,
		m.genre AS genre_string,
		m.artist AS display_artist,
		lower(m.album) AS sort_name,
		m.release_date AS release_date_string,
		m.album_artist,
		maa.musicbrainz_artist_id
	FROM album_stats a
	JOIN metadata m ON m.musicbrainz_album_id = a.musicbrainz_album_id
	LEFT JOIN user_stars s ON m.musicbrainz_album_id = s.metadata_id AND s.user_id = ?
	LEFT JOIN user_ratings ur ON m.musicbrainz_artist_id = ur.metadata_id AND ur.user_id = ?
	LEFT JOIN metadata maa ON maa.artist = m.album_artist`
	args = append(args, user.Id, user.Id)

	if sortType == "bygenre" {
		query += ` join track_genres tg on tg.file_path = m.file_path and lower(tg.genre) = lower(?)`
		args = append(args, genre)
	}

	if sortType == "starred" {
		query += ` and s.created_at is not null`
	}

	if sortType == "byyear" {
		query += ` and year >= ? and year <= ?`
		args = append(args, fromYear, toYear)
	}

	// query += ` group by m.musicbrainz_album_id`

	switch sortType {
	case "random":
		randomInt := logic.GenerateRandomInt(1, 10000000)
		query += fmt.Sprintf(" order BY ((m.rowid * %d) %% 1000000)", randomInt)
	case "newest": // recently added albums
		query += " order BY m.date_added desc"
	case "highest": // highest rated albums
		query += " order by ur.rating desc, m.musicbrainz_album_id desc"
	case "frequent": // most frequently played albums
		query += " order by pc.play_count desc, m.musicbrainz_album_id desc"
	case "recent": // recently played albums
		query += " order by pc.last_played desc, m.musicbrainz_album_id desc"
	case "alphabeticalbyname":
		query += " order by m.album asc"
	case "alphabeticalbyartist":
		query += " order by m.artist asc"
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
		var albumArtistId string

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
			{Name: albumArtistName, Id: albumArtistId},
		}

		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return albums, err
	}

	return albums, nil
}
