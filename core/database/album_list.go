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

	query := `select m.musicbrainz_album_id as id,
		m.album as name,
		m.artist as artist,
		m.musicbrainz_album_id as cover_art,
		count(m.musicbrainz_track_id) as song_count,
		cast(sum(m.duration) as integer) as duration,
		COALESCE(SUM(pc.play_count), 0) as play_count,
		min(m.date_added) as created,
		m.musicbrainz_artist_id as artist_id,
		s.created_at as starred,
		REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year,
		substr(m.genre,1,(instr(m.genre,';')-1)) as genre,
		max(pc.last_played) as played,
		COALESCE(ur.rating, 0) AS user_rating,
		m.label as label_string,
		m.musicbrainz_album_id as musicbrainz_id,
		m.genre as genre_string,
		m.artist as display_artist,
		lower(m.album) as sort_name,
		m.release_date as release_date_string,
		m.album_artist,
		maa.musicbrainz_artist_id
	from user_music_folders f
	join metadata m on m.music_folder_id = f.folder_id
	LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = f.user_id
	LEFT JOIN user_stars s ON m.musicbrainz_album_id = s.metadata_id AND s.user_id = f.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_artist_id = ur.metadata_id AND ur.user_id = f.user_id
	join metadata maa on maa.artist = m.album_artist`

	if sortType == "bygenre" {
		query += ` join track_genres tg on tg.file_path = m.file_path and lower(tg.genre) = lower(?)`
		args = append(args, genre)
	}

	query += ` where f.user_id = ?`
	args = append(args, user.Id)

	if musicFolderId != 0 {
		query += ` and m.music_folder_id = ?`
		args = append(args, musicFolderId)
	}

	if sortType == "starred" {
		query += ` and s.created_at is not null`
	}

	if sortType == "byyear" {
		query += ` and year >= ? and year <= ?`
		args = append(args, fromYear, toYear)
	}

	query += ` group by m.musicbrainz_album_id`

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
			{Name: albumArtistName, Id: albumArtistId},
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
