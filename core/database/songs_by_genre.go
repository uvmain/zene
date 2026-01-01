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

func GetSongsByGenre(ctx context.Context, genre string, count int, offset int, musicFolderInt int) ([]types.SubsonicChild, error) {
	requestUser, err := GetUserByContext(ctx)
	if err != nil {
		return []types.SubsonicChild{}, err
	}

	var args []interface{}
	query := `with track_plays AS (
    SELECT musicbrainz_track_id,
			SUM(play_count) AS play_count,
			MAX(last_played) AS last_played
    FROM play_counts
		where user_id = ?
		group by musicbrainz_track_id
	),
	starred as (
		select metadata_id,
			created_at
		from user_stars
		where user_id = ?
	),
	album_artist_map AS (
		SELECT artist,
			MIN(musicbrainz_artist_id) AS musicbrainz_artist_id
		FROM metadata
		WHERE musicbrainz_artist_id IS NOT NULL
		GROUP BY artist
	),
	gr AS (
		SELECT metadata_id, AVG(rating) AS avg_rating
		FROM user_ratings
		GROUP BY metadata_id
	)
	select m.musicbrainz_track_id as id,
		m.musicbrainz_album_id as album_id,
		m.title, m.album,
		m.artist,
		COALESCE(m.track_number, 0) as track,
		REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year,
		substr(m.genre,1,(instr(m.genre,';')-1)) as genre,
		m.musicbrainz_track_id as cover_art,
		m.size,
		m.duration,
		m.bitrate,
		m.file_path as path,
		m.date_added as created, m.disc_number, m.musicbrainz_artist_id as artist_id,
		m.genre, m.album_artist, m.bit_depth, m.sample_rate, m.channels,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(gr.avg_rating, 0.0) AS average_rating,
		COALESCE(pc.play_count, 0) AS play_count,
		pc.last_played as played,
		s.created_at AS starred,
		maa.musicbrainz_artist_id
		from metadata m
		join user_music_folders f on f.folder_id = m.music_folder_id
		join track_genres g on g.file_path = m.file_path
		LEFT JOIN starred s ON s.metadata_id = m.musicbrainz_track_id
		LEFT JOIN user_ratings ur ON ur.metadata_id = m.musicbrainz_track_id AND ur.user_id = f.user_id
		LEFT JOIN gr ON gr.metadata_id = m.musicbrainz_track_id
		LEFT JOIN track_plays pc ON pc.musicbrainz_track_id = m.musicbrainz_track_id
		left join album_artist_map maa on maa.artist = m.album_artist
		where f.user_id = ?
		and lower(g.genre) = lower(?)`

	args = append(args, requestUser.Id, requestUser.Id, requestUser.Id, genre)

	if musicFolderInt != 0 {
		query += ` and m.music_folder_id = ?`
		args = append(args, musicFolderInt)
	}

	query += ` order by m.musicbrainz_track_id limit ? offset ?`
	args = append(args, count, offset)

	rows, err := DbRead.QueryContext(ctx, query, args...)
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

		result.IsDir = false
		result.MediaType = "song"
		result.Type = "music"
		result.IsVideo = false
		result.Bpm = 0
		result.Comment = ""
		result.Contributors = []types.ChildContributors{}
		result.Moods = []string{}

		if err := rows.Scan(&result.Id, &result.AlbumId, &result.Title, &result.Album, &result.Artist,
			&result.Track, &result.Year, &result.Genre, &result.CoverArt, &result.Size,
			&durationFloat, &result.BitRate, &result.Path, &result.Created, &result.DiscNumber,
			&result.ArtistId, &genreString, &albumArtistName, &result.BitDepth, &result.SamplingRate,
			&result.ChannelCount, &result.UserRating, &result.AverageRating, &result.PlayCount,
			&played, &starred, &albumArtistId); err != nil {
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
		result.Parent = result.AlbumId
		result.SortName = strings.ToLower(result.Title)
		result.MusicBrainzId = result.Id

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

		songs = append(songs, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating song rows: %v", err)
	}

	return songs, nil
}
