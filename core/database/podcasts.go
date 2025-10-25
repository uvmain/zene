package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
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
		categories TEXT,
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
		created_at TEXT NOT NULL,
		source_url TEXT,
		UNIQUE (channel_id, guid)
		FOREIGN KEY (channel_id) REFERENCES podcast_channels(id) ON DELETE CASCADE
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_podcast_episodes_channel_id", "podcast_episodes", []string{"channel_id"}, false)
}

func CreatePodcastChannel(ctx context.Context, url string, title string, description string, original_image_url string, cover_art string, lastRefresh string, categories []string) (int, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("getting user from context: %v", err)
	}

	if !user.PodcastRole {
		return 0, fmt.Errorf("user not authorized to create Podcast channels")
	}

	createdAt := logic.GetCurrentTimeFormatted()

	query := `INSERT INTO podcast_channels (url, title, description, cover_art, original_image_url, last_refresh, user_id, created_at, categories)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	categoriesString := strings.Join(categories, ",")

	result, err := DB.ExecContext(ctx, query, url, title, description, cover_art, original_image_url, lastRefresh, user.Id, createdAt, categoriesString)
	if err != nil {
		return 0, fmt.Errorf("inserting podcast channel: %v", err)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting last insert ID: %v", err)
	}

	return int(lastInsertId), nil
}

func UpdatePodcastChannel(ctx context.Context, channelId int, url string, title string, description string, originalImageUrl string, coverArt string, lastRefresh string, categories []string) error {
	query := `UPDATE podcast_channels SET url = ?, title = ?, description = ?, original_image_url = ?, cover_art = ?, last_refresh = ?, categories = ? WHERE id = ?`

	categoriesString := strings.Join(categories, ",")

	result, err := DB.ExecContext(ctx, query, url, title, description, originalImageUrl, coverArt, lastRefresh, categoriesString, channelId)
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

func UpdatePodcastChannelStatus(ctx context.Context, channelId int, status string) error {
	query := `UPDATE podcast_channels SET status = ? WHERE id = ?`
	_, err := DB.ExecContext(ctx, query, status, channelId)
	if err != nil {
		return fmt.Errorf("updating podcast channel status: %v", err)
	}
	return nil
}

func UpdatePodcastEpisodeStatus(ctx context.Context, episodeId int, status string) error {
	query := `UPDATE podcast_episodes SET status = ? WHERE id = ?`
	_, err := DB.ExecContext(ctx, query, status, episodeId)
	if err != nil {
		return fmt.Errorf("updating podcast episode status: %v", err)
	}
	return nil
}

func AddFileDetailsToEpisode(ctx context.Context, episodeId int, filePath string, size int64, contentType string, duration string, bitRate string) error {
	query := `UPDATE podcast_episodes SET file_path = ?, size = ?, content_type = ?, duration = ?, bit_rate = ? WHERE id = ?`
	_, err := DB.ExecContext(ctx, query, filePath, size, contentType, duration, bitRate, episodeId)
	if err != nil {
		return fmt.Errorf("adding file details to podcast episode: %v", err)
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
			channelIdInt, err := strconv.Atoi(channel.Id)
			if err != nil {
				return nil, fmt.Errorf("converting channel id to int: %v", err)
			}
			episodes, err := GetPodcastEpisodesByChannelId(ctx, channelIdInt)
			if err != nil {
				return nil, fmt.Errorf("getting podcast episodes by channel id: %v", err)
			}
			channel.Episodes = episodes
		}
		if channel.Episodes == nil {
			channel.Episodes = []types.PodcastEpisode{}
		}

		if len(channel.Episodes) > 0 {
			// get first episode's stream ID
			channel.StreamId = channel.Episodes[0].StreamId
		}

		channelArray = append(channelArray, channel)
	}

	return channelArray, nil
}

func UpsertPodcastEpisode(ctx context.Context, episode types.PodcastEpisodeRow) error {
	query := `INSERT INTO podcast_episodes (
		channel_id,
		guid,
		title,
		album,
		artist,
		year,
		cover_art,
		size,
		content_type,
		suffix,
		duration,
		bit_rate,
		description,
		publish_date,
		status,
		file_path,
		created_at,
		source_url
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT (id) DO UPDATE SET
		guid = excluded.guid,
		title = excluded.title,
		album = excluded.album,
		artist = excluded.artist,
		year = excluded.year,
		cover_art = excluded.cover_art,
		size = excluded.size,
		content_type = excluded.content_type,
		suffix = excluded.suffix,
		duration = excluded.duration,
		bit_rate = excluded.bit_rate,
		description = excluded.description,
		publish_date = excluded.publish_date,
		status = excluded.status,
		file_path = excluded.file_path,
		created_at = excluded.created_at,
		source_url = excluded.source_url
	`

	suffix := strings.Split(episode.Suffix, "?")[0] //remove trailing query params
	suffix = strings.TrimPrefix(suffix, ".")        //remove leading dot

	_, err = DB.ExecContext(ctx, query,
		episode.ChannelId,
		episode.Guid,
		episode.Title,
		episode.Album,
		episode.Artist,
		episode.Year,
		episode.CoverArt,
		episode.Size,
		episode.ContentType,
		suffix,
		episode.Duration,
		episode.BitRate,
		episode.Description,
		episode.PublishDate,
		episode.Status,
		episode.FilePath,
		episode.CreatedAt,
		episode.SourceUrl,
	)

	if err != nil {
		return fmt.Errorf("inserting podcast episode: %v", err)
	}

	return nil
}

func InsertPodcastEpisodes(episodes []types.PodcastEpisodeRow) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %v", err)
	}
	defer tx.Rollback()

	for _, episode := range episodes {
		if err := UpsertPodcastEpisode(ctx, episode); err != nil {
			return fmt.Errorf("upserting podcast episode: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %v", err)
	}

	return nil
}

func UpdatePodcastChannelLastRefresh(channelId int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	query := `UPDATE podcast_channels SET last_refresh = ? WHERE id = ?`
	_, err := DB.ExecContext(ctx, query, logic.GetCurrentTimeFormatted(), channelId)
	if err != nil {
		logger.Printf("Error updating podcast channel last refresh: %v", err)
		return fmt.Errorf("updating podcast channel last refresh: %v", err)
	}
	return nil
}

func GetPodcastEpisodeById(ctx context.Context, episodeId int) (types.PodcastEpisode, error) {
	var episode types.PodcastEpisode
	query := `select pe.id, pe.guid, pe.channel_id, pe.title, pe.description, pe.publish_date,
		pe.status, pe.id, 'false', pe.year, pc.categories, pe.cover_art, pe.size, pe.content_type,
		pe.suffix, pe.duration, pe.bit_rate, pe.file_path, pe.source_url
	from podcast_episodes pe
	join podcast_channels pc on pc.id = pe.channel_id
	where pe.id = ?`
	row := DB.QueryRowContext(ctx, query, episodeId)
	var genresString sql.NullString

	if err := row.Scan(
		&episode.Id,
		&episode.StreamId,
		&episode.ChannelId,
		&episode.Title,
		&episode.Description,
		&episode.PublishDate,
		&episode.Status,
		&episode.Parent,
		&episode.IsDir,
		&episode.Year,
		&genresString,
		&episode.CoverArt,
		&episode.Size,
		&episode.ContentType,
		&episode.Suffix,
		&episode.Duration,
		&episode.BitRate,
		&episode.Path,
		&episode.SourceUrl,
	); err != nil {
		if err == sql.ErrNoRows {
			return types.PodcastEpisode{}, nil
		}
		return types.PodcastEpisode{}, fmt.Errorf("getting podcast episode by id: %v", err)
	}

	episode.Genres = []types.ChildGenre{}
	for _, genre := range strings.Split(genresString.String, ",") {
		episode.Genres = append(episode.Genres, types.ChildGenre{Name: genre})
	}

	episode.Genre = episode.Genres[0].Name

	return episode, nil
}

func GetPodcastEpisodeByGuid(ctx context.Context, episodeGuid string) (types.PodcastEpisode, error) {
	var episode types.PodcastEpisode
	query := `select pe.id, pe.guid, pe.channel_id, pe.title, pe.description, pe.publish_date,
		pe.status, pe.id, 'false', pe.year, pc.categories, pe.cover_art, pe.size, pe.content_type,
		pe.suffix, pe.duration, pe.bit_rate, pe.file_path, pe.source_url
	from podcast_episodes pe
	join podcast_channels pc on pc.id = pe.channel_id
	where pe.guid = ?`
	row := DB.QueryRowContext(ctx, query, episodeGuid)
	var genresString sql.NullString

	if err := row.Scan(
		&episode.Id,
		&episode.StreamId,
		&episode.ChannelId,
		&episode.Title,
		&episode.Description,
		&episode.PublishDate,
		&episode.Status,
		&episode.Parent,
		&episode.IsDir,
		&episode.Year,
		&genresString,
		&episode.CoverArt,
		&episode.Size,
		&episode.ContentType,
		&episode.Suffix,
		&episode.Duration,
		&episode.BitRate,
		&episode.Path,
		&episode.SourceUrl,
	); err != nil {
		if err == sql.ErrNoRows {
			return types.PodcastEpisode{}, nil
		}
		return types.PodcastEpisode{}, fmt.Errorf("getting podcast episode by id: %v", err)
	}

	episode.Genres = []types.ChildGenre{}
	for _, genre := range strings.Split(genresString.String, ",") {
		episode.Genres = append(episode.Genres, types.ChildGenre{Name: genre})
	}

	episode.Genre = episode.Genres[0].Name

	return episode, nil
}

func DeletePodcastEpisodeById(ctx context.Context, episodeId string) error {
	query := `DELETE FROM podcast_episodes WHERE id = ?;`
	_, err := DB.ExecContext(ctx, query, episodeId)
	if err != nil {
		return fmt.Errorf("deleting podcast episode by id: %v", err)
	}
	return nil
}

func GetPodcastEpisodesByChannelId(ctx context.Context, channelId int) ([]types.PodcastEpisode, error) {
	var episodes []types.PodcastEpisode
	query := `select pe.id, pe.guid, pe.channel_id, pe.title, pe.description, pe.publish_date,
		pe.status, pe.id, 'false', pe.year, pc.categories, pe.cover_art, pe.size, pe.content_type,
		pe.suffix, pe.duration, pe.bit_rate, pe.file_path, pe.source_url
	from podcast_episodes pe
	join podcast_channels pc on pc.id = pe.channel_id
	where pc.id = ?
	order by pe.publish_date desc;`
	rows, err := DB.QueryContext(ctx, query, channelId)
	if err != nil {
		return nil, fmt.Errorf("querying podcast episodes by channel id: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var episode types.PodcastEpisode
		var genresString sql.NullString
		if err := rows.Scan(
			&episode.Id,
			&episode.StreamId,
			&episode.ChannelId,
			&episode.Title,
			&episode.Description,
			&episode.PublishDate,
			&episode.Status,
			&episode.Parent,
			&episode.IsDir,
			&episode.Year,
			&genresString,
			&episode.CoverArt,
			&episode.Size,
			&episode.ContentType,
			&episode.Suffix,
			&episode.Duration,
			&episode.BitRate,
			&episode.Path,
			&episode.SourceUrl,
		); err != nil {
			return nil, fmt.Errorf("scanning podcast episode row: %v", err)
		}
		episode.Genres = []types.ChildGenre{}
		for _, genre := range strings.Split(genresString.String, ",") {
			episode.Genres = append(episode.Genres, types.ChildGenre{Name: genre})
		}

		episode.Genre = episode.Genres[0].Name

		episodes = append(episodes, episode)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating podcast episodes rows: %v", err)
	}

	return episodes, nil
}

func GetNewestPodcastEpisodes(ctx context.Context, count int) ([]types.PodcastEpisode, error) {
	var episodes []types.PodcastEpisode
	query := `select pe.id, pe.guid, pe.channel_id, pe.title, pe.description, pe.publish_date,
		pe.status, pe.id, 'false', pe.year, pc.categories, pe.cover_art, pe.size, pe.content_type,
		pe.suffix, pe.duration, pe.bit_rate, pe.file_path, pe.source_url
	from podcast_episodes pe
	join podcast_channels pc on pc.id = pe.channel_id
	order by pe.publish_date desc
	limit ?;`
	rows, err := DB.QueryContext(ctx, query, count)
	if err != nil {
		return nil, fmt.Errorf("querying newest podcast episodes: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var episode types.PodcastEpisode
		var genresString sql.NullString
		if err := rows.Scan(
			&episode.Id,
			&episode.StreamId,
			&episode.ChannelId,
			&episode.Title,
			&episode.Description,
			&episode.PublishDate,
			&episode.Status,
			&episode.Parent,
			&episode.IsDir,
			&episode.Year,
			&genresString,
			&episode.CoverArt,
			&episode.Size,
			&episode.ContentType,
			&episode.Suffix,
			&episode.Duration,
			&episode.BitRate,
			&episode.Path,
			&episode.SourceUrl,
		); err != nil {
			return nil, fmt.Errorf("scanning podcast episode row: %v", err)
		}
		episode.Genres = []types.ChildGenre{}
		for _, genre := range strings.Split(genresString.String, ",") {
			episode.Genres = append(episode.Genres, types.ChildGenre{Name: genre})
		}

		episode.Genre = episode.Genres[0].Name

		episodes = append(episodes, episode)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating podcast episodes rows: %v", err)
	}

	return episodes, nil
}

func GetPodcastEpisodesUserless(ctx context.Context) ([]types.PodcastEpisode, error) {
	var episodes []types.PodcastEpisode
	query := `select pe.id, pe.guid, pe.channel_id, pe.title, pe.description, pe.publish_date,
		pe.status, pe.id, 'false', pe.year, pc.categories, pe.cover_art, pe.size, pe.content_type,
		pe.suffix, pe.duration, pe.bit_rate, pe.file_path, pe.source_url
	from podcast_episodes pe
	join podcast_channels pc on pc.id = pe.channel_id
	order by pe.publish_date desc;`
	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying podcast episodes in GetPodcastEpisodesUserless: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var episode types.PodcastEpisode
		var genresString sql.NullString
		if err := rows.Scan(
			&episode.Id,
			&episode.StreamId,
			&episode.ChannelId,
			&episode.Title,
			&episode.Description,
			&episode.PublishDate,
			&episode.Status,
			&episode.Parent,
			&episode.IsDir,
			&episode.Year,
			&genresString,
			&episode.CoverArt,
			&episode.Size,
			&episode.ContentType,
			&episode.Suffix,
			&episode.Duration,
			&episode.BitRate,
			&episode.Path,
			&episode.SourceUrl,
		); err != nil {
			return nil, fmt.Errorf("scanning podcast episode row: %v", err)
		}
		episode.Genres = []types.ChildGenre{}
		for _, genre := range strings.Split(genresString.String, ",") {
			episode.Genres = append(episode.Genres, types.ChildGenre{Name: genre})
		}

		episode.Genre = episode.Genres[0].Name

		episodes = append(episodes, episode)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating podcast episodes rows: %v", err)
	}

	return episodes, nil
}

func GetPodcastsUserless(ctx context.Context, podcastId string) ([]types.PodcastChannel, error) {
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
	from podcast_channels p`

	var args []interface{}
	if podcastId != "" {
		query += ` where p.id = $1;`
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

		channelArray = append(channelArray, channel)
	}

	return channelArray, nil
}

func GetPodcastChannelById(ctx context.Context, channelId int) (types.PodcastChannel, error) {
	var channel types.PodcastChannel
	var errorMessage sql.NullString
	var lastRefreshed sql.NullString
	query := `select p.id, p.title, p.url, p.description, p.cover_art, p.original_image_url,
		p.last_refresh, p.status, p.created_at, p.error_message
	from podcast_channels p
	where p.id = $1;`
	err := DB.QueryRowContext(ctx, query, channelId).Scan(
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
	)
	if err != nil {
		return channel, fmt.Errorf("querying podcast channel by ID: %v", err)
	}

	if lastRefreshed.Valid {
		channel.LastRefresh = lastRefreshed.String
	}
	if errorMessage.Valid {
		channel.ErrorMessage = errorMessage.String
	}
	return channel, nil
}
