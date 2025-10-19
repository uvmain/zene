package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"zene/core/logic"
	"zene/core/types"
)

func migrateBookmarks(ctx context.Context) {
	schema := `CREATE TABLE bookmarks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		musicbrainz_track_id TEXT NOT NULL,
		created TEXT NOT NULL,
		changed TEXT NOT NULL,
		comment TEXT,
		position INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		UNIQUE (musicbrainz_track_id, user_id)
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_bookmarks_user", "bookmarks", []string{"user_id"}, false)
}

func UpsertBookmark(ctx context.Context, musicbrainzTrackId string, position int, comment string) error {
	created_at := logic.GetCurrentTimeFormatted()

	user, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("getting user from context: %v", err)
	}

	if comment != "" {
		query := `INSERT INTO bookmarks (user_id, musicbrainz_track_id, created, changed, comment, position)
			VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT(musicbrainz_track_id, user_id) DO UPDATE SET changed=excluded.changed, comment=excluded.comment, position=excluded.position`
		_, err = DB.ExecContext(ctx, query, user.Id, musicbrainzTrackId, created_at, created_at, comment, position)
	} else {
		query := `INSERT INTO bookmarks (user_id, musicbrainz_track_id, created, changed, position)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(musicbrainz_track_id, user_id) DO UPDATE SET changed=excluded.changed, position=excluded.position`
		_, err = DB.ExecContext(ctx, query, user.Id, musicbrainzTrackId, created_at, created_at, position)
	}

	if err != nil {
		return fmt.Errorf("inserting chat: %v", err)
	}
	return nil
}

func GetBookmarks(ctx context.Context) ([]types.Bookmark, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return []types.Bookmark{}, fmt.Errorf("getting user from context: %v", err)
	}

	query := `select m.musicbrainz_track_id as id, m.musicbrainz_album_id as parent, m.title, m.album, m.artist,
		COALESCE(m.track_number, 0) as track,
		REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year,
		substr(m.genre,1,(instr(m.genre,';')-1)) as genre, m.musicbrainz_track_id as cover_art,
		m.size, m.label,
		m.duration, m.bitrate, m.file_path as path, m.date_added as created, m.disc_number, m.musicbrainz_artist_id as artist_id,
		m.genre, m.album_artist, m.bit_depth, m.sample_rate, m.channels,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(AVG(gr.rating), 0.0) AS average_rating,
		COALESCE(SUM(pc.play_count), 0) AS play_count,
		max(pc.last_played) as played,
		us.created_at AS starred,
		maa.musicbrainz_artist_id as album_artist_id,
		u.username,
		b.created bookmarked_created,
		b.changed bookmark_changed,
		b.comment,
		b.position
	from bookmarks b
	join users u on u.id = b.user_id
	join user_music_folders uf on uf.user_id = b.user_id
	join metadata m on m.music_folder_id = uf.folder_id and b.musicbrainz_track_id = m.musicbrainz_track_id
	LEFT JOIN user_stars us ON m.musicbrainz_track_id = us.metadata_id AND us.user_id = u.id
	LEFT JOIN user_ratings ur ON m.musicbrainz_track_id = ur.metadata_id AND ur.user_id = u.id
	LEFT JOIN user_ratings gr ON m.musicbrainz_track_id = gr.metadata_id
	LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = u.id
	left join metadata maa on maa.artist = m.album_artist
	where b.user_id = ?
	group by m.musicbrainz_track_id`

	rows, err := DB.QueryContext(ctx, query, user.Id)
	if err != nil {
		return []types.Bookmark{}, fmt.Errorf("querying bookmarks: %v", err)
	}
	defer rows.Close()

	var result []types.Bookmark
	for rows.Next() {
		var row types.Bookmark
		var albumArtistName sql.NullString
		var albumArtistId sql.NullString
		var genreString string
		var labelString string
		var playedString sql.NullString
		var starredString sql.NullString
		var durationFloat float64

		err := rows.Scan(&row.Entry.Id, &row.Entry.Parent, &row.Entry.Title, &row.Entry.Album, &row.Entry.Artist,
			&row.Entry.Track, &row.Entry.Year, &row.Entry.Genre, &row.Entry.CoverArt, &row.Entry.Size, &labelString,
			&durationFloat, &row.Entry.BitRate, &row.Entry.Path, &row.Entry.Created, &row.Entry.DiscNumber, &row.Entry.ArtistId,
			&genreString, &albumArtistName, &row.Entry.BitDepth, &row.Entry.SamplingRate, &row.Entry.ChannelCount,
			&row.Entry.UserRating, &row.Entry.AverageRating, &row.Entry.PlayCount, &playedString,
			&starredString, &albumArtistId,
			&row.Username,
			&row.Created,
			&row.Changed,
			&row.Comment,
			&row.Position,
		)
		if err != nil {
			return []types.Bookmark{}, fmt.Errorf("scanning bookmark row: %v", err)
		}

		row.Entry.Genres = []types.ChildGenre{}
		for _, genre := range strings.Split(genreString, ";") {
			row.Entry.Genres = append(row.Entry.Genres, types.ChildGenre{Name: genre})
		}

		if playedString.Valid {
			row.Entry.Played = playedString.String
		}

		if starredString.Valid {
			row.Entry.Starred = starredString.String
		}

		row.Entry.RecordLabels = []types.ChildRecordLabel{}
		row.Entry.RecordLabels = append(row.Entry.RecordLabels, types.ChildRecordLabel{Name: labelString})

		row.Entry.Duration = int(durationFloat)
		row.Entry.Title = row.Entry.Album
		row.Entry.IsDir = false

		row.Entry.Artists = []types.ChildArtist{}
		row.Entry.Artists = append(row.Entry.Artists, types.ChildArtist{Id: row.Entry.ArtistId, Name: row.Entry.Artist})

		row.Entry.DisplayArtist = row.Entry.Artist

		row.Entry.AlbumArtists = []types.ChildArtist{}
		if albumArtistId.Valid && albumArtistName.Valid {
			row.Entry.AlbumArtists = append(row.Entry.AlbumArtists, types.ChildArtist{Id: albumArtistId.String, Name: albumArtistName.String})
		}

		row.Entry.DisplayAlbumArtist = albumArtistName.String

		result = append(result, row)
	}
	return result, nil
}

func DeleteBookmark(ctx context.Context, musicbrainzTrackId string) error {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("getting user from context: %v", err)
	}

	query := `DELETE FROM bookmarks WHERE musicbrainz_track_id = ? AND user_id = ?`
	_, err = DB.ExecContext(ctx, query, musicbrainzTrackId, user.Id)
	if err != nil {
		return fmt.Errorf("deleting bookmark: %v", err)
	}
	return nil
}
