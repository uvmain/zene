package database

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

type CreateShareOptions struct {
	Description string
	ExpiresAt   time.Time
	MediaIds    []string
}

type UpdateShareOptions struct {
	ShareId           int
	UpdateDescription bool
	Description       string
	UpdateExpiresAt   bool
	ExpiresAt         time.Time
}

type Share struct {
	Id          int
	OwnerUserId int
	Token       string
	Description string
	CreatedAt   string
	ExpiresAt   string
	VisitCount  int
	LastVisited string
	MediaIds    []string
}

func createSharesTable(ctx context.Context) {
	schema := `CREATE TABLE shares (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		owner_user_id INTEGER NOT NULL,
		token TEXT NOT NULL,
		description TEXT,
		created_at TEXT NOT NULL,
		expires_at TEXT,
		visit_count INTEGER NOT NULL DEFAULT 0,
		last_visited TEXT,
		FOREIGN KEY (owner_user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(token)
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_shares_owner_user_id", "shares", []string{"owner_user_id"}, false)
}

func createSharedMediaTable(ctx context.Context) {
	schema := `CREATE TABLE shared_media (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		share_id INTEGER NOT NULL,
		media_id TEXT NOT NULL,
		FOREIGN KEY (share_id) REFERENCES shares(id) ON DELETE CASCADE
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_shared_media_share_id", "shared_media", []string{"share_id"}, false)
	createIndex(ctx, "idx_shared_media_share_id_media_id", "shared_media", []string{"share_id", "media_id"}, false)
}

func CreateShare(ctx context.Context, options CreateShareOptions) (types.ShareRow, error) {
	owner, err := GetUserByContext(ctx)
	if err != nil {
		return types.ShareRow{}, err
	}

	createdDate := logic.GetCurrentTimeFormatted()
	token, err := logic.GenerateNewApiKey()
	if err != nil {
		return types.ShareRow{}, fmt.Errorf("generating new token: %v", err)
	}

	if !options.ExpiresAt.IsZero() && logic.GetTimeFromString(createdDate).After(options.ExpiresAt) {
		return types.ShareRow{}, fmt.Errorf("expires_at must be in the future")
	}

	expiresAtForDb := sql.NullString{}
	if !options.ExpiresAt.IsZero() {
		expiresAtForDb = sql.NullString{
			String: logic.FormatTimeAsString(options.ExpiresAt),
			Valid:  true,
		}
	}

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return types.ShareRow{}, fmt.Errorf("starting transaction: %v", err)
	}

	shareQuery := `INSERT INTO shares (owner_user_id, token, description, created_at, expires_at)
		VALUES (?, ?, ?, ?, ?)`

	result, err := tx.ExecContext(ctx, shareQuery, owner.Id, token, options.Description, createdDate, expiresAtForDb)
	if err != nil {
		tx.Rollback()
		return types.ShareRow{}, fmt.Errorf("inserting share row: %v", err)
	}

	lastInserted, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return types.ShareRow{}, fmt.Errorf("getting last inserted ID: %v", err)
	}

	mediaQuery := `INSERT INTO shared_media (share_id, media_id) VALUES (?, ?)`
	for _, mediaId := range options.MediaIds {
		_, err := tx.ExecContext(ctx, mediaQuery, lastInserted, mediaId)
		if err != nil {
			tx.Rollback()
			return types.ShareRow{}, fmt.Errorf("inserting shared media row: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return types.ShareRow{}, fmt.Errorf("committing transaction: %v", err)
	}

	share, err := GetShareById(ctx, int(lastInserted))
	if err != nil {
		return types.ShareRow{}, fmt.Errorf("retrieving created share: %v", err)
	}

	return share, nil
}

func UpdateShare(ctx context.Context, options UpdateShareOptions) error {
	if !options.UpdateExpiresAt && !options.UpdateDescription {
		return fmt.Errorf("UpdateShare called with neither UpdateExpiresAt nor UpdateDescription set to true")
	}

	user, err := GetUserByContext(ctx)
	if err != nil {
		return err
	}

	query := `SELECT owner_user_id FROM shares WHERE id = ?`
	var ownerUserId int
	err = DB.QueryRowContext(ctx, query, options.ShareId).Scan(&ownerUserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("share with id %d not found", options.ShareId)
		}
		return fmt.Errorf("querying share owner: %v", err)
	}

	if ownerUserId != user.Id {
		return fmt.Errorf("user does not own the share")
	}

	if options.UpdateDescription {
		updateQuery := `UPDATE shares SET description = ? WHERE id = ?`
		_, err := DB.ExecContext(ctx, updateQuery, options.Description, options.ShareId)
		if err != nil {
			return fmt.Errorf("updating share description: %v", err)
		}
	}

	if options.UpdateExpiresAt {
		expiresAtForDb := sql.NullString{}
		if !options.ExpiresAt.IsZero() {
			expiresAtForDb = sql.NullString{
				String: logic.FormatTimeAsString(options.ExpiresAt),
				Valid:  true,
			}
		}

		updateQuery := `UPDATE shares SET expires_at = ? WHERE id = ?`
		_, err := DB.ExecContext(ctx, updateQuery, expiresAtForDb, options.ShareId)
		if err != nil {
			return fmt.Errorf("updating share expires_at: %v", err)
		}
	}

	return nil
}

func IncrementShareVisitCount(ctx context.Context, shareId int) error {
	lastVisited := logic.GetCurrentTimeFormatted()
	updateQuery := `UPDATE shares SET visit_count = visit_count + 1,
		last_visited = ?
		WHERE id = ?`
	_, err := DB.ExecContext(ctx, updateQuery, lastVisited, shareId)
	if err != nil {
		return fmt.Errorf("incrementing share visit count: %v", err)
	}

	return nil
}

func DeleteShare(ctx context.Context, id int) error {
	owner, err := GetUserByContext(ctx)
	if err != nil {
		return err
	}

	query := `DELETE FROM shares WHERE id = ? AND owner_user_id = ?`
	result, err := DB.ExecContext(ctx, query, id, owner.Id)
	if err != nil {
		return fmt.Errorf("deleting share: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no share found with id %d for the current user", id)
	}

	return nil
}

func GetAllShares(ctx context.Context) ([]types.ShareRow, error) {
	query := `WITH gr AS (
    SELECT metadata_id, AVG(rating) AS avg_rating
    FROM user_ratings
    GROUP BY metadata_id
		),
		plays AS (
				SELECT musicbrainz_track_id, SUM(play_count) AS play_count, MAX(last_played) AS last_played
				FROM play_counts
				GROUP BY musicbrainz_track_id
		),
		share_media AS (
				SELECT sh.id AS share_id, sh.owner_user_id, sh.token, sh.description, sh.created_at, sh.expires_at, sh.visit_count, sh.last_visited, sm.media_id
				FROM shares sh
				JOIN shared_media sm ON sm.share_id = sh.id
		),
		shared_metadata AS (
				SELECT sm.share_id, sm.owner_user_id, sm.token, sm.description, sm.created_at, sm.expires_at, sm.visit_count, sm.last_visited, m.*
				FROM share_media sm
				JOIN metadata m ON m.musicbrainz_track_id = sm.media_id
				UNION
				SELECT sm.share_id, sm.owner_user_id, sm.token, sm.description, sm.created_at, sm.expires_at, sm.visit_count, sm.last_visited, m.*
				FROM share_media sm
				JOIN metadata m ON m.musicbrainz_album_id = sm.media_id
				UNION
				SELECT sm.share_id, sm.owner_user_id, sm.token, sm.description, sm.created_at, sm.expires_at, sm.visit_count, sm.last_visited, m.*
				FROM share_media sm
				JOIN metadata m ON m.musicbrainz_artist_id = sm.media_id
		)
		SELECT
				sm.share_id AS id, sm.token, sm.description, sm.created_at, sm.expires_at, sm.visit_count, sm.last_visited, u.username,
				m.musicbrainz_track_id AS id, m.musicbrainz_album_id AS album_id, m.title, m.album, m.artist,
				COALESCE(m.track_number, 0) AS track, REPLACE(PRINTF('%4s', substr(m.release_date, 1, 4)),' ','0') AS year,
				substr(m.genre,1,instr(m.genre, ';') - 1) AS genre, m.musicbrainz_track_id AS cover_art, m.size,
				m.duration, m.bitrate, m.file_path AS path, m.date_added AS created, m.disc_number, m.musicbrainz_artist_id AS artist_id,
				m.album_artist, m.bit_depth, m.sample_rate, m.channels, COALESCE(ur.rating, 0) AS user_rating,
				COALESCE(gr.avg_rating, 0.0) AS average_rating, COALESCE(plays.play_count, 0) AS play_count, plays.last_played AS played,
			us.created_at AS starred, maa.musicbrainz_artist_id, maa.artist AS album_artist_name
		FROM shared_metadata sm
		JOIN metadata m ON m.musicbrainz_track_id = sm.musicbrainz_track_id
		JOIN user_music_folders f ON f.folder_id = m.music_folder_id AND f.user_id = sm.owner_user_id
		JOIN users u ON u.id = sm.owner_user_id
		LEFT JOIN user_stars us ON us.metadata_id = m.musicbrainz_track_id AND us.user_id = sm.owner_user_id
		LEFT JOIN user_ratings ur ON ur.metadata_id = m.musicbrainz_track_id AND ur.user_id = sm.owner_user_id
		LEFT JOIN gr ON gr.metadata_id = m.musicbrainz_artist_id
		LEFT JOIN plays ON plays.musicbrainz_track_id = m.musicbrainz_track_id
		LEFT JOIN metadata maa ON maa.artist = m.album_artist
		group by sm.share_id, m.musicbrainz_track_id
		ORDER BY sm.created_at DESC, m.musicbrainz_artist_id, m.musicbrainz_album_id, m.musicbrainz_track_id`

	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying shares: %v", err)
	}
	defer rows.Close()

	var shares []types.ShareRow

	for rows.Next() {
		var row types.DbUserShare
		var lastPlayed sql.NullString
		var dateStarred sql.NullString
		var durationFloat float64
		var albumArtistName sql.NullString
		var albumArtistId sql.NullString
		var shareExpiresAt sql.NullString
		var shareLastVisited sql.NullString
		if err := rows.Scan(
			&row.ShareId, &row.Token, &row.Description, &row.ShareCreated, &shareExpiresAt, &row.VisitCount, &shareLastVisited,
			&row.ShareOwner, &row.TrackId, &row.AlbumId, &row.Title, &row.Album, &row.Artist, &row.TrackNumber, &row.Year,
			&row.Genre, &row.CoverArt, &row.Size, &durationFloat, &row.Bitrate, &row.Path, &row.DateAdded,
			&row.DiscNumber, &row.ArtistId, &row.AlbumArtist, &row.BitDepth, &row.SampleRate, &row.Channels,
			&row.UserRating, &row.AverageRating, &row.PlayCount, &lastPlayed, &dateStarred, &albumArtistId, &albumArtistName,
		); err != nil {
			return nil, fmt.Errorf("scanning share row: %v", err)
		}

		var entry types.SubsonicChild
		entry.Genres = []types.ChildGenre{}
		for _, genre := range strings.Split(row.Genre, ";") {
			entry.Genres = append(entry.Genres, types.ChildGenre{Name: genre})
		}

		if lastPlayed.Valid {
			entry.Played = lastPlayed.String
		}

		if dateStarred.Valid {
			entry.Starred = dateStarred.String
		}

		if dateStarred.Valid {
			entry.Starred = dateStarred.String
		}

		entry.Id = row.TrackId
		entry.Duration = int(durationFloat)
		entry.IsDir = false
		entry.MusicBrainzId = entry.Id
		entry.AlbumId = entry.Parent
		entry.Title = row.Title
		entry.Album = row.Album
		entry.Artist = row.Artist
		entry.Track = row.TrackNumber
		entry.Year = row.Year
		entry.Genre = row.Genre
		entry.CoverArt = row.CoverArt
		entry.Size = row.Size
		entry.BitRate = row.Bitrate
		entry.Path = row.Path
		entry.Created = row.DateAdded
		entry.DiscNumber = row.DiscNumber
		entry.ArtistId = row.ArtistId
		entry.BitDepth = row.BitDepth
		entry.SamplingRate = row.SampleRate
		entry.ChannelCount = row.Channels
		entry.UserRating = row.UserRating
		entry.AverageRating = row.AverageRating
		entry.PlayCount = row.PlayCount

		entry.Artists = []types.ChildArtist{}
		entry.Artists = append(entry.Artists, types.ChildArtist{Id: entry.ArtistId, Name: entry.Artist})

		entry.DisplayArtist = entry.Artist

		entry.AlbumArtists = []types.ChildArtist{}
		if albumArtistId.Valid && albumArtistName.Valid {
			entry.AlbumArtists = append(entry.AlbumArtists, types.ChildArtist{Id: albumArtistId.String, Name: albumArtistName.String})
		}

		entry.DisplayAlbumArtist = albumArtistName.String

		currentShare := slices.IndexFunc(shares, func(s types.ShareRow) bool { return s.Id == row.ShareId })

		var expiresAt string
		if shareExpiresAt.Valid {
			expiresAt = shareExpiresAt.String
		}

		var lastVisited string
		if shareLastVisited.Valid {
			lastVisited = shareLastVisited.String
		}

		if currentShare == -1 {
			shares = append(shares, types.ShareRow{
				Id:          row.ShareId,
				Description: row.Description,
				Username:    row.ShareOwner,
				Url:         logic.GetShareUrl(row.Token),
				Created:     row.ShareCreated,
				VisitCount:  row.VisitCount,
				LastVisited: lastVisited,
				Expires:     expiresAt,
				Entries:     []types.SubsonicChild{},
			})
			currentShare = len(shares) - 1
		}

		shares[currentShare].Entries = append(shares[currentShare].Entries, entry)
	}

	return shares, nil
}

func GetSharesByUser(ctx context.Context) ([]types.ShareRow, error) {
	owner, err := GetUserByContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `WITH gr AS (
    SELECT metadata_id, AVG(rating) AS avg_rating
    FROM user_ratings
    GROUP BY metadata_id
		),
		plays AS (
				SELECT musicbrainz_track_id, SUM(play_count) AS play_count, MAX(last_played) AS last_played
				FROM play_counts
				WHERE user_id = ?
				GROUP BY musicbrainz_track_id
		),
		share_media AS (
				SELECT sh.id AS share_id, sh.owner_user_id, sh.token, sh.description, sh.created_at, sh.expires_at, sh.visit_count, sh.last_visited, sm.media_id
				FROM shares sh
				JOIN shared_media sm ON sm.share_id = sh.id
				WHERE sh.owner_user_id = ?
		),
		shared_metadata AS (
				SELECT sm.share_id, sm.owner_user_id, sm.token, sm.description, sm.created_at, sm.expires_at, sm.visit_count, sm.last_visited, m.*
				FROM share_media sm
				JOIN metadata m ON m.musicbrainz_track_id = sm.media_id
				UNION
				SELECT sm.share_id, sm.owner_user_id, sm.token, sm.description, sm.created_at, sm.expires_at, sm.visit_count, sm.last_visited, m.*
				FROM share_media sm
				JOIN metadata m ON m.musicbrainz_album_id = sm.media_id
				UNION
				SELECT sm.share_id, sm.owner_user_id, sm.token, sm.description, sm.created_at, sm.expires_at, sm.visit_count, sm.last_visited, m.*
				FROM share_media sm
				JOIN metadata m ON m.musicbrainz_artist_id = sm.media_id
		)
		SELECT
				sm.share_id AS id, sm.token, sm.description, sm.created_at, sm.expires_at, sm.visit_count, sm.last_visited, u.username,
				m.musicbrainz_track_id AS id, m.musicbrainz_album_id AS album_id, m.title, m.album, m.artist,
				COALESCE(m.track_number, 0) AS track, REPLACE(PRINTF('%4s', substr(m.release_date, 1, 4)),' ','0') AS year,
				substr(m.genre,1,instr(m.genre, ';') - 1) AS genre, m.musicbrainz_track_id AS cover_art, m.size,
				m.duration, m.bitrate, m.file_path AS path, m.date_added AS created, m.disc_number, m.musicbrainz_artist_id AS artist_id,
				m.album_artist, m.bit_depth, m.sample_rate, m.channels, COALESCE(ur.rating, 0) AS user_rating,
				COALESCE(gr.avg_rating, 0.0) AS average_rating, COALESCE(plays.play_count, 0) AS play_count, plays.last_played AS played,
			us.created_at AS starred, maa.musicbrainz_artist_id, maa.artist AS album_artist_name
		FROM shared_metadata sm
		JOIN metadata m ON m.musicbrainz_track_id = sm.musicbrainz_track_id
		JOIN user_music_folders f ON f.folder_id = m.music_folder_id AND f.user_id = sm.owner_user_id
		JOIN users u ON u.id = sm.owner_user_id
		LEFT JOIN user_stars us ON us.metadata_id = m.musicbrainz_track_id AND us.user_id = sm.owner_user_id
		LEFT JOIN user_ratings ur ON ur.metadata_id = m.musicbrainz_track_id AND ur.user_id = sm.owner_user_id
		LEFT JOIN gr ON gr.metadata_id = m.musicbrainz_artist_id
		LEFT JOIN plays ON plays.musicbrainz_track_id = m.musicbrainz_track_id
		LEFT JOIN metadata maa ON maa.artist = m.album_artist
		group by sm.share_id, m.musicbrainz_track_id
		ORDER BY sm.created_at DESC, m.musicbrainz_artist_id, m.musicbrainz_album_id, m.musicbrainz_track_id`

	rows, err := DB.QueryContext(ctx, query, owner.Id, owner.Id)
	if err != nil {
		return nil, fmt.Errorf("querying shares: %v", err)
	}
	defer rows.Close()

	var shares []types.ShareRow

	for rows.Next() {
		var row types.DbUserShare
		var lastPlayed sql.NullString
		var dateStarred sql.NullString
		var durationFloat float64
		var albumArtistName sql.NullString
		var albumArtistId sql.NullString
		var shareExpiresAt sql.NullString
		var shareLastVisited sql.NullString
		if err := rows.Scan(
			&row.ShareId, &row.Token, &row.Description, &row.ShareCreated, &shareExpiresAt, &row.VisitCount, &shareLastVisited, &row.ShareOwner,
			&row.TrackId, &row.AlbumId, &row.Title, &row.Album, &row.Artist, &row.TrackNumber, &row.Year,
			&row.Genre, &row.CoverArt, &row.Size, &durationFloat, &row.Bitrate, &row.Path, &row.DateAdded,
			&row.DiscNumber, &row.ArtistId, &row.AlbumArtist, &row.BitDepth, &row.SampleRate, &row.Channels,
			&row.UserRating, &row.AverageRating, &row.PlayCount, &lastPlayed, &dateStarred, &albumArtistId, &albumArtistName,
		); err != nil {
			return nil, fmt.Errorf("scanning share row: %v", err)
		}

		var entry types.SubsonicChild
		entry.Genres = []types.ChildGenre{}
		for _, genre := range strings.Split(row.Genre, ";") {
			entry.Genres = append(entry.Genres, types.ChildGenre{Name: genre})
		}

		if lastPlayed.Valid {
			entry.Played = lastPlayed.String
		}

		if dateStarred.Valid {
			entry.Starred = dateStarred.String
		}

		if dateStarred.Valid {
			entry.Starred = dateStarred.String
		}

		entry.Id = row.TrackId
		entry.Duration = int(durationFloat)
		entry.IsDir = false
		entry.MusicBrainzId = entry.Id
		entry.AlbumId = entry.Parent
		entry.Title = row.Title
		entry.Album = row.Album
		entry.Artist = row.Artist
		entry.Track = row.TrackNumber
		entry.Year = row.Year
		entry.Genre = row.Genre
		entry.CoverArt = row.CoverArt
		entry.Size = row.Size
		entry.BitRate = row.Bitrate
		entry.Path = row.Path
		entry.Created = row.DateAdded
		entry.DiscNumber = row.DiscNumber
		entry.ArtistId = row.ArtistId
		entry.BitDepth = row.BitDepth
		entry.SamplingRate = row.SampleRate
		entry.ChannelCount = row.Channels
		entry.UserRating = row.UserRating
		entry.AverageRating = row.AverageRating
		entry.PlayCount = row.PlayCount

		entry.Artists = []types.ChildArtist{}
		entry.Artists = append(entry.Artists, types.ChildArtist{Id: entry.ArtistId, Name: entry.Artist})

		entry.DisplayArtist = entry.Artist

		entry.AlbumArtists = []types.ChildArtist{}
		if albumArtistId.Valid && albumArtistName.Valid {
			entry.AlbumArtists = append(entry.AlbumArtists, types.ChildArtist{Id: albumArtistId.String, Name: albumArtistName.String})
		}

		entry.DisplayAlbumArtist = albumArtistName.String

		currentShare := slices.IndexFunc(shares, func(s types.ShareRow) bool { return s.Id == row.ShareId })

		var expiresAt string
		if shareExpiresAt.Valid {
			expiresAt = shareExpiresAt.String
		}

		var lastVisited string
		if shareLastVisited.Valid {
			lastVisited = shareLastVisited.String
		}

		if currentShare == -1 {
			shares = append(shares, types.ShareRow{
				Id:          row.ShareId,
				Description: row.Description,
				Username:    row.ShareOwner,
				Url:         logic.GetShareUrl(row.Token),
				Created:     row.ShareCreated,
				VisitCount:  row.VisitCount,
				Expires:     expiresAt,
				LastVisited: lastVisited,
				Entries:     []types.SubsonicChild{},
			})
			currentShare = len(shares) - 1
		}

		shares[currentShare].Entries = append(shares[currentShare].Entries, entry)
	}

	return shares, nil
}

func GetShareById(ctx context.Context, id int) (types.ShareRow, error) {
	query := `select s.id, u.username, s.token, s.description, s.created_at, s.expires_at, s.visit_count, s.last_visited
		from shares s
		join users u on u.id = s.owner_user_id
		where s.id = ?
		limit 1`

	var share types.ShareRow
	var token string
	var expires sql.NullString
	var lastVisited sql.NullString

	err := DB.QueryRowContext(ctx, query, id).Scan(
		&share.Id, &share.Username, &token, &share.Description, &share.Created, &expires, &share.VisitCount, &lastVisited,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.ShareRow{}, fmt.Errorf("no share found with id %d", id)
		}
		return types.ShareRow{}, fmt.Errorf("querying share by id: %v", err)
	}

	if expires.Valid && expires.String != "" && logic.GetTimeFromString(expires.String).Before(time.Now()) {
		return types.ShareRow{}, fmt.Errorf("share has expired")
	}

	mediaQuery := `SELECT media_id FROM shared_media WHERE share_id = ?`
	mediaRows, err := DB.QueryContext(ctx, mediaQuery, share.Id)
	if err != nil {
		if mediaRows != nil {
			mediaRows.Close()
		}
		return types.ShareRow{}, fmt.Errorf("querying shared media: %v", err)
	}

	var mediaIds []string
	for mediaRows.Next() {
		var mediaId string
		if err := mediaRows.Scan(&mediaId); err != nil {
			mediaRows.Close()
			return types.ShareRow{}, fmt.Errorf("scanning shared media: %v", err)
		}
		mediaIds = append(mediaIds, mediaId)
	}
	if err := mediaRows.Err(); err != nil {
		mediaRows.Close()
		return types.ShareRow{}, fmt.Errorf("iterating over shared media: %v", err)
	}
	mediaRows.Close()

	entries, err := GetSongsByShareId(ctx, share.Id)
	if err != nil {
		return types.ShareRow{}, fmt.Errorf("getting songs by share ID: %v", err)
	}

	if expires.Valid {
		share.Expires = expires.String
	}

	if lastVisited.Valid {
		share.LastVisited = lastVisited.String
	}

	share.Entries = entries

	return share, nil
}

func GetShareByToken(ctx context.Context, token string) (types.ShareRow, error) {
	query := `select s.id, u.username, s.token, s.description, s.created_at, s.expires_at, s.visit_count, s.last_visited
		from shares s
		join users u on u.id = s.owner_user_id
		where s.token = ?
		limit 1`

	var share types.ShareRow
	var expires sql.NullString
	var lastVisited sql.NullString

	err := DB.QueryRowContext(ctx, query, token).Scan(
		&share.Id, &share.Username, &token, &share.Description, &share.Created, &expires, &share.VisitCount, &lastVisited,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.ShareRow{}, fmt.Errorf("no share found with token %s", token)
		}
		return types.ShareRow{}, fmt.Errorf("querying share by token %s: %v", token, err)
	}

	if expires.Valid && expires.String != "" && logic.GetTimeFromString(expires.String).Before(time.Now()) {
		return types.ShareRow{}, fmt.Errorf("share has expired")
	}

	mediaQuery := `SELECT media_id FROM shared_media WHERE share_id = ?`
	mediaRows, err := DB.QueryContext(ctx, mediaQuery, share.Id)
	if err != nil {
		if mediaRows != nil {
			mediaRows.Close()
		}
		return types.ShareRow{}, fmt.Errorf("querying shared media: %v", err)
	}

	var mediaIds []string
	for mediaRows.Next() {
		var mediaId string
		if err := mediaRows.Scan(&mediaId); err != nil {
			mediaRows.Close()
			return types.ShareRow{}, fmt.Errorf("scanning shared media: %v", err)
		}
		mediaIds = append(mediaIds, mediaId)
	}
	if err := mediaRows.Err(); err != nil {
		mediaRows.Close()
		return types.ShareRow{}, fmt.Errorf("iterating over shared media: %v", err)
	}
	mediaRows.Close()

	entries, err := GetSongsByShareId(ctx, share.Id)
	if err != nil {
		return types.ShareRow{}, fmt.Errorf("getting songs by share ID: %v", err)
	}

	share.Entries = entries
	share.Url = logic.GetShareUrl(token)

	if expires.Valid {
		share.Expires = expires.String
	}

	if lastVisited.Valid {
		share.LastVisited = lastVisited.String
	}

	return share, nil
}

func ClearExpiredShares(ctx context.Context) error {
	query := `DELETE FROM shares WHERE expires_at IS NOT NULL AND expires_at < ?`
	now := logic.GetCurrentTimeFormatted()
	_, err := DB.ExecContext(ctx, query, now)
	if err != nil {
		return err
	}
	return nil
}

func GetSongsByShareId(ctx context.Context, shareId int) ([]types.SubsonicChild, error) {
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
		from shares sh
		join shared_media sm on sm.share_id = sh.id
		join metadata m on (m.musicbrainz_track_id = sm.media_id or m.musicbrainz_album_id = sm.media_id or m.musicbrainz_artist_id = sm.media_id)
		join user_music_folders f on f.folder_id = m.music_folder_id
		join track_genres g on m.file_path = g.file_path
		LEFT JOIN user_stars us ON m.musicbrainz_track_id = us.metadata_id AND us.user_id = sh.owner_user_id
		LEFT JOIN user_ratings ur ON m.musicbrainz_track_id = ur.metadata_id AND ur.user_id = sh.owner_user_id
		LEFT JOIN user_ratings gr ON m.musicbrainz_track_id = gr.metadata_id
		LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = sh.owner_user_id
		left join metadata maa on maa.artist = m.album_artist
		where sh.id = ?
		group by m.musicbrainz_track_id order by m.musicbrainz_artist_id, m.musicbrainz_album_id, m.musicbrainz_track_id`

	rows, err := DB.QueryContext(ctx, query, shareId)
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

func ValidateShareToken(ctx context.Context, shareToken string) (int, bool) {
	query := `SELECT s.id,
			s.expires_at
		FROM shares s
		JOIN users u ON u.id = s.owner_user_id
		WHERE u.share_role = 1
		AND s.token = ?
		LIMIT 1`

	var shareId int
	var expires sql.NullString

	err := DB.QueryRowContext(ctx, query, shareToken).Scan(&shareId, &expires)

	if err != nil {
		logger.Printf("Query failed: %v", err)
		return 0, false
	}

	if expires.Valid && expires.String != "" {
		expiresTime := logic.GetTimeFromString(expires.String)
		if !expiresTime.IsZero() && expiresTime.Before(time.Now()) {
			return 0, false
		}
	}

	return shareId, true
}

func ValidateTokenAndMediaId(ctx context.Context, token string, mediaId string) (int, bool) {
	query := `SELECT s.id, s.expires_at
		FROM shares s
		JOIN users u ON u.id = s.owner_user_id
		JOIN shared_media sm ON sm.share_id = s.id
		WHERE u.share_role = 1
		AND s.token = ?
		AND sm.media_id = ?
		LIMIT 1`

	var shareId int
	var expires sql.NullString

	err := DB.QueryRowContext(ctx, query, token, mediaId).Scan(&shareId, &expires)
	if err != nil {
		logger.Printf("Error validating mediaId %s for token %s: %v", mediaId, token, err)
		return 0, false
	}

	if expires.Valid && expires.String != "" {
		expiresTime := logic.GetTimeFromString(expires.String)
		if !expiresTime.IsZero() && expiresTime.Before(time.Now()) {
			return 0, false
		}
	}
	return shareId, true
}
