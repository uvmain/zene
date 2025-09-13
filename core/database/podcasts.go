package database

import (
	"context"
	"fmt"
	"zene/core/logic"
)

func migratePodcasts(ctx context.Context) {
	schema := `CREATE TABLE podcast_channels (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		title TEXT,
		description TEXT,
    	cover_art TEXT,
		original_image_url TEXT,
    	last_refresh TEXT,
		status TEXT DEFAULT 'new', -- new / downloading / completed / error / deleted / skipped
	    created_at TEXT,
		error_message TEXT,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_podcast_channels_user_id", "podcast_channels", []string{"user_id"}, false)

	schema = `CREATE TABLE podcast_episodes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		channel_id INTEGER NOT NULL REFERENCES podcast_channels(id) ON DELETE CASCADE,
		title TEXT,
		album TEXT,
		artist TEXT,
		year TEXT,
		cover_art TEXT,
		size TEXT,
		content_type TEXT,
		suffix TEXT,
		duration INTEGER,
		bit_rate TEXT,
		description TEXT,
		publish_date TEXT,
		status TEXT DEFAULT 'new',  -- new / downloading / completed / error / deleted / skipped
		file_path TEXT,
		stream_id TEXT,
		created_at TEXT NOT NULL
	);`
}

func CreatePodcastChannel(ctx context.Context, url string, title string, description string, original_image_url string, lastRefresh string) error {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("getting user from context: %v", err)
	}

	if !user.PodcastRole {
		return fmt.Errorf("user not authorized to create Podcast channels")
	}

	createdAt := logic.GetCurrentTimeFormatted()

	query := `INSERT INTO podcast_channels (url, title, description, original_image_url, last_refresh, user_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = DB.ExecContext(ctx, query, url, title, description, original_image_url, lastRefresh, user.Id, createdAt)
	if err != nil {
		return fmt.Errorf("inserting podcast channel: %v", err)
	}

	return nil
}

func UpdatePodcastChannel(ctx context.Context, channelId int, url string, title string, description string, coverArt string, lastRefresh string) error {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("getting user from context: %v", err)
	}

	if !user.PodcastRole {
		return fmt.Errorf("user not authorized to update Podcast channels")
	}
	query := `UPDATE podcast_channels SET url = ?, title = ?, description = ?, cover_art = ?, last_refresh = ? WHERE id = ? AND user_id = ?`

	result, err := DB.ExecContext(ctx, query, url, title, description, coverArt, lastRefresh, channelId, user.Id)
	if err != nil {
		return fmt.Errorf("updating podcast channel: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no podcast channel found with the given ID for this user")
	}

	return nil
}

func DeletePodcastChannel(ctx context.Context, channelId int) error {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("getting user from context: %v", err)
	}

	if !user.PodcastRole {
		return fmt.Errorf("user not authorized to delete Podcast channels")
	}

	query := `DELETE FROM podcast_channels WHERE id = ? AND user_id = ?`

	result, err := DB.ExecContext(ctx, query, channelId, user.Id)
	if err != nil {
		return fmt.Errorf("deleting podcast channel: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no podcast channel found with the given ID for this user")
	}
	return nil
}
