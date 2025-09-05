package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func migratePlaylists(ctx context.Context) {
	createPlaylistsTable(ctx)
	createPlaylistsAllowedUsersTable(ctx)
	createPlaylistEntriesTable(ctx)
}

func createPlaylistsTable(ctx context.Context) {
	schema := `CREATE TABLE playlists (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		name        TEXT NOT NULL,
    comment     TEXT,
    user_id     INTEGER,
    public      BOOLEAN DEFAULT FALSE,
    created     TEXT NOT NULL,
    changed     TEXT NOT NULL,
    cover_art   TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		UNIQUE (name)
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_playlists_user", "playlists", []string{"user_id"}, false)
	createIndex(ctx, "idx_playlists_name", "playlists", []string{"name"}, false)
}

func createPlaylistsAllowedUsersTable(ctx context.Context) {
	schema := `CREATE TABLE playlist_allowed_users (
    playlist_id INTEGER NOT NULL,
    user_id     INTEGER NOT NULL,
    PRIMARY KEY (playlist_id, user_id),
    FOREIGN KEY (playlist_id) REFERENCES playlists(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_playlists_user", "playlist_allowed_users", []string{"user_id"}, false)
}

func createPlaylistEntriesTable(ctx context.Context) {
	schema := `CREATE TABLE playlist_entries (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    playlist_id          INTEGER NOT NULL,
    musicbrainz_track_id TEXT NOT NULL,
    sort_order           INTEGER NOT NULL,
    FOREIGN KEY (playlist_id) REFERENCES playlists(id) ON DELETE CASCADE,
    UNIQUE (playlist_id, sort_order)
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_playlist_entries_playlist", "playlist_entries", []string{"playlist_id"}, false)
}

func CreatePlaylist(ctx context.Context, playlistName string, playlistId int, songIds []string) (types.PlaylistRow, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return types.PlaylistRow{}, err
	}

	// if the playlist already exists, update it
	exists, err := PlaylistExists(ctx, playlistId, playlistName)
	if err != nil {
		return types.PlaylistRow{}, err
	}
	if exists && playlistId > 0 && len(songIds) > 0 {
		err := addPlaylistEntries(ctx, playlistId, songIds)
		if err != nil {
			return types.PlaylistRow{}, fmt.Errorf("updating playlist via CreatePlaylist: %v", err)
		}

		err = updatePlaylistChangedDate(ctx, playlistId)
		if err != nil {
			return types.PlaylistRow{}, fmt.Errorf("updating playlist changed date: %v", err)
		}

		changePlaylist, err := GetPlaylist(ctx, playlistId)
		if err != nil {
			return types.PlaylistRow{}, fmt.Errorf("getting playlist after updating: %v", err)
		}

		entries, err := GetPlaylistEntries(ctx, playlistId)
		if err != nil {
			return types.PlaylistRow{}, fmt.Errorf("getting playlist entries after updating: %v", err)
		}

		changePlaylist.Entries = entries

		return changePlaylist, nil

	} else if exists && len(songIds) == 0 {
		return types.PlaylistRow{}, fmt.Errorf("existing playlist provided with no new songIds")
	} else if exists && playlistId == 0 {
		return types.PlaylistRow{}, fmt.Errorf("existing playlists should be referenced by playlistId, not name")
	}

	var newPlaylistId int

	if playlistName != "" {
		logger.Printf("Creating new playlist for user %s with name %s", user.Username, playlistName)

		query := `INSERT INTO playlists (name, user_id, created, changed) VALUES (?, ?, ?, ?);`
		result, err := DB.ExecContext(ctx, query,
			playlistName,
			user.Id,
			logic.GetCurrentTimeFormatted(),
			logic.GetCurrentTimeFormatted())
		if err != nil {
			return types.PlaylistRow{}, fmt.Errorf("creating playlist: %v", err)
		}
		lastInserted, err := result.LastInsertId()
		newPlaylistId = int(lastInserted)

		if len(songIds) > 0 {
			err = addPlaylistEntries(ctx, newPlaylistId, songIds)
			if err != nil {
				return types.PlaylistRow{}, fmt.Errorf("adding entries to new playlist via CreatePlaylist: %v", err)
			}
		}

		err = updateAllowedUsersForPlaylist(ctx, newPlaylistId, []int{user.Id})
		if err != nil {
			return types.PlaylistRow{}, fmt.Errorf("updating allowed users for new playlist: %v", err)
		}
	} else {
		return types.PlaylistRow{}, fmt.Errorf("either existing playlistId or new name parameter must be provided")
	}

	logger.Printf("Created new playlist with id %d", newPlaylistId)

	newPlaylist, err := GetPlaylist(ctx, newPlaylistId)
	if err != nil {
		return types.PlaylistRow{}, fmt.Errorf("getting playlist after creation: %v", err)
	}

	entries, err := GetPlaylistEntries(ctx, newPlaylistId)
	if err != nil {
		return types.PlaylistRow{}, fmt.Errorf("getting playlist entries after creation: %v", err)
	}

	newPlaylist.Entries = entries

	return newPlaylist, nil
}

func PlaylistExists(ctx context.Context, playlistId int, playlistName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM playlists WHERE id = ? OR name = ?);`
	err := DB.QueryRowContext(ctx, query, playlistId, playlistName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking if playlist exists: %v", err)
	}
	return exists, nil
}

func GetPlaylistIdByName(ctx context.Context, playlistName string) (int, error) {
	var playlistId int
	query := `SELECT id FROM playlists WHERE name = ?;`
	err := DB.QueryRowContext(ctx, query, playlistName).Scan(&playlistId)
	if err != nil {
		return 0, fmt.Errorf("getting playlist id by name: %v", err)
	}
	return playlistId, nil
}

func addPlaylistEntry(ctx context.Context, playlistId int, musicbrainzTrackId string) error {
	query := `INSERT INTO playlist_entries (playlist_id, musicbrainz_track_id, sort_order) VALUES (?, ?, (select COALESCE(max(sort_order)+1, 1) from playlist_entries where playlist_id = ?));`
	_, err := DB.ExecContext(ctx, query,
		playlistId,
		musicbrainzTrackId,
		playlistId)
	if err != nil {
		return fmt.Errorf("adding playlist entry: %v", err)
	}
	return nil
}

func addPlaylistEntries(ctx context.Context, playlistId int, songIds []string) error {
	// batch insert using a transaction
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %v", err)
	}
	defer tx.Rollback()

	for _, songId := range songIds {
		err := addPlaylistEntry(ctx, playlistId, songId)
		if err != nil {
			return fmt.Errorf("adding playlist entries: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %v", err)
	}
	return nil
}

func removePlaylistEntriesByIndexes(ctx context.Context, playlistId int, songIdIndexes []int) error {
	existingEntries, err := GetPlaylistEntries(ctx, playlistId)
	if err != nil {
		return fmt.Errorf("getting existing playlist entries: %v", err)
	}

	entryCount := len(existingEntries)
	for _, songIdIndex := range songIdIndexes {
		if songIdIndex < 0 || songIdIndex >= entryCount {
			return fmt.Errorf("songIdIndex %d out of range (playlist has %d entries)", songIdIndex, entryCount)
		}
		trackId := existingEntries[songIdIndex].Id
		_, err := DB.ExecContext(ctx, `DELETE FROM playlist_entries WHERE playlist_id = ? AND musicbrainz_track_id = ?`, playlistId, trackId)
		if err != nil {
			logger.Printf("Error removing track %s from playlist %d: %v", existingEntries[songIdIndex].Title, playlistId, err)
			return fmt.Errorf("removing playlist entry: %v", err)
		} else {
			logger.Printf("Successfully removed track %s from playlist %d", existingEntries[songIdIndex].Title, playlistId)
		}
	}

	return nil
}

func updateAllowedUsersForPlaylist(ctx context.Context, playlistId int, allowedUserIds []int) error {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("getting user by context: %v", err)
	}

	// remove unused user access
	if len(allowedUserIds) == 0 {
		// if no allowed users, remove all except the owner
		_, err := DB.ExecContext(ctx, `DELETE FROM playlist_allowed_users WHERE playlist_id = ? AND user_id != ?`, playlistId, user.Id)
		if err != nil {
			return fmt.Errorf("removing all allowed users: %v", err)
		}
	} else {
		// Build placeholders for NOT IN clause
		placeholders := make([]string, len(allowedUserIds))
		args := make([]interface{}, 0, len(allowedUserIds)+2)
		args = append(args, playlistId)
		for i, uid := range allowedUserIds {
			placeholders[i] = "?"
			args = append(args, uid)
		}
		args = append(args, user.Id)
		query := "DELETE FROM playlist_allowed_users WHERE playlist_id = ? AND user_id NOT IN (" + strings.Join(placeholders, ",") + ") AND user_id != ?"
		_, err := DB.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("removing old allowed users: %v", err)
		}
	}

	// add new allowed users
	for _, userId := range allowedUserIds {
		_, err := DB.ExecContext(ctx, `INSERT OR IGNORE INTO playlist_allowed_users (playlist_id, user_id) VALUES (?, ?)`, playlistId, userId)
		if err != nil {
			return fmt.Errorf("adding allowed user to playlist: %v", err)
		}
	}

	return nil
}

func RemoveOrphanedPlaylistEntries(ctx context.Context) error {
	query := `DELETE FROM playlist_entries
	WHERE musicbrainz_track_id NOT IN (SELECT musicbrainz_track_id FROM metadata);`
	_, err := DB.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("removing orphaned playlist entries: %v", err)
	}
	return nil
}

func GetPlaylists(ctx context.Context, username string) ([]types.PlaylistRow, error) {
	query := `select 
    p.id as id,
    p.name as name,
    u.username as owner,
    p.public,
    p.created,
    p.changed,
    count(pe.musicbrainz_track_id) as song_count,
    coalesce(cast(sum(m.duration) as integer), 0) as duration,
    coalesce(p.comment, '') as comment,
    coalesce(coalesce(p.cover_art, min(pe.musicbrainz_track_id)), '') as cover_art,
    au.allowed_users
	from playlists p
	join users u on u.id = p.user_id
	left join playlist_entries pe on pe.playlist_id = p.id
	left join metadata m on m.musicbrainz_track_id = pe.musicbrainz_track_id
	left join (
		select playlist_id, group_concat(u.username, ',') as allowed_users
		from playlist_allowed_users pau
		join users u on u.id = pau.user_id
		group by playlist_id
	) au on au.playlist_id = p.id
	where u.username = ?
	group by p.id, au.allowed_users;`

	rows, err := DB.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("querying playlists: %v", err)
	}
	defer rows.Close()

	var playlists []types.PlaylistRow
	for rows.Next() {
		var playlist types.PlaylistRow
		var allowedUsersString string
		if err := rows.Scan(&playlist.Id, &playlist.Name, &playlist.Owner, &playlist.Public, &playlist.Created, &playlist.Changed,
			&playlist.SongCount, &playlist.Duration, &playlist.Comment, &playlist.CoverArt, &allowedUsersString); err != nil {
			return nil, fmt.Errorf("scanning row in GetPlaylists: %v", err)
		}
		playlist.AllowedUsers = strings.Split(allowedUsersString, ",")
		playlists = append(playlists, playlist)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over playlists in GetPlaylists: %v", err)
	}

	return playlists, nil
}

func GetPlaylist(ctx context.Context, playlistId int) (types.PlaylistRow, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return types.PlaylistRow{}, err
	}

	query := `select 
    p.id as id,
    p.name as name,
    u.username as owner,
    p.public,
    p.created,
    p.changed,
    count(pe.musicbrainz_track_id) as song_count,
    coalesce(cast(sum(m.duration) as integer), 0) as duration,
    coalesce(p.comment, '') as comment,
    coalesce(coalesce(p.cover_art, min(pe.musicbrainz_track_id)), '') as cover_art,
    au.allowed_users
	from playlists p
	join users u on u.id = p.user_id
	left join playlist_entries pe on pe.playlist_id = p.id
	left join metadata m on m.musicbrainz_track_id = pe.musicbrainz_track_id
	left join (
		select playlist_id, group_concat(u.username, ',') as allowed_users
		from playlist_allowed_users pau
		join users u on u.id = pau.user_id
		group by playlist_id
	) au on au.playlist_id = p.id
	where p.id = ?
	group by p.id, au.allowed_users;`

	var result types.PlaylistRow
	var allowedUsersString string

	err = DB.QueryRowContext(ctx, query, playlistId).Scan(&result.Id, &result.Name, &result.Owner, &result.Public, &result.Created, &result.Changed,
		&result.SongCount, &result.Duration, &result.Comment, &result.CoverArt, &allowedUsersString)
	if err == sql.ErrNoRows {
		return types.PlaylistRow{}, nil
	} else if err != nil {
		return types.PlaylistRow{}, err
	}

	if allowedUsersString != "" {
		result.AllowedUsers = strings.Split(allowedUsersString, ",")
	}

	if result.Owner != user.Username && !user.AdminRole {
		return types.PlaylistRow{}, fmt.Errorf("user %s not authorized to access playlist %d owned by %s", user.Username, playlistId, result.Owner)
	}

	return result, nil
}

func GetPlaylistEntries(ctx context.Context, playlistId int) ([]types.SubsonicChild, error) {
	query := `select m.musicbrainz_track_id as id, m.musicbrainz_album_id as parent, m.title, m.album, m.artist, COALESCE(m.track_number, 0) as track,
		REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year, substr(m.genre,1,(instr(m.genre,';')-1)) as genre, m.musicbrainz_track_id as cover_art,
		m.size, m.duration, m.bitrate, m.file_path as path, m.date_added as created, m.disc_number, m.musicbrainz_artist_id as artist_id,
		m.genre, m.album_artist, m.bit_depth, m.sample_rate, m.channels,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(AVG(gr.rating), 0.0) AS average_rating,
		COALESCE(SUM(pc.play_count), 0) AS play_count,
		max(pc.last_played) as played,
		us.created_at AS starred,
		maa.musicbrainz_artist_id as album_artist_id
	from playlists p
	join playlist_entries pe on pe.playlist_id = p.id
	join users u on u.id = p.user_id
	join user_music_folders uf on uf.user_id = u.id
	join metadata m on m.music_folder_id = uf.folder_id and m.musicbrainz_track_id = pe.musicbrainz_track_id
	LEFT JOIN user_stars us ON m.musicbrainz_album_id = us.metadata_id AND us.user_id = uf.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_album_id = ur.metadata_id AND ur.user_id = u.id
	LEFT JOIN user_ratings gr ON m.musicbrainz_album_id = gr.metadata_id
	LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = u.id
	left join metadata maa on maa.artist = m.album_artist
	where p.id = ?
	group by m.musicbrainz_track_id
	order by pe.sort_order asc`

	var results []types.SubsonicChild

	rows, err := DB.QueryContext(ctx, query, playlistId)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.SubsonicChild{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var result types.SubsonicChild

		var albumArtistName sql.NullString
		var albumArtistId sql.NullString
		var genreString string
		var durationFloat float64
		var played sql.NullString
		var starred sql.NullString

		if err := rows.Scan(&result.Id, &result.Parent, &result.Title, &result.Album, &result.Artist, &result.Track,
			&result.Year, &result.Genre, &result.CoverArt,
			&result.Size, &durationFloat, &result.BitRate, &result.Path, &result.Created, &result.DiscNumber, &result.ArtistId,
			&genreString, &albumArtistName, &result.BitDepth, &result.SamplingRate, &result.ChannelCount,
			&result.UserRating, &result.AverageRating, &result.PlayCount, &played, &starred, &albumArtistId); err != nil {
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

func DeletePlaylist(ctx context.Context, playlistId int) error {
	query := `DELETE FROM playlists WHERE id = ?`
	_, err := DB.ExecContext(ctx, query, playlistId)
	return err
}

func UpdatePlaylist(ctx context.Context, playlistId int, playlistName, comment string, public string, coverArt string, allowedUsers []int, songIdsToAdd []string, songIndexesToRemove []int) error {
	if playlistId == 0 && playlistName == "" {
		return fmt.Errorf("either existing playlistId or new name parameter must be provided")
	}

	if playlistId == 0 {
		var err error
		playlistId, err = GetPlaylistIdByName(ctx, playlistName)
		if err != nil {
			return fmt.Errorf("getting playlist id by name: %v", err)
		}
	}

	user, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("getting user by context: %v", err)
	}

	currentPlaylist, err := GetPlaylist(ctx, playlistId)
	if err != nil {
		return fmt.Errorf("getting current playlist: %v", err)
	}

	if currentPlaylist.Owner != user.Username && !user.AdminRole {
		return fmt.Errorf("user is not the owner of the playlist")
	}

	if playlistName != "" || comment != "" || public != "" || coverArt != "" {
		var args []interface{}

		query := `UPDATE playlists SET changed = ?`
		args = append(args, logic.GetCurrentTimeFormatted())

		if playlistName != "" {
			query += ` name = ?,`
			args = append(args, playlistName)
		}
		if comment != "" {
			query += ` comment = ?,`
			args = append(args, comment)
		}
		if public != "" {
			query += ` public = ?,`
			args = append(args, public)
		}
		if coverArt != "" {
			query += ` cover_art = ?,`
			args = append(args, coverArt)
		}

		// trim trailing comma from query if there is one
		query = strings.TrimSuffix(query, ",")
		query += ` WHERE id = ?`
		args = append(args, playlistId)

		_, err = DB.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("updating playlist: %v", err)
		}
	}

	if len(allowedUsers) > 0 {
		err := updateAllowedUsersForPlaylist(ctx, playlistId, allowedUsers)
		if err != nil {
			return fmt.Errorf("updating allowed users for playlist: %v", err)
		}
	}

	if len(songIdsToAdd) > 0 {
		err := addPlaylistEntries(ctx, playlistId, songIdsToAdd)
		if err != nil {
			return fmt.Errorf("adding entries to playlist: %v", err)
		}
	}

	if len(songIndexesToRemove) > 0 {
		err := removePlaylistEntriesByIndexes(ctx, playlistId, songIndexesToRemove)
		if err != nil {
			return fmt.Errorf("removing entries from playlist: %v", err)
		}
	}

	return nil
}

func updatePlaylistChangedDate(ctx context.Context, playlistId int) error {
	currentTime := logic.GetCurrentTimeFormatted()
	query := `UPDATE playlists SET changed = ? WHERE id = ?`
	_, err := DB.ExecContext(ctx, query, currentTime, playlistId)
	return err
}
