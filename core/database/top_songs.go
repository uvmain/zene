package database

import (
	"context"
	"database/sql"
	"regexp"
	"strings"
	"zene/core/logger"
	"zene/core/types"
)

func migrateTopSongs(ctx context.Context) {
	schema := `CREATE TABLE top_songs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		musicbrainz_track_id TEXT NOT NULL,
		musicbrainz_album_id TEXT NOT NULL,
		musicbrainz_artist_id TEXT NOT NULL,
		sort_order INTEGER NOT NULL,
		UNIQUE (musicbrainz_artist_id, musicbrainz_album_id, musicbrainz_track_id)
	);`

	createTable(ctx, schema)
	createIndex(ctx, "idx_top_songs_artist", "top_songs", []string{"musicbrainz_artist_id", "sort_order"}, false)
}

func SelectTopSongsForArtistName(ctx context.Context, artistName string, limit int, offset int) ([]types.SubsonicChild, error) {
	user, err := GetUserByContext(ctx)
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
		album_artist_map AS (
			SELECT artist,
				MIN(musicbrainz_artist_id) AS musicbrainz_artist_id
			FROM metadata
			WHERE musicbrainz_artist_id IS NOT NULL
		GROUP BY artist
		),
		starred as (
			select metadata_id,
				created_at
			from user_stars
			where user_id = ?
		)
		SELECT m.musicbrainz_track_id,
			m.musicbrainz_album_id,
			m.title,
			m.album,
			m.artist,
			COALESCE(m.track_number, 0) as track_number,
			REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year,
			substr(m.genre,1,(instr(m.genre,';')-1)) as genre,
			m.musicbrainz_track_id,
			m.size,
			m.duration,
			m.bitrate,
			m.file_path,
			m.date_added,
			m.disc_number,
			m.musicbrainz_artist_id,
			m.genre,
			m.album_artist,
			maa.musicbrainz_artist_id as album_artist_id,
			m.bit_depth,
			m.sample_rate,
			m.channels,
			COALESCE(ur.rating, 0),
			COALESCE(AVG(gr.rating), 0.0),
			COALESCE(tp.play_count, 0) as play_count,
			tp.last_played,
			us.created_at
		FROM metadata m
		JOIN user_music_folders f ON f.folder_id = m.music_folder_id
		LEFT JOIN user_ratings ur ON m.musicbrainz_track_id = ur.metadata_id AND ur.user_id = f.user_id
		LEFT JOIN user_ratings gr ON m.musicbrainz_track_id = gr.metadata_id
		LEFT JOIN starred us ON us.metadata_id = m.musicbrainz_track_id
		LEFT JOIN album_artist_map maa ON maa.artist = m.album_artist
		left join track_plays tp on tp.musicbrainz_track_id = m.musicbrainz_track_id
		left join top_songs ts on ts.musicbrainz_track_id = m.musicbrainz_track_id
		WHERE f.user_id = ? and lower(m.artist) = lower(?)
		GROUP BY m.musicbrainz_track_id
		order by coalesce(ts.sort_order, 1) desc, us.created_at desc, play_count desc, release_date desc
		limit ? offset ?`
	args = append(args, user.Id, user.Id, user.Id, strings.ToLower(artistName), limit, offset)

	rows, err := DB.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.SubsonicChild{}, err
	}
	defer rows.Close()

	var results []types.SubsonicChild

	for rows.Next() {
		var result types.SubsonicChild

		var albumArtistName sql.NullString
		var albumArtistId sql.NullString
		var genreString string
		var durationFloat float64
		var played sql.NullString
		var starred sql.NullString

		if err := rows.Scan(&result.Id, &result.Parent, &result.Title, &result.Album, &result.Artist, &result.Track,
			&result.Year, &result.Genre, &result.CoverArt, &result.Size, &durationFloat, &result.BitRate, &result.Path,
			&result.Created, &result.DiscNumber, &result.ArtistId,
			&genreString, &albumArtistName, &albumArtistId, &result.BitDepth, &result.SamplingRate, &result.ChannelCount,
			&result.UserRating, &result.AverageRating, &result.PlayCount, &played, &starred); err != nil {
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

func InsertTopSongs(ctx context.Context, topSongs []types.TopSongRow) error {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	for _, song := range topSongs {
		// Deezer sometimes only has special editions or re-releases of an album, eg "Superunknown (Super Deluxe)"
		// so we remove any parenthesis suffixes when matching album names
		cleanAlbumName := regexp.MustCompile(`([ ]\(.*\))?$`).ReplaceAllString(song.AlbumName, "")
		query := `INSERT INTO top_songs (musicbrainz_album_id, musicbrainz_track_id, musicbrainz_artist_id, sort_order)
			SELECT musicbrainz_album_id, musicbrainz_track_id, musicbrainz_artist_id, ?
			FROM metadata
			WHERE lower(title) = ?
				AND lower(artist) = ?
				AND lower(album) = ?
			LIMIT 1
			ON CONFLICT(musicbrainz_artist_id, musicbrainz_album_id, musicbrainz_track_id) DO NOTHING`

		_, err = tx.ExecContext(ctx, query, song.SortOrder, strings.ToLower(song.TrackName), strings.ToLower(song.ArtistName), strings.ToLower(cleanAlbumName))
		if err != nil {
			logger.Printf("Failed to insert top song: %v", err)
			return err
		}
	}

	return nil
}
