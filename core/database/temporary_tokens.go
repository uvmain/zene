package database

import (
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
	dbMutex.Lock()
	defer dbMutex.Unlock()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		return "", fmt.Errorf("taking a db conn from the pool in SaveTemporaryToken: %v", err)
	}
	defer DbPool.Put(conn)

	expiresAt := time.Now().Add(duration).Format(time.RFC3339Nano)
	stmt := conn.Prep("INSERT INTO temporary_tokens (user_id, temporary_token, expires) VALUES ($user_id, $temporary_token, $expires);")
	defer stmt.Finalize()
	stmt.SetInt64("$user_id", userId)
	stmt.SetText("$temporary_token", temporary_token)
	stmt.SetText("$expires", expiresAt)

	_, err = stmt.Step()
	if err != nil {
		return "", fmt.Errorf("saving temporary token: %v", err)
	}

	return expiresAt, nil
}

func ExtendTemporaryToken(ctx context.Context, userId int64, temporary_token string, duration time.Duration) (string, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		return "", fmt.Errorf("taking a db conn from the pool in ExtendTemporaryToken: %v", err)
	}
	defer DbPool.Put(conn)

	expiresAt := time.Now().Add(duration).Format(time.RFC3339Nano)
	stmt := conn.Prep("Update temporary_tokens set expires = $expires where temporary_token = $temporary_token and user_id = $user_id;")
	defer stmt.Finalize()
	stmt.SetText("$temporary_token", temporary_token)
	stmt.SetInt64("$user_id", userId)
	stmt.SetText("$expires", expiresAt)

	_, err = stmt.Step()
	if err != nil {
		return "", fmt.Errorf("extending temporary token: %v", err)
	}

	return expiresAt, nil
}

func IsTemporaryTokenValid(ctx context.Context, temporary_token string) (bool, error) {
	var expiresAt time.Time
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return false, fmt.Errorf("taking a db conn from the pool in IsTemporaryTokenValid: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep("SELECT expires FROM temporary_tokens WHERE temporary_token = $temporary_token;")
	defer stmt.Finalize()
	stmt.SetText("$temporary_token", temporary_token)

	if hasRow, err := stmt.Step(); err != nil {
		return false, fmt.Errorf("taking a db conn from the pool in IsTemporaryTokenValid: %v", err)
	} else if !hasRow {
		return false, nil
	} else {
		expiresAt, err = time.Parse(time.RFC3339Nano, stmt.GetText("expires"))
		if err != nil {
			return false, fmt.Errorf("parsing expiry time for temporary_token %s: %v", temporary_token, err)
		}
		return time.Now().Before(expiresAt), nil
	}
}

func CleanupExpiredTemporaryTokens(ctx context.Context) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		logger.Printf("taking a db conn from the pool in CleanupExpiredTemporaryTokens: %v", err)
		return
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`DELETE FROM temporary_tokens WHERE expires < $expiry`)
	defer stmt.Finalize()
	stmt.SetText("$expiry", time.Now().Format(time.RFC3339Nano))

	_, err = stmt.Step()
	if err != nil {
		logger.Printf("Failed to run temporary_tokens cleanup: %v", err)
	}
	logger.Printf("temporary_tokens cleanup finished")
}
