package database

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func GetRandomSongs(ctx context.Context, count int, genre string, fromYear string, toYear string, musicFolderInt int, seed int, offset int) ([]types.SubsonicChild, error) {
	requestUser, err := GetUserByContext(ctx)
	if err != nil {
		return []types.SubsonicChild{}, err
	}

	var args []interface{}

	query := `WITH rand AS (
		SELECT m.rowid
		FROM metadata m
		join track_genres g
		where g.file_path = m.file_path`

	if genre != "" {
		query += ` and lower(g.genre) = lower(?)`
		args = append(args, genre)
	}

	if musicFolderInt != 0 {
		query += ` and m.music_folder_id = ?`
		args = append(args, musicFolderInt)
	}

	if fromYear != "" {
		query += ` and CAST(REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') AS INTEGER) >= ?`
		args = append(args, fromYear)
	}

	if toYear != "" {
		query += ` and CAST(REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') AS INTEGER) <= ?`
		args = append(args, toYear)
	}

	query += ` group by m.musicbrainz_track_id
		ORDER BY (m.rowid * ?) % 1000000
	  limit ? offset ?
	),`

	if seed == 0 {
		seed = logic.GenerateRandomInt(1, 10000000)
	}

	args = append(args, seed, count, offset)

	query += `
	base_tracks as (
		select m.musicbrainz_track_id as musicbrainz_track_id ,
			m.musicbrainz_album_id as musicbrainz_album_id,
			m.title as title,
			m.album as album,
			m.artist as artist,
			m.track_number as track_number,
			REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year,
			substr(m.genre,1,(instr(m.genre,';')-1)) as genre,
			m.musicbrainz_track_id as cover_art,
			m.size as size,
			m.duration as duration,
			m.bitrate as bit_rate,
			m.file_path as file_path,
			m.date_added as date_created,
			m.disc_number as disc_number,
			m.musicbrainz_artist_id as musicbrainz_artist_id,
			m.genre as genre_string,
			m.album_artist as album_artist,
			maa.musicbrainz_artist_id as album_artist_id,
			m.bit_depth as bit_depth,
			m.sample_rate as sample_rate,
			m.channels as channel_count,
			us.created_at AS starred_date,
			m.label,
			u.id as user_id,
			m.rowid as row_id
		from rand
		join metadata m on rand.rowid = m.rowid
		join user_music_folders f on f.folder_id = m.music_folder_id
		join users u on f.user_id = u.id
		left join track_genres g on m.file_path = g.file_path
		LEFT JOIN user_ratings ur ON m.musicbrainz_track_id = ur.metadata_id AND ur.user_id = f.user_id
		LEFT JOIN gr ON m.musicbrainz_track_id = gr.metadata_id
		LEFT JOIN pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND f.user_id = pc.user_id
		LEFT JOIN user_stars us ON m.musicbrainz_track_id = us.metadata_id AND us.user_id = f.user_id
		left join metadata maa on maa.artist = m.album_artist
		group by m.musicbrainz_track_id
	),
	gr as (
		SELECT metadata_id,
			AVG(rating) AS avg_rating
		FROM user_ratings
		GROUP BY metadata_id
	),
	pc AS (
		SELECT musicbrainz_track_id, user_id, SUM(play_count) AS play_count, MAX(last_played) AS played
		FROM play_counts 
		GROUP BY musicbrainz_track_id, user_id
	)
	select bt.musicbrainz_track_id,
		bt.musicbrainz_album_id,
		bt.title,
		bt.album,
		bt.artist,
		coalesce(bt.track_number,0) as track_number,
		bt.year,
		bt.genre,
		bt.cover_art,
		bt.size,
		bt.duration,
		bt.bit_rate,
		bt.file_path,
		bt.date_created,
		bt.disc_number,
		bt.musicbrainz_artist_id,
		bt.genre_string,
		bt.album_artist,
		bt.album_artist_id,
		bt.bit_depth,
		bt.sample_rate,
		bt.channel_count,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(gr.avg_rating, 0) AS average_rating,
		coalesce(pc.play_count, 0) as play_count,
		pc.played as last_played,
		bt.starred_date as starred_date,
		bt.label
	from rand
	join base_tracks bt on bt.row_id = rand.rowid
	join users u on bt.user_id = u.id
	LEFT JOIN user_ratings ur ON bt.musicbrainz_track_id = ur.metadata_id AND ur.user_id = bt.user_id
	LEFT JOIN gr ON bt.musicbrainz_track_id = gr.metadata_id
	LEFT JOIN pc ON bt.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = u.id
	where u.id = ?
	`

	args = append(args, requestUser.Id)

	rows, err := DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("getting scans: %v", err)
	}
	defer rows.Close()

	var songs []types.SubsonicChild
	for rows.Next() {
		var result types.SubsonicChild

		var genreString string
		var durationFloat float64
		var albumArtistName sql.NullString
		var albumArtistId sql.NullString
		var starred sql.NullString
		var played sql.NullString
		var labels sql.NullString

		result.IsDir = false
		result.MediaType = "song"
		result.Type = "music"
		result.IsVideo = false
		result.Bpm = 0
		result.Comment = ""
		result.Contributors = []types.ChildContributors{}
		result.Moods = []string{}

		if err := rows.Scan(&result.Id, &result.Parent, &result.Title, &result.Album, &result.Artist,
			&result.Track, &result.Year, &result.Genre, &result.CoverArt, &result.Size,
			&durationFloat, &result.BitRate, &result.Path, &result.Created, &result.DiscNumber,
			&result.ArtistId, &genreString, &albumArtistName, &albumArtistId, &result.BitDepth, &result.SamplingRate,
			&result.ChannelCount, &result.UserRating, &result.AverageRating, &result.PlayCount,
			&played, &starred, &labels); err != nil {
			logger.Printf("Failed to scan row in GetSongsByGenre: %v", err)
			return []types.SubsonicChild{}, err
		}
		if starred.Valid {
			result.Starred = starred.String
		}

		if played.Valid {
			result.Played = played.String
		}

		result.ContentType = logic.InferMimeTypeFromFileExtension(result.Path)
		result.Suffix = strings.Replace(filepath.Ext(result.Path), ".", "", 1)
		result.Duration = int(durationFloat)
		result.SortName = strings.ToLower(result.Title)
		result.MusicBrainzId = result.Id
		result.AlbumId = result.Parent

		result.Genres = []types.ChildGenre{}
		for _, genre := range strings.Split(genreString, ";") {
			result.Genres = append(result.Genres, types.ChildGenre{Name: genre})
		}

		result.Artists = []types.ChildArtist{}
		result.Artists = append(result.Artists, types.ChildArtist{Id: result.ArtistId, Name: result.Artist})

		result.DisplayArtist = result.Artist

		result.AlbumArtists = []types.ChildArtist{}
		if albumArtistId.Valid && albumArtistName.Valid {
			result.AlbumArtists = append(result.AlbumArtists, types.ChildArtist{Id: albumArtistId.String, Name: albumArtistName.String})
		}

		result.DisplayAlbumArtist = albumArtistName.String

		if labels.Valid {
			result.RecordLabels = []types.ChildRecordLabel{
				{Name: labels.String},
			}
		}

		songs = append(songs, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating song rows: %v", err)
	}

	return songs, nil
}
