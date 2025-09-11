package database

import (
	"context"
	"fmt"
)

type UserRating struct {
	MetadataId string
	Rating     int
}

func migrateUserRatings(ctx context.Context) {
	schema := `CREATE TABLE user_ratings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		metadata_id TEXT NOT NULL,
		rating INTEGER NOT NULL CHECK (rating >= 0 AND rating <= 5),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE (user_id, metadata_id)
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_user_ratings_metadata_user", "user_ratings", []string{"metadata_id", "user_id"}, false)
}

func UpsertUserRating(ctx context.Context, userId int, metadataId string, rating int) error {
	if rating < 0 || rating > 5 {
		return fmt.Errorf("invalid rating: %d must be between 0 and 5", rating)
	}

	if rating == 0 {
		return DeleteUserRating(ctx, userId, metadataId)
	}

	isValidMetadataResponse, _, err := IsValidMetadataId(ctx, metadataId)
	if !isValidMetadataResponse {
		return fmt.Errorf("invalid metadata ID: %s", metadataId)
	}
	if err != nil {
		return err
	}

	query := `INSERT OR REPLACE INTO user_ratings (user_id, metadata_id, rating)
		VALUES (?, ?, ?);`

	_, err = DB.ExecContext(ctx, query, userId, metadataId, rating)
	if err != nil {
		return fmt.Errorf("upserting user rating row in UpsertUserRating: %v", err)
	}
	return nil
}

func DeleteUserRating(ctx context.Context, userId int, metadataId string) error {
	query := `DELETE FROM user_ratings WHERE user_id = ? AND metadata_id = ?`
	_, err := DB.ExecContext(ctx, query, userId, metadataId)
	if err != nil {
		return fmt.Errorf("deleting user rating row in DeleteUserRating: %v", err)
	}
	return nil
}
