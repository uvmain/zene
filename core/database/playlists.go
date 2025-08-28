package database

import (
	"context"
	"fmt"
	"strings"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func createPlaylistTables(ctx context.Context) {
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
    position             INTEGER NOT NULL,
    FOREIGN KEY (playlist_id) REFERENCES playlists(id) ON DELETE CASCADE,
    UNIQUE (playlist_id, position)
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
	exists, err := playlistExists(ctx, playlistId, playlistName)
	if err != nil {
		return types.PlaylistRow{}, err
	}
	if exists && playlistId > 0 && len(songIds) > 0 {
		err := addPlaylistEntries(ctx, playlistId, songIds)
		if err != nil {
			return types.PlaylistRow{}, fmt.Errorf("updating playlist via CreatePlaylist: %v", err)
		}
		return types.PlaylistRow{}, nil
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

	return types.PlaylistRow{}, nil
}

func playlistExists(ctx context.Context, playlistId int, playlistName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM playlists WHERE id = ? OR name = ?);`
	err := DB.QueryRowContext(ctx, query, playlistId, playlistName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking if playlist exists: %v", err)
	}
	return exists, nil
}

func addPlaylistEntry(ctx context.Context, playlistId int, musicbrainzTrackId string) error {
	query := `INSERT INTO playlist_entries (playlist_id, musicbrainz_track_id, position) VALUES (?, ?, (select COALESCE(max(position)+1, 1) from playlist_entries where playlist_id = ?));`
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

func updateAllowedUsersForPlaylist(ctx context.Context, playlistId int, allowedUserIds []int) error {
	// remove unused user access
	if len(allowedUserIds) == 0 {
		// if no allowed users, remove all
		_, err := DB.ExecContext(ctx, `DELETE FROM playlist_allowed_users WHERE playlist_id = ?`, playlistId)
		if err != nil {
			return fmt.Errorf("removing all allowed users: %v", err)
		}
	} else {
		placeholders := ""
		args := make([]interface{}, 0, len(allowedUserIds)+1)
		args = append(args, playlistId)
		for i := range allowedUserIds {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			args = append(args, allowedUserIds[i])
		}
		query := "DELETE FROM playlist_allowed_users WHERE playlist_id = ? AND user_id NOT IN (" + placeholders + ")"
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
	query := `select p.id as id,
		p.name as name,
		u.username as owner,
		p.public,
		p.created,
		p.changed,
		count(pe.musicbrainz_track_id) as song_count,
		cast(sum(m.duration) as integer) as duration,
		coalesce(p.comment, '') as comment,
		coalesce(p.cover_art, min(pe.musicbrainz_track_id)) as cover_art,
		au.allowed_users
	from playlists p
	join users u on u.id = p.user_id
	join playlist_entries pe on pe.playlist_id = p.id
	join (
		select playlist_id, group_concat(u.username, ',') as allowed_users
		from playlist_allowed_users pau
		join users u on u.id = pau.user_id
		group by playlist_id
		order by playlist_id
	) au on au.playlist_id = p.id
	join metadata m on m.musicbrainz_track_id = pe.musicbrainz_track_id
	where u.username = ?
	group by p.id, au.playlist_id;`

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
