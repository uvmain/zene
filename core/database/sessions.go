package database

import (
	"database/sql"
	"context"
	"fmt"
	"time"
	"zene/core/logger"

	"github.com/patrickmn/go-cache"
)

var sessionCache = cache.New(5*time.Minute, 10*time.Minute)

func createSessionsTable(ctx context.Context) error {
	tableName := "sessions"
	schema := `CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		expires TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	err := createTable(ctx, tableName, schema)
	return err
}

func SaveSessionToken(ctx context.Context, userId int64, token string, duration time.Duration) error {
	expiresAt := time.Now().Add(duration)
	query := `INSERT INTO sessions (user_id, token, expires) VALUES (?, ?, ?)`
	_, err := DB.ExecContext(ctx, query, userId, token, expiresAt.Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("saving session token: %v", err)
	}
	return nil
}

func IsSessionValid(ctx context.Context, userId int, token string) (bool, error) {
	var expiresAt time.Time
	query := `SELECT expires FROM sessions WHERE user_id = ? AND token = ?`
	err := DB.QueryRowContext(ctx, query, userId, token).Scan(&expiresAt)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("checking session validity: %v", err)
	}
	return time.Now().Before(expiresAt), nil
}

func DeleteSessionToken(ctx context.Context, token string) error {
	query := `DELETE FROM sessions WHERE token = ?`
	_, err := DB.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("deleting session for token %s: %v", token, err)
	}
	logger.Printf("Deleted session for token %s", token)
	return nil
}

func CleanupExpiredSessions(ctx context.Context) {
	query := `DELETE FROM sessions WHERE expires < ?`
	_, err := DB.ExecContext(ctx, query, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		logger.Printf("Failed to run session cleanup: %v", err)
	}
	logger.Printf("Session cleanup finished")
}

func DeleteAllSessionsForUserId(ctx context.Context, userId int) error {
	query := `DELETE FROM sessions WHERE user_id = ?`
	_, err := DB.ExecContext(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("deleting all sessions for user %d: %v", userId, err)
	}
	return nil
}

type cachedSession struct {
	UserID    int64
	ExpiresAt time.Time
}

func GetUserIdFromSession(ctx context.Context, token string) (int64, bool, error) {
	if cachedVal, found := sessionCache.Get(token); found {
		entry := cachedVal.(cachedSession)
		return entry.UserID, time.Now().Before(entry.ExpiresAt), nil
	}

	var userID int64
	var expiresAtStr string
	query := `SELECT user_id, expires FROM sessions WHERE token = ?`
	err := DB.QueryRowContext(ctx, query, token).Scan(&userID, &expiresAtStr)
	if err == sql.ErrNoRows {
		return 0, false, nil
	} else if err != nil {
		return 0, false, fmt.Errorf("querying session: %v", err)
	}

	expiresAt, err := time.Parse(time.RFC3339Nano, expiresAtStr)
	if err != nil {
		return 0, false, fmt.Errorf("parsing session expiry: %v", err)
	}

	ttl := time.Until(expiresAt)
	if ttl > 0 {
		sessionCache.Set(token, cachedSession{
			UserID:    userID,
			ExpiresAt: expiresAt,
		}, ttl)
	}

	return userID, time.Now().Before(expiresAt), nil
}
