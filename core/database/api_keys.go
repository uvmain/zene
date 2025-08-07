package database

import (
	"context"
	"database/sql"
	"fmt"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func createApiKeysTable(ctx context.Context) {
	schema := `CREATE TABLE api_keys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		api_key TEXT NOT NULL,
		date_created TEXT NOT NULL,
		last_used TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	createTable(ctx, schema)
}

func ValidateApiKey(ctx context.Context, apiKey string) (types.User, error) {
	query := `SELECT u.* FROM users_with_folders u join api_keys k on u.user_id = k.user_id WHERE k.api_key = ?`
	var row types.User
	var foldersString string

	err := DB.QueryRowContext(ctx, query, apiKey).Scan(&row.Id, &row.Username, &row.Email, &row.Password, &row.ScrobblingEnabled, &row.LdapAuthenticated,
		&row.AdminRole, &row.SettingsRole, &row.StreamRole, &row.JukeboxRole, &row.DownloadRole, &row.UploadRole, &row.PlaylistRole,
		&row.CoverArtRole, &row.CommentRole, &row.PodcastRole, &row.ShareRole, &row.VideoConversionRole, &foldersString)
	if err == sql.ErrNoRows {
		return types.User{}, fmt.Errorf("user not found")
	} else if err != nil {
		return types.User{}, fmt.Errorf("selecting user from users in GetUserByApiKey: %v", err)
	}
	row.Folders = logic.StringToIntSlice(foldersString)
	return row, nil
}

func InsertApiKey(ctx context.Context, userId int64, apiKey string) error {
	query := `
		INSERT INTO api_keys (user_id, apiKey, dateCreated)
		VALUES (?, ?, ?)`
	dateCreated := logic.GetCurrentTimeFormatted()
	result, err := DB.ExecContext(ctx, query, userId, apiKey, dateCreated)
	if err != nil || result == nil {
		return fmt.Errorf("inserting API key: %v", err)
	}

	return nil
}

func UpdateApiKeyLastUsed(ctx context.Context, apiKey string) error {
	query := `UPDATE api_keys SET last_used = ? WHERE api_key = ?`
	_, err := DB.ExecContext(ctx, query, logic.GetCurrentTimeFormatted(), apiKey)
	if err != nil {
		return fmt.Errorf("updating API key last used: %v", err)
	}
	return nil
}

func DeleteApiKey(ctx context.Context, apiKey string) error {
	query := `DELETE FROM api_keys WHERE api_key = ?`
	_, err := DB.ExecContext(ctx, query, apiKey)
	if err != nil {
		logger.Printf("Error deleting API key: %v", err)
		return fmt.Errorf("deleting API key: %v", err)
	}
	return nil
}
