package database

import (
	"context"
	"fmt"
	"zene/core/types"
)

func createAudioCacheTable(ctx context.Context) error {
	tableName := "audio_cache"
	schema := `CREATE TABLE IF NOT EXISTS audio_cache (
id INTEGER PRIMARY KEY AUTOINCREMENT,
musicbrainz_track_id TEXT NOT NULL UNIQUE,
file_path TEXT NOT NULL,
format TEXT NOT NULL,
quality TEXT NOT NULL,
size INTEGER NOT NULL,
date_created TEXT NOT NULL
);`
	err := createTable(ctx, tableName, schema)
	if err != nil {
		return err
	}
	createIndex(ctx, "idx_audio_cache_track_id", "audio_cache", "musicbrainz_track_id", false)
	return nil
}

func SelectCachedFile(ctx context.Context, musicbrainzTrackId string, format string, quality string) (types.AudioCacheEntry, error) {
	// TODO: Migrate to standard database/sql patterns
	return types.AudioCacheEntry{}, fmt.Errorf("function not yet migrated - core migration complete")
}

func GetAudioCacheEntry(ctx context.Context, format string) (types.AudioCacheEntry, error) {
	// TODO: Migrate to standard database/sql patterns
	return types.AudioCacheEntry{}, fmt.Errorf("function not yet migrated - core migration complete")
}

func InsertCachedFile(ctx context.Context, musicbrainzTrackId string, filePath string, format string, quality string, size int64) error {
	// TODO: Migrate to standard database/sql patterns
	return fmt.Errorf("function not yet migrated - core migration complete")
}

func DeleteCachedFile(ctx context.Context, musicbrainzTrackId string, format string, quality string) error {
	// TODO: Migrate to standard database/sql patterns
	return fmt.Errorf("function not yet migrated - core migration complete")
}

func GetOldestCachedFiles(ctx context.Context, format string, limit int) ([]types.AudioCacheEntry, error) {
	// TODO: Migrate to standard database/sql patterns
	return []types.AudioCacheEntry{}, fmt.Errorf("function not yet migrated - core migration complete")
}
