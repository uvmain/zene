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

func GetSimilarSongs(ctx context.Context, count int, musicbrainzId string) ([]types.SubsonicChild, error) {
	valid, metadataRow, err := IsValidMetadataId(ctx, musicbrainzId)
	if err != nil || !valid {
		return []types.SubsonicChild{}, fmt.Errorf("invalid musicbrainz id '%s': %v", musicbrainzId, err)
	}

	var query string
	var artistId string

	if metadataRow.MusicbrainzTrackId {
		query = `select musicbrainz_artist_id from metadata where musicbrainz_track_id = ? limit 1`
	} else if metadataRow.MusicbrainzAlbumId {
		query = `select musicbrainz_artist_id from metadata where musicbrainz_album_id = ? limit 1`
	}

	if !metadataRow.MusicbrainzArtistId {
		err := DB.QueryRow(query, musicbrainzId).Scan(&artistId)
		if err == sql.ErrNoRows {
			return []types.SubsonicChild{}, fmt.Errorf("no artist found for musicbrainz ID: %s", musicbrainzId)
		} else if err != nil {
			return []types.SubsonicChild{}, fmt.Errorf("error querying artist ID: %v", err)
		}
	} else {
		artistId = musicbrainzId
	}

	requestUser, err := GetUserByContext(ctx)
	if err != nil {
		return []types.SubsonicChild{}, err
	}

	var args []interface{}

	query = `select m.musicbrainz_track_id as id, m.musicbrainz_album_id as parent, m.title, m.album, m.artist,
		COALESCE(m.track_number, 0) as track,
		REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year, substr(m.genre,1,(instr(m.genre,';')-1)) as genre, m.musicbrainz_track_id as cover_art,
		m.size, m.duration, m.bitrate, m.file_path as path, m.date_added as created, m.disc_number, m.musicbrainz_artist_id as artist_id,
		m.genre, m.album_artist, m.bit_depth, m.sample_rate, m.channels,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(AVG(gr.rating), 0.0) AS average_rating,
		COALESCE(SUM(pc.play_count), 0) AS play_count,
		max(pc.last_played) as played,
		us.created_at AS starred,
		maa.musicbrainz_artist_id as album_artist_id
	from user_music_folders u
	join metadata m on m.music_folder_id = u.folder_id
	LEFT JOIN user_stars us ON m.musicbrainz_album_id = us.metadata_id AND us.user_id = u.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_album_id = ur.metadata_id AND ur.user_id = u.user_id
	LEFT JOIN user_ratings gr ON m.musicbrainz_album_id = gr.metadata_id
	LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = u.user_id
	left join metadata maa on maa.artist = m.album_artist
	where u.user_id = ?`

	args = append(args, requestUser.Id)

	query += ` and (
		m.musicbrainz_artist_id = ?
		or m.musicbrainz_artist_id in (
			select similar_artist_id
			from similar_artists
			where artist_id = ?
			)
		)
	group by m.musicbrainz_track_id`
	args = append(args, artistId, artistId)

	randomInt := logic.GenerateRandomInt(1, 10000000)
	query += fmt.Sprintf(" order BY ((m.rowid * %d) %% 1000000)", randomInt)

	query += ` limit ?`
	args = append(args, count)

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
		var albumArtistName string
		var albumArtistId string
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
		result.AlbumArtists = append(result.AlbumArtists, types.ChildArtist{Id: albumArtistId, Name: albumArtistName})

		result.DisplayAlbumArtist = albumArtistName

		songs = append(songs, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating song rows: %v", err)
	}

	return songs, nil
}
