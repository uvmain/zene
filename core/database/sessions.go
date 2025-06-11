package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"zene/core/logic"
)

func CreateSessionsTable(ctx context.Context) {
	tableName := "sessions"
	schema := `CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY,
		expires TEXT NOT NULL
	);`
	createTable(ctx, tableName, schema)
}

func SaveSessionToken(ctx context.Context, token string, duration time.Duration) (int, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	if err := logic.CheckContext(ctx); err != nil {
		return 0, err
	}
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	expiresAt := time.Now().Add(duration)
	stmt := conn.Prep(`INSERT INTO sessions (token, expires) VALUES ($token, $expires)`)
	defer stmt.Finalize()
	stmt.SetText("$token", token)
	stmt.SetText("$expires", expiresAt.Format(time.RFC3339Nano))

	if err := logic.CheckContext(ctx); err != nil {
		return 0, err
	}

	_, err = stmt.Step()
	if err != nil {
		return 0, fmt.Errorf("failed to save session token: %v", err)
	}

	rowId := int(conn.LastInsertRowID())
	return rowId, nil
}

func IsSessionValid(ctx context.Context, token string) bool {
	var expiresAt time.Time
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	if err := logic.CheckContext(ctx); err != nil {
		return false
	}
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT expires FROM sessions WHERE token = $token`)
	defer stmt.Finalize()
	stmt.SetText("$token", token)

	if err := logic.CheckContext(ctx); err != nil {
		return false
	}

	if hasRow, err := stmt.Step(); err != nil {
		return false
	} else if !hasRow {
		return false
	} else {
		expiresAt, err = time.Parse(time.RFC3339Nano, stmt.GetText("expires"))
		if err != nil {
			log.Printf("Error parsing session expiry time for token %s: %v", token, err)
			return false
		}
		return time.Now().Before(expiresAt)
	}
}

func DeleteSessionToken(ctx context.Context, token string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if err := logic.CheckContext(ctx); err != nil {
		return err
	}
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`DELETE FROM sessions WHERE token = $token`)
	defer stmt.Finalize()
	stmt.SetText("$token", token)

	if err := logic.CheckContext(ctx); err != nil {
		return err
	}

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to delete session for token %s: %v", token, err)
	}
	log.Printf("Deleted session for token %s", token)
	return nil
}

func CleanupExpiredSessions(ctx context.Context) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if err := logic.CheckContext(ctx); err != nil {
		return
	}
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`DELETE FROM sessions WHERE expires < $expiry`)
	defer stmt.Finalize()
	stmt.SetText("$expiry", time.Now().Format(time.RFC3339Nano))

	if err := logic.CheckContext(ctx); err != nil {
		return
	}

	_, err = stmt.Step()
	if err != nil {
		log.Printf("failed to run session cleanup: %v", err)
	}
	log.Printf("Session cleanup finished")
}

func StartSessionCleanupRoutine(ctx context.Context) {
	log.Println("Starting session cleanup routine")
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			CleanupExpiredSessions(ctx)
		}
	}()
}
