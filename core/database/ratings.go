package database

import (
	"context"
	"fmt"
)

type UserRating struct {
	MetadataId string
	Rating     int
}

func createUserRatingsTable(ctx context.Context) {
	schema := `CREATE TABLE user_ratings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		metadata_id TEXT NOT NULL,
		rating INTEGER NOT NULL CHECK (rating >= 0 AND rating <= 5),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE (user_id, metadata_id)
	);`
	createTable(ctx, schema)
}

func UpsertUserRating(ctx context.Context, userId int64, metadataId string, rating int64) error {
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

	query := `INSERT OR REPLACE INTO user_ratings (user_id, metadata_id, rating)
		VALUES (?, ?, ?);`

	_, err = DB.ExecContext(ctx, query, userId, metadataId, rating)
	if err != nil {
		return fmt.Errorf("upserting user rating row in UpsertUserRating: %v", err)
	}
	return nil
}

func GetUserRatingsForUser(ctx context.Context, userId int64) (UserRating, error) {
	query := "SELECT metadata_id, rating FROM user_ratings WHERE user_id = ?"
	rows, err := DB.QueryContext(ctx, query, userId)
	if err != nil {
		return UserRating{}, fmt.Errorf("querying user ratings in GetUserRatingsForUser: %v", err)
	}
	defer rows.Close()

	var userRatings []UserRating
	for rows.Next() {
		var rating UserRating
		if err := rows.Scan(&rating.MetadataId, &rating.Rating); err != nil {
			return UserRating{}, fmt.Errorf("scanning user rating row in GetUserRatingsForUser: %v", err)
		}
		userRatings = append(userRatings, rating)
	}
	if err := rows.Err(); err != nil {
		return UserRating{}, fmt.Errorf("error iterating user rating rows in GetUserRatingsForUser: %v", err)
	}

	if len(userRatings) == 0 {
		return UserRating{}, nil
	}
	return userRatings[0], nil
}

func GetUserRatingsForMetadataId(ctx context.Context, userId int64, metadataId string) (UserRating, error) {
	isValidMetadataResponse, _, err := IsValidMetadataId(ctx, metadataId)
	if !isValidMetadataResponse {
		return UserRating{}, fmt.Errorf("invalid metadata ID: %s", metadataId)
	}

	query := "SELECT metadata_id, rating FROM user_ratings WHERE user_id = ? AND metadata_id = ?"
	rows, err := DB.QueryContext(ctx, query, userId, metadataId)
	if err != nil {
		return UserRating{}, fmt.Errorf("querying user ratings in GetUserRatingsForMetadataId: %v", err)
	}
	defer rows.Close()

	var userRatings []UserRating
	for rows.Next() {
		var rating UserRating
		if err := rows.Scan(&rating.MetadataId, &rating.Rating); err != nil {
			return UserRating{}, fmt.Errorf("scanning user rating row in GetUserRatingsForMetadataId: %v", err)
		}
		userRatings = append(userRatings, rating)
	}
	if err := rows.Err(); err != nil {
		return UserRating{}, fmt.Errorf("error iterating user rating rows in GetUserRatingsForMetadataId: %v", err)
	}

	if len(userRatings) == 0 {
		return UserRating{}, nil
	}
	return userRatings[0], nil
}

func DeleteUserRating(ctx context.Context, userId int64, metadataId string) error {
	query := `DELETE FROM user_ratings WHERE user_id = ? AND metadata_id = ?`
	_, err := DB.ExecContext(ctx, query, userId, metadataId)
	if err != nil {
		return fmt.Errorf("deleting user rating row in DeleteUserRating: %v", err)
	}
	return nil
}
