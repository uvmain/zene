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
		&row.CoverArtRole, &row.CommentRole, &row.PodcastRole, &row.ShareRole, &row.VideoConversionRole, &row.MaxBitRate, &foldersString)
	if err == sql.ErrNoRows {
		return types.User{}, fmt.Errorf("user not found")
	} else if err != nil {
		return types.User{}, fmt.Errorf("selecting user from users in GetUserByApiKey: %v", err)
	}
	row.Folders = logic.StringToIntSlice(foldersString)
	return row, nil
}

func GetApiKeys(ctx context.Context, userId int) ([]types.ApiKey, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting user from context: %v", err)
	}

	if user.Id != userId && !user.AdminRole {
		return nil, fmt.Errorf("user not authorized to access these API keys")
	}

	var lastUsed sql.NullString

	query := `SELECT id, user_id, api_key, date_created, last_used FROM api_keys`

	var args []interface{}

	if userId != 0 {
		query += ` WHERE user_id = ?`
		args = append(args, userId)
	} else {
		query += ` WHERE user_id = ?`
		args = append(args, user.Id)
	}

	query += ` ORDER BY last_used desc, date_created DESC`

	rows, err := DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying API keys: %v", err)
	}
	defer rows.Close()

	var apiKeys []types.ApiKey
	for rows.Next() {
		var apiKey types.ApiKey
		if err := rows.Scan(&apiKey.Id, &apiKey.UserId, &apiKey.ApiKey, &apiKey.DateCreated, &lastUsed); err != nil {
			return nil, fmt.Errorf("scanning API key row: %v", err)
		}
		if lastUsed.Valid {
			apiKey.LastUsed = lastUsed.String
		}
		apiKeys = append(apiKeys, apiKey)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating API key rows: %v", err)
	}

	return apiKeys, nil
}

func CreateApiKey(ctx context.Context, userId int) (types.ApiKey, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return types.ApiKey{}, fmt.Errorf("getting user from context: %v", err)
	}

	if user.Id != userId && !user.AdminRole {
		return types.ApiKey{}, fmt.Errorf("user not authorized to create API keys for other users")
	}

	newApiKey, err := logic.GenerateNewApiKey()
	if err != nil {
		return types.ApiKey{}, fmt.Errorf("generating new API key: %v", err)
	}

	query := `INSERT INTO api_keys (user_id, api_key, date_created)
		VALUES (?, ?, ?)`

	dateCreated := logic.GetCurrentTimeFormatted()

	result, err := DB.ExecContext(ctx, query, userId, newApiKey, dateCreated)
	if err != nil {
		return types.ApiKey{}, fmt.Errorf("inserting API key: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return types.ApiKey{}, fmt.Errorf("getting last insert ID: %v", err)
	}

	return types.ApiKey{
		Id:          int(id),
		UserId:      userId,
		ApiKey:      newApiKey,
		DateCreated: dateCreated,
	}, nil
}

func UpdateApiKeyLastUsed(ctx context.Context, apiKey string) error {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("getting user from context: %v", err)
	}

	var args []interface{}
	query := `UPDATE api_keys SET last_used = ? WHERE api_key = ?  AND user_id = ?`
	args = append(args, logic.GetCurrentTimeFormatted(), apiKey, user.Id)

	_, err = DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("updating API key last used: %v", err)
	}
	return nil
}

func DeleteApiKey(ctx context.Context, apiKeyId int, userId int) error {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("getting user from context: %v", err)
	}

	if user.Id != userId && !user.AdminRole {
		return fmt.Errorf("user not authorized to delete this API key")
	}

	var args []interface{}
	query := `DELETE FROM api_keys WHERE id = ?`
	args = append(args, apiKeyId)

	if user.Id != 0 && user.AdminRole {
		query += ` AND user_id = ?`
		args = append(args, userId)
	} else {
		query += ` AND user_id = ?`
		args = append(args, user.Id)
	}

	_, err = DB.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Printf("Error deleting API key: %v", err)
		return fmt.Errorf("deleting API key: %v", err)
	}
	return nil
}
