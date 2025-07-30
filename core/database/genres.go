package database

import (
	"context"
	"fmt"
	"zene/core/types"
)

func createGenresTable(ctx context.Context) error {
	tableName := "track_genres"
	schema := `CREATE TABLE IF NOT EXISTS track_genres (
id INTEGER PRIMARY KEY AUTOINCREMENT,
file_path TEXT NOT NULL,
genre TEXT NOT NULL,
FOREIGN KEY (file_path) REFERENCES metadata(file_path) ON DELETE CASCADE
);`
	err := createTable(ctx, tableName, schema)
	if err != nil {
		return err
	}
	createIndex(ctx, "idx_track_genres_file_path", "track_genres", "file_path", false)
	createIndex(ctx, "idx_track_genres_genre", "track_genres", "genre", false)
	return nil
}

func SelectGenres(ctx context.Context, searchTerm string) ([]string, error) {
	// TODO: Migrate to standard database/sql patterns
	return []string{}, fmt.Errorf("function not yet migrated - core migration complete")
}

func UpsertGenresForTrack(ctx context.Context, filePath string, genres []string) error {
	// TODO: Migrate to standard database/sql patterns
	return fmt.Errorf("function not yet migrated - core migration complete")
}

func SelectTracksByGenres(ctx context.Context, genres []string, condition string, userId int64, limit string, random string) ([]types.MetadataWithPlaycounts, error) {
	// TODO: Migrate to standard database/sql patterns
	return []types.MetadataWithPlaycounts{}, fmt.Errorf("function not yet migrated - core migration complete")
}
