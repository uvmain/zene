package database

import (
	"context"
	"database/sql"
	"strings"
	"zene/core/logger"
	"zene/core/types"
)

func createTopSongsTable(ctx context.Context) {
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

func SelectTopSongsForArtistId(ctx context.Context, artistId string) ([]types.SubsonicChild, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return []types.SubsonicChild{}, err
	}

	query := `select m.musicbrainz_track_id as id, m.musicbrainz_album_id as album_id, m.title, m.album, m.artist, COALESCE(m.track_number, 0) as track,
		REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year, substr(m.genre,1,(instr(m.genre,';')-1)) as genre, m.musicbrainz_track_id as cover_art,
		m.size, m.duration, m.bitrate, m.file_path as path, m.date_added as created, m.disc_number, m.musicbrainz_artist_id as artist_id,
		m.genre, m.album_artist, m.bit_depth, m.sample_rate, m.channels,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(AVG(gr.rating), 0.0) AS average_rating,
		COALESCE(SUM(pc.play_count), 0) AS play_count,
		max(pc.last_played) as played,
		us.created_at AS starred
	from metadata m
	join top_songs t on t.musicbrainz_album_id = m.musicbrainz_album_id and t.musicbrainz_track_id  = m.musicbrainz_track_id and t.musicbrainz_artist_id = m.musicbrainz_artist_id 
	join user_music_folders f on f.folder_id = m.music_folder_id
	LEFT JOIN user_stars s ON m.musicbrainz_album_id = s.metadata_id AND s.user_id = f.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_album_id = ur.metadata_id AND ur.user_id = f.user_id
	LEFT JOIN user_ratings gr ON m.musicbrainz_album_id = gr.metadata_id
	LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = f.user_id
	LEFT JOIN user_stars us ON m.musicbrainz_track_id = us.metadata_id AND us.user_id = f.user_id
	where f.user_id = ?
	and t.musicbrainz_artist_id = ?
	group by m.musicbrainz_track_id
	order by t.sort_order asc;`

	rows, err := DB.QueryContext(ctx, query, user.Id, artistId)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.SubsonicChild{}, err
	}
	defer rows.Close()

	var results []types.SubsonicChild

	for rows.Next() {
		var result types.SubsonicChild

		var albumArtist string
		var genreString string
		var durationFloat float64
		var played sql.NullString
		var starred sql.NullString

		if err := rows.Scan(&result.Id, &result.Parent, &result.Title, &result.Album, &result.Artist, &result.Track,
			&result.Year, &result.Genre, &result.CoverArt,
			&result.Size, &durationFloat, &result.BitRate, &result.Path, &result.Created, &result.DiscNumber, &result.ArtistId,
			&genreString, &albumArtist, &result.BitDepth, &result.SamplingRate, &result.ChannelCount,
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

		result.Artists = []types.ChildArtist{}
		result.Artists = append(result.Artists, types.ChildArtist{Id: result.ArtistId, Name: result.Artist})

		result.DisplayArtist = result.Artist

		result.AlbumArtists = []types.ChildArtist{}
		result.AlbumArtists = append(result.AlbumArtists, types.ChildArtist{Id: result.ArtistId, Name: albumArtist})

		result.DisplayAlbumArtist = albumArtist

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return results, err
	}

	return results, nil
}

func SelectTopSongsForArtistName(ctx context.Context, artistName string, limit int) ([]types.SubsonicChild, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return []types.SubsonicChild{}, err
	}

	var args []interface{}

	query := `WITH ranked AS (`

	// priority 1: songs from top_songs
	query += `
	SELECT m.musicbrainz_track_id as id, m.musicbrainz_album_id as album_id, m.title as title,
		m.album as album, m.artist as artist, COALESCE(m.track_number, 0) as track,
		REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year, 
		substr(m.genre,1,(instr(m.genre,';')-1)) as genre, m.musicbrainz_track_id as cover_art,
		m.size, m.duration, m.bitrate, m.file_path as path, m.date_added as created,
		m.disc_number as disc_number, m.musicbrainz_artist_id as artist_id, m.genre as genre_string,
		m.album_artist as album_artist, m.bit_depth as bit_depth, m.sample_rate as sample_rate,
		m.channels as channels, COALESCE(ur.rating, 0) AS user_rating, COALESCE(AVG(gr.rating), 0.0) AS average_rating,
		COALESCE(SUM(pc.play_count), 0) AS play_count, max(pc.last_played) as played, us.created_at AS starred,
		1 as priority, t.sort_order as sort_key
    FROM metadata m
    JOIN top_songs t 
			ON t.musicbrainz_album_id = m.musicbrainz_album_id
			AND t.musicbrainz_track_id  = m.musicbrainz_track_id
			AND t.musicbrainz_artist_id = m.musicbrainz_artist_id 
    JOIN user_music_folders f ON f.folder_id = m.music_folder_id
    LEFT JOIN user_stars s ON m.musicbrainz_album_id = s.metadata_id AND s.user_id = f.user_id
    LEFT JOIN user_ratings ur ON m.musicbrainz_album_id = ur.metadata_id AND ur.user_id = f.user_id
    LEFT JOIN user_ratings gr ON m.musicbrainz_album_id = gr.metadata_id
    LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = f.user_id
    LEFT JOIN user_stars us ON m.musicbrainz_track_id = us.metadata_id AND us.user_id = f.user_id
    WHERE f.user_id = ?
      AND lower(m.artist) = lower(?)
    GROUP BY m.musicbrainz_track_id`
	args = append(args, user.Id, artistName)

	// priority 2: songs with play_counts
	query += `
	UNION
		SELECT m.musicbrainz_track_id, m.musicbrainz_album_id, m.title, m.album, m.artist, COALESCE(m.track_number, 0),
			REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0'), substr(m.genre,1,(instr(m.genre,';')-1)),
			m.musicbrainz_track_id, m.size, m.duration, m.bitrate, m.file_path, m.date_added, m.disc_number,
			m.musicbrainz_artist_id, m.genre, m.album_artist, m.bit_depth, m.sample_rate, m.channels, COALESCE(ur.rating, 0),
			COALESCE(AVG(gr.rating), 0.0), COALESCE(SUM(pc.play_count), 0), max(pc.last_played), us.created_at,
			2 as priority, COALESCE(SUM(pc.play_count), 0) as sort_key
    FROM metadata m
    JOIN user_music_folders f ON f.folder_id = m.music_folder_id
    JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = f.user_id
    LEFT JOIN user_stars s ON m.musicbrainz_album_id = s.metadata_id AND s.user_id = f.user_id
    LEFT JOIN user_ratings ur ON m.musicbrainz_album_id = ur.metadata_id AND ur.user_id = f.user_id
    LEFT JOIN user_ratings gr ON m.musicbrainz_album_id = gr.metadata_id
    LEFT JOIN user_stars us ON m.musicbrainz_track_id = us.metadata_id AND us.user_id = f.user_id
    WHERE f.user_id = ?
      AND lower(m.artist) = lower(?)
    GROUP BY m.musicbrainz_track_id`
	args = append(args, user.Id, artistName)

	// priority 3: fallback songs by release_date desc
	query += `
	UNION
    SELECT m.musicbrainz_track_id, m.musicbrainz_album_id, m.title, m.album, m.artist, COALESCE(m.track_number, 0),
			REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0'), substr(m.genre,1,(instr(m.genre,';')-1)), m.musicbrainz_track_id,
			m.size, m.duration, m.bitrate, m.file_path, m.date_added, m.disc_number, m.musicbrainz_artist_id, m.genre, m.album_artist,
			m.bit_depth, m.sample_rate, m.channels, COALESCE(ur.rating, 0), COALESCE(AVG(gr.rating), 0.0),COALESCE(SUM(pc.play_count), 0),
			max(pc.last_played), us.created_at,
			3 as priority, random() as sort_key
    FROM metadata m
    JOIN user_music_folders f ON f.folder_id = m.music_folder_id
    LEFT JOIN user_stars s ON m.musicbrainz_album_id = s.metadata_id AND s.user_id = f.user_id
    LEFT JOIN user_ratings ur ON m.musicbrainz_album_id = ur.metadata_id AND ur.user_id = f.user_id
    LEFT JOIN user_ratings gr ON m.musicbrainz_album_id = gr.metadata_id
    LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = f.user_id
    LEFT JOIN user_stars us ON m.musicbrainz_track_id = us.metadata_id AND us.user_id = f.user_id
    WHERE f.user_id = ?
			AND lower(m.artist) = lower(?)
    GROUP BY m.musicbrainz_track_id`
	args = append(args, user.Id, artistName)

	// now stick the 3 queries together
	query += `
	)
		SELECT id, album_id, title, album, artist, track, year, genre, cover_art, size, duration, bitrate, path, 
		created, disc_number, artist_id, genre_string, album_artist, bit_depth, sample_rate, channels,
		user_rating, average_rating, play_count, played, starred
		FROM ranked
		GROUP BY id
		ORDER BY priority, sort_key ASC
		limit ?;`
	args = append(args, limit)

	rows, err := DB.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.SubsonicChild{}, err
	}
	defer rows.Close()

	var results []types.SubsonicChild

	for rows.Next() {
		var result types.SubsonicChild

		var albumArtist string
		var genreString string
		var durationFloat float64
		var played sql.NullString
		var starred sql.NullString

		if err := rows.Scan(&result.Id, &result.Parent, &result.Title, &result.Album, &result.Artist, &result.Track,
			&result.Year, &result.Genre, &result.CoverArt,
			&result.Size, &durationFloat, &result.BitRate, &result.Path, &result.Created, &result.DiscNumber, &result.ArtistId,
			&genreString, &albumArtist, &result.BitDepth, &result.SamplingRate, &result.ChannelCount,
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

		result.Artists = []types.ChildArtist{}
		result.Artists = append(result.Artists, types.ChildArtist{Id: result.ArtistId, Name: result.Artist})

		result.DisplayArtist = result.Artist

		result.AlbumArtists = []types.ChildArtist{}
		result.AlbumArtists = append(result.AlbumArtists, types.ChildArtist{Id: result.ArtistId, Name: albumArtist})

		result.DisplayAlbumArtist = albumArtist

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

	query := `INSERT INTO top_songs (musicbrainz_album_id, musicbrainz_track_id, musicbrainz_artist_id, sort_order)
		SELECT musicbrainz_album_id, musicbrainz_track_id, musicbrainz_artist_id, ?
		FROM metadata
		WHERE lower(title) = lower(?)
			AND lower(artist) = lower(?)
			AND lower(album) = lower(?)
		ON CONFLICT(musicbrainz_artist_id, musicbrainz_album_id, musicbrainz_track_id) DO NOTHING`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, song := range topSongs {
		_, err = stmt.ExecContext(ctx,
			song.SortOrder,
			song.TrackName,
			song.ArtistName,
			song.AlbumName,
		)
		if err != nil {
			logger.Printf("Failed to insert top song: %v", err)
			return err
		}
	}

	return nil
}
