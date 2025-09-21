package database

import (
	"context"
	"database/sql"
	"fmt"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
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
		guid TEXT,
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
	createTable(ctx, schema)
	createIndex(ctx, "idx_podcast_episodes_channel_id", "podcast_episodes", []string{"channel_id"}, false)
}

func CreatePodcastChannel(ctx context.Context, url string, title string, description string, original_image_url string, cover_art string, lastRefresh string) error {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("getting user from context: %v", err)
	}

	if !user.PodcastRole {
		return fmt.Errorf("user not authorized to create Podcast channels")
	}

	createdAt := logic.GetCurrentTimeFormatted()

	query := `INSERT INTO podcast_channels (url, title, description, cover_art, original_image_url, last_refresh, user_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = DB.ExecContext(ctx, query, url, title, description, cover_art, original_image_url, lastRefresh, user.Id, createdAt)
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

func IsValidPodcastCover(ctx context.Context, coverArtId string) (bool, error) {
	query := `SELECT cover_art
		FROM (
				SELECT cover_art
				FROM podcast_channels
				WHERE cover_art = ?
				UNION ALL
				SELECT cover_art
				FROM podcast_episodes
				WHERE cover_art = ?
		)
		LIMIT 1;`
	row := DB.QueryRowContext(ctx, query, coverArtId, coverArtId)

	var dbCoverArtId string
	if err := row.Scan(&dbCoverArtId); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("checking podcast cover validity: %v", err)
	}

	return dbCoverArtId == coverArtId, nil
}

func GetPodcasts(ctx context.Context, podcastId int, includeEpisodes bool) ([]types.PodcastChannel, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return []types.PodcastChannel{}, err
	}

	var args []interface{}
	query := `select p.id,
		p.title,
		p.url,
		p.description,
		p.cover_art,
		p.original_image_url,
		p.last_refresh,
		p.status,
		p.created_at,
		p.error_message
	from users u
	join podcast_channels p on p.user_id = u.id
	where u.id = ?`
	args = append(args, user.Id)

	if podcastId != 0 {
		query += ` AND p.id = ?`
		args = append(args, podcastId)
	}

	rows, err := DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying podcasts: %v", err)
	}
	defer rows.Close()

	var channelArray []types.PodcastChannel

	for rows.Next() {
		var channel types.PodcastChannel
		var errorMessage sql.NullString
		var lastRefreshed sql.NullString
		if err := rows.Scan(
			&channel.Id,
			&channel.Title,
			&channel.Url,
			&channel.Description,
			&channel.CoverArt,
			&channel.OriginalImageUrl,
			&lastRefreshed,
			&channel.Status,
			&channel.CreatedAt,
			&errorMessage,
		); err != nil {
			return nil, fmt.Errorf("scanning channel row: %v", err)
		}

		channel.ParentId = channel.Id
		channel.ChannelId = channel.Id
		channel.Type = "podcast"
		channel.IsDir = "false"
		channel.IsVideo = "false"

		if lastRefreshed.Valid {
			channel.LastRefresh = lastRefreshed.String
		}
		if errorMessage.Valid {
			channel.ErrorMessage = errorMessage.String
		}

		if includeEpisodes {
			logger.Printf("fetch episodes and add to channel.Episodes")
		}

		channelArray = append(channelArray, channel)
	}

	return channelArray, nil
}
