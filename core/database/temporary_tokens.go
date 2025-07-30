package database

import (
	"database/sql"
	"context"
	"fmt"
	"time"
	"zene/core/logger"
)

func createTemporaryTokensTable(ctx context.Context) error {
	tableName := "temporary_tokens"
	schema := `CREATE TABLE IF NOT EXISTS temporary_tokens (
		temporary_token TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		expires TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	err := createTable(ctx, tableName, schema)
	return err
}

func SaveTemporaryToken(ctx context.Context, userId int64, temporary_token string, duration time.Duration) (string, error) {
	expiresAt := time.Now().Add(duration).Format(time.RFC3339Nano)
	query := "INSERT INTO temporary_tokens (user_id, temporary_token, expires) VALUES (?, ?, ?)"
	_, err := DB.ExecContext(ctx, query, userId, temporary_token, expiresAt)
	if err != nil {
		return "", fmt.Errorf("saving temporary token: %v", err)
	}
	return expiresAt, nil
}

func ExtendTemporaryToken(ctx context.Context, userId int64, temporary_token string, duration time.Duration) (string, error) {
	expiresAt := time.Now().Add(duration).Format(time.RFC3339Nano)
	query := "UPDATE temporary_tokens SET expires = ? WHERE temporary_token = ? AND user_id = ?"
	_, err := DB.ExecContext(ctx, query, expiresAt, temporary_token, userId)
	if err != nil {
		return "", fmt.Errorf("extending temporary token: %v", err)
	}
	return expiresAt, nil
}

func IsTemporaryTokenValid(ctx context.Context, temporary_token string) (bool, error) {
	var expiresAtStr string
	query := "SELECT expires FROM temporary_tokens WHERE temporary_token = ?"
	err := DB.QueryRowContext(ctx, query, temporary_token).Scan(&expiresAtStr)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("checking temporary token validity: %v", err)
	}
	
	expiresAt, err := time.Parse(time.RFC3339Nano, expiresAtStr)
	if err != nil {
		return false, fmt.Errorf("parsing expiry time for temporary_token %s: %v", temporary_token, err)
	}
	return time.Now().Before(expiresAt), nil
}

func CleanupExpiredTemporaryTokens(ctx context.Context) {
	query := `DELETE FROM temporary_tokens WHERE expires < ?`
	_, err := DB.ExecContext(ctx, query, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		logger.Printf("Failed to run temporary_tokens cleanup: %v", err)
	}
	logger.Printf("temporary_tokens cleanup finished")
}
