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

func migrateNowPlaying(ctx context.Context) {
	schema := `CREATE TABLE now_playing (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		track_id TEXT,
		played_at INTEGER,
		player_id INTEGER,
		player_name TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE (user_id, track_id)
	);`
	createTable(ctx, schema)
}

func UpsertNowPlaying(ctx context.Context, userId int, trackId string, playedAt int, playerId int, playerName string) error {
	query := `DELETE FROM now_playing where user_id = ? and player_id = ? and player_name = ?`

	_, err := DbWrite.ExecContext(ctx, query, userId, playerId, playerName)
	if err != nil {
		return fmt.Errorf("removing existing playing row: %v", err)
	}

	query = `INSERT INTO now_playing (user_id, track_id, played_at, player_id, player_name)
		VALUES (?, ?, ?, ?, ?)`

	_, err = DbWrite.ExecContext(ctx, query, userId, trackId, playedAt, playerId, playerName)
	if err != nil {
		return fmt.Errorf("upserting now playing row: %v", err)
	}
	return nil
}

func CleanupNowPlaying(ctx context.Context) error {
	query := `DELETE FROM now_playing WHERE played_at < ?`
	tenMinutesAgo := time.Now().Add(-10 * time.Minute).UnixMilli()
	_, err := DbWrite.ExecContext(ctx, query, tenMinutesAgo)
	if err != nil {
		return fmt.Errorf("cleaning up now playing: %v", err)
	}
	return nil
}

func GetNowPlaying(ctx context.Context) ([]types.SubsonicNowPlayingEntry, error) {
	query := `select m.musicbrainz_track_id as id, m.musicbrainz_album_id as parent, m.title, m.album, m.artist, COALESCE(m.track_number, 0) as track,
		REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year, substr(m.genre,1,(instr(m.genre,';')-1)) as genre, m.musicbrainz_track_id as cover_art,
		m.size, m.duration, m.bitrate, m.file_path as path, m.date_added as created, m.disc_number, m.musicbrainz_artist_id as artist_id,
		m.genre, m.album_artist, m.bit_depth, m.sample_rate, m.channels,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(AVG(gr.rating), 0.0) AS average_rating,
		COALESCE(SUM(pc.play_count), 0) AS play_count,
		max(pc.last_played) as played,
		us.created_at AS starred,
		u.username,
		np.played_at,
		np.player_id,
		np.player_name
	from now_playing np
	join users u on np.user_id = u.id
	join metadata m on m.musicbrainz_track_id = np.track_id
	LEFT JOIN user_stars us ON m.musicbrainz_track_id = us.metadata_id AND us.user_id = np.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_track_id = ur.metadata_id AND ur.user_id = np.user_id
	LEFT JOIN user_ratings gr ON m.musicbrainz_track_id = gr.metadata_id
	LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = np.user_id
	group by m.musicbrainz_track_id;`

	rows, err := DbRead.Query(query)
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

		if err := rows.Scan(&nowPlayingEntry.Id, &nowPlayingEntry.Parent, &nowPlayingEntry.Title, &nowPlayingEntry.Album, &nowPlayingEntry.Artist, &nowPlayingEntry.Track,
			&nowPlayingEntry.Year, &nowPlayingEntry.Genre, &nowPlayingEntry.CoverArt,
			&nowPlayingEntry.Size, &durationFloat, &nowPlayingEntry.BitRate, &nowPlayingEntry.Path, &nowPlayingEntry.Created, &nowPlayingEntry.DiscNumber, &nowPlayingEntry.ArtistId,
			&genreString, &albumArtist, &nowPlayingEntry.BitDepth, &nowPlayingEntry.SamplingRate, &nowPlayingEntry.ChannelCount,
			&nowPlayingEntry.UserRating, &nowPlayingEntry.AverageRating, &nowPlayingEntry.PlayCount, &played, &starred,
			&nowPlayingEntry.Username, &playedAt, &nowPlayingEntry.PlayerId, &nowPlayingEntry.PlayerName); err != nil {
			return nil, err
		}
		nowPlayingEntry.Genres = []types.ChildGenre{}
		for _, genre := range strings.Split(genreString, ";") {
			nowPlayingEntry.Genres = append(nowPlayingEntry.Genres, types.ChildGenre{Name: genre})
		}

		if played.Valid {
			nowPlayingEntry.Played = played.String
		}
		if starred.Valid {
			nowPlayingEntry.Starred = starred.String
		}

		nowPlayingEntry.Duration = int(durationFloat)
		nowPlayingEntry.Title = nowPlayingEntry.Album
		nowPlayingEntry.IsDir = true

		playedAtTime := time.UnixMilli(playedAt)
		playedAtMinutesAgo := int(time.Since(playedAtTime).Minutes())
		nowPlayingEntry.MinutesAgo = playedAtMinutesAgo

		nowPlayingEntry.Artists = []types.ChildArtist{}
		nowPlayingEntry.Artists = append(nowPlayingEntry.Artists, types.ChildArtist{Id: nowPlayingEntry.ArtistId, Name: nowPlayingEntry.Artist})

		nowPlayingEntry.DisplayArtist = nowPlayingEntry.Artist

		nowPlayingEntry.AlbumArtists = []types.ChildArtist{}
		nowPlayingEntry.AlbumArtists = append(nowPlayingEntry.AlbumArtists, types.ChildArtist{Id: nowPlayingEntry.ArtistId, Name: albumArtist})

		nowPlayingEntry.DisplayAlbumArtist = albumArtist
		nowPlaying = append(nowPlaying, nowPlayingEntry)
	}

	return nowPlaying, nil
}
