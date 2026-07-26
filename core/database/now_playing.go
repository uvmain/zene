package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
	"zene/core/logger"
	"zene/core/types"
)

func createNowPlayingTable(ctx context.Context) {
	schema := `CREATE TABLE now_playing (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		track_id TEXT,
		played_at INTEGER,
		player_id INTEGER,
		player_name TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE (user_id, track_id, player_id, player_name)
	);`
	createTable(ctx, schema)
}

func UpsertNowPlaying(ctx context.Context, userId int, trackId string, playedAt int, playerId int, playerName string) (int, error) {
	query := `INSERT INTO now_playing (user_id, track_id, played_at, player_id, player_name)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(user_id, track_id, player_id, player_name) DO UPDATE SET player_name=excluded.player_name`

	_, err := DB.ExecContext(ctx, query, userId, trackId, playedAt, playerId, playerName)
	if err != nil {
		return 0, fmt.Errorf("upserting now playing row: %v", err)
	}

	selectQuery := `SELECT id FROM now_playing WHERE user_id = ? AND track_id = ? AND player_id = ? AND player_name = ?`
	var nowPlayingId int
	err = DB.QueryRowContext(ctx, selectQuery, userId, trackId, playerId, playerName).Scan(&nowPlayingId)
	if err != nil {
		return 0, fmt.Errorf("getting now playing row id: %v", err)
	}
	return nowPlayingId, nil
}

func MediaIsCurrentlyPlaying(ctx context.Context, userId int, trackId string, playerName string) (bool, error) {
	query := `SELECT 1 FROM now_playing WHERE user_id = ? AND track_id = ? AND player_name = ? LIMIT 1`
	var exists int
	err := DB.QueryRowContext(ctx, query, userId, trackId, playerName).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("checking now playing session: %v", err)
	}
	return true, nil
}

func CleanupNowPlaying(ctx context.Context) error {
	query := `DELETE FROM now_playing
		WHERE id IN (
			SELECT np.id
			FROM now_playing np
			LEFT JOIN playback_reports pr ON pr.now_playing_id = np.id
			WHERE COALESCE(pr.updated_at, np.played_at) < ?
		)`
	tenMinutesAgo := time.Now().Add(-10 * time.Minute).UnixMilli()
	_, err := DB.ExecContext(ctx, query, tenMinutesAgo)
	if err != nil {
		return fmt.Errorf("cleaning up now playing: %v", err)
	}
	return nil
}

func GetNowPlaying(ctx context.Context) ([]types.SubsonicNowPlayingEntry, error) {
	query := `select COALESCE(m.musicbrainz_track_id, pe.guid) as id,
		COALESCE(m.musicbrainz_album_id, CAST(pe.channel_id AS TEXT)) as parent,
		COALESCE(m.title, pe.title, '') as title,
		COALESCE(m.album, pe.album, '') as album,
		COALESCE(m.artist, pe.artist, '') as artist,
		COALESCE(m.track_number, 0) as track,
		CAST(COALESCE(substr(m.release_date,1,4), pe.year, '0') AS INTEGER) as year,
		COALESCE(substr(m.genre,1,(instr(m.genre,';')-1)), podcast_channels.categories, '') as genre,
		COALESCE(m.musicbrainz_track_id, pe.cover_art, pe.guid, '') as cover_art,
		CAST(COALESCE(m.size, pe.size, '0') AS INTEGER) as size,
		CAST(COALESCE(m.duration, pe.duration, 0) AS REAL) as duration,
		CAST(COALESCE(m.bitrate, pe.bit_rate, '0') AS INTEGER) as bitrate,
		COALESCE(m.file_path, pe.file_path, '') as path,
		COALESCE(m.date_added, pe.created_at, '') as created,
		COALESCE(m.disc_number, 0) as disc_number,
		COALESCE(m.musicbrainz_artist_id, CAST(pe.channel_id AS TEXT), '') as artist_id,
		COALESCE(m.genre, podcast_channels.categories, '') as genre_string,
		COALESCE(m.album_artist, podcast_channels.title, '') as album_artist,
		COALESCE(m.bit_depth, 0) as bit_depth,
		COALESCE(m.sample_rate, 0) as sample_rate,
		COALESCE(m.channels, 0) as channels,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(AVG(gr.rating), 0.0) AS average_rating,
		COALESCE(SUM(play_counts.play_count), 0) AS play_count,
		max(play_counts.last_played) as played,
		us.created_at AS starred,
		u.username,
		np.played_at,
		COALESCE(pr.updated_at, np.played_at) as report_updated_at,
		np.player_id,
		np.player_name,
		COALESCE(pr.state, '') as state,
		COALESCE(pr.position_ms, 0) as position_ms,
		CAST(COALESCE(pr.playback_rate, '1.0') AS REAL) as playback_rate
	from now_playing np
	join users u on np.user_id = u.id
	left join metadata m on m.musicbrainz_track_id = np.track_id
	left join podcast_episodes pe on pe.guid = np.track_id
	left join podcast_channels on podcast_channels.id = pe.channel_id
	LEFT JOIN user_stars us ON m.musicbrainz_track_id = us.metadata_id AND us.user_id = np.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_track_id = ur.metadata_id AND ur.user_id = np.user_id
	LEFT JOIN user_ratings gr ON m.musicbrainz_track_id = gr.metadata_id
	LEFT JOIN play_counts play_counts ON m.musicbrainz_track_id = play_counts.musicbrainz_track_id AND play_counts.user_id = np.user_id
	LEFT JOIN playback_reports pr ON np.id = pr.now_playing_id
	group by np.id;`

	rows, err := DB.Query(query)
	if err != nil {
		logger.Printf("Error querying now playing: %v", err)
		return []types.SubsonicNowPlayingEntry{}, err
	}
	defer rows.Close()

	var nowPlaying []types.SubsonicNowPlayingEntry

	for rows.Next() {
		var nowPlayingEntry types.SubsonicNowPlayingEntry
		var albumArtist string
		var genreString string
		var durationFloat float64
		var played sql.NullString
		var starred sql.NullString
		var playedAt int64
		var reportUpdatedAt int64
		var positionMs int

		if err := rows.Scan(&nowPlayingEntry.Id, &nowPlayingEntry.Parent, &nowPlayingEntry.Title, &nowPlayingEntry.Album, &nowPlayingEntry.Artist, &nowPlayingEntry.Track,
			&nowPlayingEntry.Year, &nowPlayingEntry.Genre, &nowPlayingEntry.CoverArt,
			&nowPlayingEntry.Size, &durationFloat, &nowPlayingEntry.BitRate, &nowPlayingEntry.Path, &nowPlayingEntry.Created, &nowPlayingEntry.DiscNumber, &nowPlayingEntry.ArtistId,
			&genreString, &albumArtist, &nowPlayingEntry.BitDepth, &nowPlayingEntry.SamplingRate, &nowPlayingEntry.ChannelCount,
			&nowPlayingEntry.UserRating, &nowPlayingEntry.AverageRating, &nowPlayingEntry.PlayCount, &played, &starred,
			&nowPlayingEntry.Username, &playedAt, &reportUpdatedAt, &nowPlayingEntry.PlayerId, &nowPlayingEntry.PlayerName,
			&nowPlayingEntry.State, &positionMs, &nowPlayingEntry.PlaybackRate); err != nil {
			return nil, err
		}
		nowPlayingEntry.Genres = []types.ChildGenre{}
		for _, genre := range strings.FieldsFunc(genreString, func(r rune) bool {
			return r == ';' || r == ','
		}) {
			nowPlayingEntry.Genres = append(nowPlayingEntry.Genres, types.ChildGenre{Name: genre})
		}

		if played.Valid {
			nowPlayingEntry.Played = played.String
		}
		if starred.Valid {
			nowPlayingEntry.Starred = starred.String
		}

		nowPlayingEntry.Duration = int(durationFloat)
		nowPlayingEntry.IsDir = false

		playedAtTime := time.UnixMilli(playedAt)
		playedAtMinutesAgo := int(time.Since(playedAtTime).Minutes())
		nowPlayingEntry.MinutesAgo = playedAtMinutesAgo

		nowPlayingEntry.Artists = []types.ChildArtist{}
		nowPlayingEntry.Artists = append(nowPlayingEntry.Artists, types.ChildArtist{Id: nowPlayingEntry.ArtistId, Name: nowPlayingEntry.Artist})

		nowPlayingEntry.DisplayArtist = nowPlayingEntry.Artist

		nowPlayingEntry.AlbumArtists = []types.ChildArtist{}
		nowPlayingEntry.AlbumArtists = append(nowPlayingEntry.AlbumArtists, types.ChildArtist{Id: nowPlayingEntry.ArtistId, Name: albumArtist})

		nowPlayingEntry.DisplayAlbumArtist = albumArtist

		// calculate position in milliseconds based on positionMs, playbackRate, and duration and time elapsed since playedAt
		// as per https://opensubsonic.netlify.app/docs/endpoints/reportplayback/
		if nowPlayingEntry.State == "playing" {
			activityTime := playedAtTime
			if reportUpdatedAt > 0 {
				activityTime = time.UnixMilli(reportUpdatedAt)
			}
			elapsedTime := int(time.Since(activityTime).Milliseconds())
			positionMs = positionMs + int(float64(elapsedTime)*nowPlayingEntry.PlaybackRate)
			if positionMs > nowPlayingEntry.Duration*1000 {
				positionMs = nowPlayingEntry.Duration * 1000
			}
		}
		nowPlayingEntry.PositionMs = positionMs

		nowPlaying = append(nowPlaying, nowPlayingEntry)
	}

	return nowPlaying, nil
}
