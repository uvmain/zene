package database

import (
	"context"
	"database/sql"
	"fmt"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func createApiKeysTable(ctx context.Context) error {
	tableName := "api_keys"
	schema := `CREATE TABLE IF NOT EXISTS api_keys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		api_key TEXT NOT NULL,
		date_created TEXT NOT NULL,
		last_used TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	);`
	err := createTable(ctx, tableName, schema)
	return err
}

func ValidateApiKey(ctx context.Context, apiKey string) (types.User, error) {
	query := `SELECT u.Id, u.username, u.encrypted_password, u.created_at, u.is_admin, u.is_disabled FROM users u join api_keys k on u.id = k.user_id WHERE k.api_key = ?`
	var user types.User

	err := DB.QueryRowContext(ctx, query, apiKey).Scan(&user.Id, &user.Username, &user.EncryptedPassword, &user.CreatedAt, &user.IsAdmin, &user.IsDisabled)
	if err == sql.ErrNoRows {
		return types.User{}, fmt.Errorf("user not found")
	} else if user.IsDisabled {
		return types.User{}, fmt.Errorf("user is disabled")
	} else if err != nil {
		return types.User{}, fmt.Errorf("selecting user from users in GetUserByApiKey: %v", err)
	}
	return user, nil
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
