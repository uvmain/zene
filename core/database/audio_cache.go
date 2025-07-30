package database

import (
	"database/sql"
	"context"
	"fmt"
	"time"
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


	var query = `SELECT last_accessed FROM audio_cache WHERE cache_key = ?`)
	

	// param: cache_key = cache_key

	// TODO: Query single row
	if err != nil {
		return time.Time{}, fmt.Errorf("querying audio_cache: %v", err)
	}
	if !hasRow {
		return time.Time{}, fmt.Errorf("cache_key %s not found", cache_key)
	}

	lastAccessedString := stmt.GetText("last_accessed")
	lastAccessed, err := time.Parse(time.RFC3339Nano, lastAccessedString)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing last_accessed time: %v", err)
	}

	return lastAccessed, nil
}

func SelectStaleAudioCacheEntries(ctx context.Context, olderThan time.Time) ([]string, error) {


	var query = `
		SELECT cache_key FROM audio_cache
		WHERE last_accessed < ?
	`)
	

	// param: older_than = olderThan.Format(time.RFC3339Nano)

	var staleKeys []string
	for {
		// TODO: Query single row
		if err != nil {
			return nil, fmt.Errorf("stepping through stale cache query: %v", err)
		}
		if !hasRow {
			break
		}
		staleKeys = append(staleKeys, stmt.ColumnText(0))
	}

	return staleKeys, nil
}

func UpsertAudioCacheEntry(ctx context.Context, cache_key string) error {


	var query = `
		INSERT INTO audio_cache (cache_key, last_accessed)
		VALUES (?, ?)
		ON CONFLICT(cache_key) DO UPDATE SET last_accessed = ?
	`)
	

	// param: cache_key = cache_key
	// param: lastAccessed = time.Now(.Format(time.RFC3339Nano))

	// TODO: Execute query
	if err != nil {
		return fmt.Errorf("upserting audio_cache: %v", err)
	}

	return nil
}

func DeleteAudioCacheEntry(ctx context.Context, cache_key string) error {


	var query = `DELETE FROM audio_cache WHERE cache_key = ?`)
	

	// param: cache_key = cache_key

	// TODO: Execute query
	if err != nil {
		return fmt.Errorf("deleting from audio_cache: %v", err)
	}

	return nil
}

func SelectAllAudioCacheEntries(ctx context.Context) ([]types.AudioCacheEntry, error) {


	var query = `SELECT cache_key, last_accessed FROM audio_cache`)
	

	var rows []types.AudioCacheEntry
	for {
		// TODO: Query single row
		if err != nil {
			return nil, err
		}
		if !hasRow {
			break
		}
		lastAccessedString := stmt.GetText("last_accessed")
		var lastAccessed time.Time
		if lastAccessedString != "" {
			lastAccessed, err = time.Parse(time.RFC3339Nano, lastAccessedString)
			if err != nil {
				return nil, fmt.Errorf("parsing last_accessed time: %v", err)
			}
		}

		row := types.AudioCacheEntry{
			CacheKey:     stmt.GetText("cache_key"),
			LastAccessed: lastAccessed,
		}
		rows = append(rows, row)
	}

	return rows, nil
}
