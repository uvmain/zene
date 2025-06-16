package database

import (
	"context"
	"fmt"
	"time"
	"zene/core/logger"
)

func CreateSessionsTable(ctx context.Context) {
	tableName := "sessions"
	schema := `CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		expires TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	createTable(ctx, tableName, schema)
}

func SaveSessionToken(ctx context.Context, userId int64, token string, duration time.Duration) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("failed to take a db conn from the pool in SaveSessionToken: %v", err)
	}
	defer DbPool.Put(conn)

	expiresAt := time.Now().Add(duration)
	stmt := conn.Prep(`INSERT INTO sessions (user_id, token, expires) VALUES ($user_id, $token, $expires)`)
	defer stmt.Finalize()
	stmt.SetInt64("$user_id", userId)
	stmt.SetText("$token", token)
	stmt.SetText("$expires", expiresAt.Format(time.RFC3339Nano))

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to save session token: %v", err)
	}

	return nil
}

func IsSessionValid(ctx context.Context, userId int, token string) (bool, error) {
	var expiresAt time.Time
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to take a db conn from the pool in IsSessionValid: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT expires FROM sessions WHERE user_id = $user_id and token = $token`)
	defer stmt.Finalize()
	stmt.SetInt64("$user_id", int64(userId))
	stmt.SetText("$token", token)

	if hasRow, err := stmt.Step(); err != nil {
		return false, fmt.Errorf("failed to take a db conn from the pool in IsSessionValid: %v", err)
	} else if !hasRow {
		return false, nil
	} else {
		expiresAt, err = time.Parse(time.RFC3339Nano, stmt.GetText("expires"))
		if err != nil {
			return false, fmt.Errorf("Error parsing session expiry time for token %s: %v", token, err)
		}
		return time.Now().Before(expiresAt), nil
	}
}

func DeleteSessionToken(ctx context.Context, token string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("failed to take a db conn from the pool in DeleteSessionToken: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`DELETE FROM sessions WHERE token = $token`)
	defer stmt.Finalize()
	stmt.SetText("$token", token)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to delete session for token %s: %v", token, err)
	}
	logger.Printf("Deleted session for token %s", token)
	return nil
}

func CleanupExpiredSessions(ctx context.Context) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		logger.Printf("failed to take a db conn from the pool in CleanupExpiredSessions: %v", err)
		return
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`DELETE FROM sessions WHERE expires < $expiry`)
	defer stmt.Finalize()
	stmt.SetText("$expiry", time.Now().Format(time.RFC3339Nano))

	_, err = stmt.Step()
	if err != nil {
		logger.Printf("failed to run session cleanup: %v", err)
	}
	logger.Printf("Session cleanup finished")
}

func DeleteAllSessionsForUserId(ctx context.Context, userId int) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("failed to take a db conn from the pool in DeleteAllUserSessions: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`DELETE FROM sessions WHERE user_id = $user_id`)
	defer stmt.Finalize()
	stmt.SetInt64("$user_id", int64(userId))

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to delete all sessions for user %d: %v", userId, err)
	}
	return nil
}

func GetUserIDFromSession(ctx context.Context, token string) (int64, bool, error) {
	var userID int64
	var expiresAt time.Time

	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return 0, false, fmt.Errorf("failed to get DB conn: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT user_id, expires FROM sessions WHERE token = $token`)
	defer stmt.Finalize()
	stmt.SetText("$token", token)

	hasRow, err := stmt.Step()
	if err != nil {
		return 0, false, fmt.Errorf("failed to query session: %v", err)
	} else if !hasRow {
		return 0, false, nil
	}

	userID = stmt.GetInt64("user_id")
	expiresAt, err = time.Parse(time.RFC3339Nano, stmt.GetText("expires"))
	if err != nil {
		return 0, false, fmt.Errorf("error parsing session expiry: %v", err)
	}

	return userID, time.Now().Before(expiresAt), nil
}
