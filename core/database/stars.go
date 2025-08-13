package database

import (
	"context"
	"fmt"
)

func createUserStarsTable(ctx context.Context) {
	schema := `CREATE TABLE user_stars (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		metadata_id TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE (user_id, metadata_id)
	);`
	createTable(ctx, schema)
}

func UpsertUserStar(ctx context.Context, userId int64, metadataId string) error {
	isValidMetadataResponse, _, err := IsValidMetadataId(ctx, metadataId)
	if !isValidMetadataResponse {
		return fmt.Errorf("invalid metadata ID: %s", metadataId)
	}

	query := `INSERT OR IGNORE INTO user_stars (user_id, metadata_id)
		VALUES (?, ?);`

	_, err = DB.ExecContext(ctx, query, userId, metadataId)
	if err != nil {
		return fmt.Errorf("upserting user star row: %v", err)
	}
	return nil
}

func GetUserStarsForUser(ctx context.Context, userId int64) ([]string, error) {
	query := "SELECT metadata_id FROM user_stars WHERE user_id = ?"
	rows, err := DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("querying user stars: %v", err)
	}
	defer rows.Close()

	var metadataIds []string
	for rows.Next() {
		var metadataId string
		if err := rows.Scan(&metadataId); err != nil {
			return nil, fmt.Errorf("scanning user star row: %v", err)
		}
		metadataIds = append(metadataIds, metadataId)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user star rows: %v", err)
	}

	return metadataIds, nil
}

func DeleteUserStar(ctx context.Context, userId int64, metadataId string) error {
	query := `DELETE FROM user_stars WHERE user_id = ? AND metadata_id = ?`
	_, err := DB.ExecContext(ctx, query, userId, metadataId)
	if err != nil {
		return fmt.Errorf("deleting user star row: %v", err)
	}
	return nil
}
