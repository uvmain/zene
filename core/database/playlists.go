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
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name        TEXT NOT NULL,
    comment     TEXT,
    owner       INTEGER,
    public      BOOLEAN DEFAULT FALSE,
    created     TEXT NOT NULL,
    changed     TEXT NOT NULL,
    cover_art   TEXT,
		FOREIGN KEY (owner) REFERENCES users(id) ON DELETE CASCADE
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_playlists_user", "playlists", []string{"owner"}, false)
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
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    playlist_id INTEGER NOT NULL,
    track_id    INTEGER NOT NULL,
    position    INTEGER NOT NULL,
    FOREIGN KEY (playlist_id) REFERENCES playlists(id) ON DELETE CASCADE,
    FOREIGN KEY (track_id) REFERENCES tracks(id) ON DELETE CASCADE,
		UNIQUE (playlist_id, position)
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_playlist_entries_playlist", "playlist_entries", []string{"playlist_id"}, false)
}

func CreatePlaylist(ctx context.Context, playlistName string, playlistId int, songId string, public bool, allowedUserIds []int) (types.PlaylistRow, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return types.PlaylistRow{}, err
	}

	var newPlaylistId int64

	if playlistId > 0 {
		err := addPlaylistEntry(ctx, playlistId, songId)
		if err != nil {
			return types.PlaylistRow{}, fmt.Errorf("updating playlist via CreatePlaylist: %v", err)
		}
	} else if playlistName != "" {
		logger.Printf("Creating new playlist for user %s with name %s", user.Username, playlistName)

		query := `INSERT INTO playlists (name, owner, public, created, changed) VALUES (?, ?, ?, ?, ?);`
		result, err := DB.ExecContext(ctx, query,
			playlistName,
			user.Id,
			public,
			logic.GetCurrentTimeFormatted(),
			logic.GetCurrentTimeFormatted())
		if err != nil {
			return types.PlaylistRow{}, fmt.Errorf("creating playlist: %v", err)
		}
		newPlaylistId, err = result.LastInsertId()
		logger.Printf("Created playlist with ID %d", newPlaylistId)

		err = addPlaylistEntry(ctx, int(newPlaylistId), songId)
		if err != nil {
			return types.PlaylistRow{}, fmt.Errorf("adding entry to new playlist via CreatePlaylist: %v", err)
		}
	} else {
		return types.PlaylistRow{}, fmt.Errorf("either playlistId or playlistName must be provided")
	}

	playlistToUpdate := playlistId
	if playlistToUpdate == 0 {
		playlistToUpdate = int(newPlaylistId)
	}
	if len(allowedUserIds) > 0 {
		err = updateAllowedUsersForPlaylist(ctx, playlistToUpdate, allowedUserIds)
		if err != nil {
			return types.PlaylistRow{}, fmt.Errorf("updating allowed users for playlist: %v", err)
		}
	} else {
		err = updateAllowedUsersForPlaylist(ctx, playlistToUpdate, []int{user.Id})
		if err != nil {
			return types.PlaylistRow{}, fmt.Errorf("updating allowed users for playlist: %v", err)
		}
	}

	return types.PlaylistRow{}, nil
}

func addPlaylistEntry(ctx context.Context, playlistId int, songId string) error {
	query := `INSERT INTO playlist_entries (playlist_id, track_id, position) VALUES (?, ?, (select max(position)+1 from playlist_entries where playlist_id = ?));`
	_, err := DB.ExecContext(ctx, query,
		playlistId,
		songId,
		playlistId)
	if err != nil {
		return fmt.Errorf("adding playlist entry: %v", err)
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
