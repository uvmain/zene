package database

import (
	"context"
	"fmt"
	"zene/core/logic"
)

func createUserStarsTable(ctx context.Context) {
	schema := `CREATE TABLE user_stars (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		metadata_id TEXT NOT NULL,
		created_at TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE (user_id, metadata_id)
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_user_stars_metadata_user", "user_stars", []string{"metadata_id", "user_id"}, false)
}

func UpsertUserStar(ctx context.Context, userId int, metadataId string) error {
	isValidMetadataResponse, _, err := IsValidMetadataId(ctx, metadataId)
	if !isValidMetadataResponse {
		return fmt.Errorf("invalid metadata ID: %s", metadataId)
	}

	query := `INSERT OR IGNORE INTO user_stars (user_id, metadata_id, created_at)
		VALUES (?, ?, ?);`

	createdAt := logic.GetCurrentTimeFormatted()

	_, err = DB.ExecContext(ctx, query, userId, metadataId, createdAt)
	if err != nil {
		return fmt.Errorf("upserting user star row: %v", err)
	}
	return nil
}

func DeleteUserStar(ctx context.Context, userId int, metadataId string) error {
	query := `DELETE FROM user_stars WHERE user_id = ? AND metadata_id = ?`
	_, err := DB.ExecContext(ctx, query, userId, metadataId)
	if err != nil {
		return fmt.Errorf("deleting user star row: %v", err)
	}
	return nil
}
