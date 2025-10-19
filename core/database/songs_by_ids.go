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

func GetSongsByIDs(ctx context.Context, ids []string) ([]types.SubsonicChild, error) {
	requestUser, err := GetUserByContext(ctx)
	if err != nil {
		return []types.SubsonicChild{}, err
	}

	var args []interface{}
	query := `select m.musicbrainz_track_id as id, m.musicbrainz_album_id as album_id, m.title, m.album, m.artist, COALESCE(m.track_number, 0) as track,
		REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year, substr(m.genre,1,(instr(m.genre,';')-1)) as genre, m.musicbrainz_track_id as cover_art,
		m.size, m.duration, m.bitrate, m.file_path as path, m.date_added as created, m.disc_number, m.musicbrainz_artist_id as artist_id,
		m.genre, m.album_artist, m.bit_depth, m.sample_rate, m.channels,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(AVG(gr.rating), 0.0) AS average_rating,
		COALESCE(SUM(pc.play_count), 0) AS play_count,
		max(pc.last_played) as played,
		us.created_at AS starred,
		maa.musicbrainz_artist_id
	from metadata m
	join user_music_folders f on f.folder_id = m.music_folder_id
	join track_genres g on m.file_path = g.file_path
	LEFT JOIN user_stars s ON m.musicbrainz_track_id = s.metadata_id AND s.user_id = f.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_track_id = ur.metadata_id AND ur.user_id = f.user_id
	LEFT JOIN user_ratings gr ON m.musicbrainz_track_id = gr.metadata_id
	LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = f.user_id
	LEFT JOIN user_stars us ON m.musicbrainz_track_id = us.metadata_id AND us.user_id = f.user_id
	left join metadata maa on maa.artist = m.album_artist
	where f.user_id = ?`

	args = append(args, requestUser.Id)

	query += ` and m.musicbrainz_track_id IN (`
	for i := 0; i < len(ids); i++ {
		if i > 0 {
			query += `,`
		}
		query += `?`
		args = append(args, ids[i])
	}
	query += ` ) group by m.musicbrainz_track_id order by m.musicbrainz_track_id`

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
			logger.Printf("Failed to scan row in GetSongsByIDs: %v", err)
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
