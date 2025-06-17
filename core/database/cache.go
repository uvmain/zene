package database

import (
	"context"
	"fmt"
	"time"
)

func createAudioCacheTable(ctx context.Context) {
	tableName := "audio_cache"
	schema := `CREATE TABLE IF NOT EXISTS audio_cache (
		cacheKey TEXT PRIMARY KEY,
		last_accessed TEXT NOT NULL
	);`
	createTable(ctx, tableName, schema)
}

func SelectAudioCacheEntry(ctx context.Context, cacheKey string) (time.Time, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return time.Time{}, fmt.Errorf("Failed to take a db conn from the pool in SelectAudioCacheEntry: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT last_accessed FROM audio_cache WHERE cacheKey = $cacheKey`)
	defer stmt.Finalize()

	stmt.SetText("$cacheKey", cacheKey)

	hasRow, err := stmt.Step()
	if err != nil {
		return time.Time{}, fmt.Errorf("Failed to query audio_cache: %v", err)
	}
	if !hasRow {
		return time.Time{}, fmt.Errorf("cacheKey %s not found", cacheKey)
	}

	lastAccessedString := stmt.GetText("last_accessed")
	lastAccessed, err := time.Parse(time.RFC3339Nano, lastAccessedString)
	if err != nil {
		return time.Time{}, fmt.Errorf("Failed to parse last_accessed time: %v", err)
	}

	return lastAccessed, nil
}

func SelectStaleAudioCacheEntries(ctx context.Context, olderThan time.Time) ([]string, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to take a db conn from the pool in SelectStaleAudioCacheEntries: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`
		SELECT cacheKey FROM audio_cache
		WHERE last_accessed < $older_than
	`)
	defer stmt.Finalize()

	stmt.SetText("$older_than", olderThan.Format(time.RFC3339Nano))

	var staleKeys []string
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return nil, fmt.Errorf("Error stepping through stale cache query: %v", err)
		}
		if !hasRow {
			break
		}
		staleKeys = append(staleKeys, stmt.ColumnText(0))
	}

	return staleKeys, nil
}

func UpsertAudioCacheEntry(ctx context.Context, cacheKey string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("Failed to take a db conn from the pool in UpsertAudioCacheEntry: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`
		INSERT INTO audio_cache (cacheKey, last_accessed)
		VALUES ($cacheKey, $lastAccessed)
		ON CONFLICT(cacheKey) DO UPDATE SET last_accessed = $lastAccessed
	`)
	defer stmt.Finalize()

	stmt.SetText("$cacheKey", cacheKey)
	stmt.SetText("$lastAccessed", time.Now().Format(time.RFC3339Nano))

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("Failed to upsert audio_cache: %v", err)
	}

	return nil
}

func DeleteAudioCacheEntry(ctx context.Context, cacheKey string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("Failed to take a db conn from the pool in DeleteAudioCacheEntry: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`DELETE FROM audio_cache WHERE cacheKey = $cacheKey`)
	defer stmt.Finalize()

	stmt.SetText("$cacheKey", cacheKey)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("Failed to delete from audio_cache: %v", err)
	}

	return nil
}
