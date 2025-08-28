package database

import (
	"context"
	"fmt"
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
		logger.Printf("Created playlist with ID %d", newPlaylistId)

		if len(songIds) > 0 {
			err = addPlaylistEntries(ctx, int(newPlaylistId), songIds)
			if err != nil {
				return types.PlaylistRow{}, fmt.Errorf("adding entries to new playlist via CreatePlaylist: %v", err)
			}
		}
	} else {
		return types.PlaylistRow{}, fmt.Errorf("either existing playlistId or new name parameter must be provided")
	}

	logger.Printf("Playlist ID: %d", newPlaylistId)

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
	// remove existing allowed users
	_, err := DB.ExecContext(ctx, `DELETE FROM playlist_allowed_users WHERE playlist_id = ?`, playlistId)
	if err != nil {
		return fmt.Errorf("removing old allowed users: %v", err)
	}

	// add new allowed users
	for _, userId := range allowedUserIds {
		_, err := DB.ExecContext(ctx, `INSERT INTO playlist_allowed_users (playlist_id, user_id) VALUES (?, ?)`, playlistId, userId)
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
