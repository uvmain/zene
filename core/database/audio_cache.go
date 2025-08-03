package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"zene/core/logger"
	"zene/core/types"
)

func createAudioCacheTable(ctx context.Context) error {
	tableName := "audio_cache"
	schema := `CREATE TABLE IF NOT EXISTS audio_cache (
		cache_key TEXT PRIMARY KEY,
		last_accessed TEXT NOT NULL
	);`
	err := createTable(ctx, tableName, schema)
	return err
}

func SelectAudioCacheEntry(ctx context.Context, cache_key string) (time.Time, error) {
	query := "SELECT last_accessed FROM audio_cache WHERE cache_key = ?"

	var lastAccessedString string

	err := DB.QueryRowContext(ctx, query, cache_key).Scan(&lastAccessedString)

	if err == sql.ErrNoRows {
		return time.Time{}, fmt.Errorf("cache_key %s not found", cache_key)
	} else if err != nil {
		return time.Time{}, fmt.Errorf("querying audio_cache: %v", err)
	}

	lastAccessed, err := time.Parse(time.RFC3339Nano, lastAccessedString)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing last_accessed time: %v", err)
	}

	return lastAccessed, nil
}

func SelectStaleAudioCacheEntries(ctx context.Context, olderThan time.Time) ([]string, error) {
	query := "SELECT cache_key FROM audio_cache WHERE last_accessed < ?"

	rows, err := DB.QueryContext(ctx, query, olderThan.Format(time.RFC3339Nano))
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []string{}, err
	}
	defer rows.Close()

	var results []string

	for rows.Next() {
		var result string
		if err := rows.Scan(&result); err != nil {
			logger.Printf("Failed to scan row in SelectStaleAudioCacheEntries: %v", err)
			return []string{}, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return results, err
	}

	return results, nil
}

func UpsertAudioCacheEntry(ctx context.Context, cache_key string) error {
	stmt := `
		INSERT INTO audio_cache (cache_key, last_accessed)
		VALUES (?, ?)
		ON CONFLICT(cache_key) DO UPDATE SET last_accessed = ?
	`

	lastAccessed := time.Now().Format(time.RFC3339Nano)
	_, err := DB.ExecContext(ctx, stmt, cache_key, lastAccessed, lastAccessed)

	if err != nil {
		return fmt.Errorf("upserting audio cache row: %v", err)
	}

	return nil
}

func DeleteAudioCacheEntry(ctx context.Context, cache_key string) error {
	stmt := "DELETE FROM audio_cache WHERE cache_key = ?"
	_, err := DB.ExecContext(ctx, stmt, cache_key)

	if err != nil {
		return fmt.Errorf("deleting audio cache row: %v", err)
	}

	return nil
}

func SelectAllAudioCacheEntries(ctx context.Context) ([]types.AudioCacheEntry, error) {
	query := "SELECT cache_key, last_accessed FROM audio_cache"

	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.AudioCacheEntry{}, err
	}
	defer rows.Close()

	var results []types.AudioCacheEntry

	for rows.Next() {
		var result types.AudioCacheEntry
		var cacheKey string
		var lastAccessedString string
		if err := rows.Scan(&cacheKey, &lastAccessedString); err != nil {
			logger.Printf("Failed to scan row in SelectAllAudioCacheEntries: %v", err)
			return []types.AudioCacheEntry{}, err
		}
		result.CacheKey = cacheKey
		var lastAccessed time.Time
		if lastAccessedString != "" {
			lastAccessed, err = time.Parse(time.RFC3339Nano, lastAccessedString)
			if err != nil {
				return nil, fmt.Errorf("parsing last_accessed time: %v", err)
			}
		}
		result.LastAccessed = lastAccessed
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return results, err
	}

	return results, nil
}
