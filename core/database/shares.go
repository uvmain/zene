package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"zene/core/logic"
)

type CreateShareOptions struct {
	Description string
	ExpiresAt   time.Time
	MediaIds    []string
}

type Share struct {
	Id          int
	OwnerUserId int
	Token       string
	Description string
	CreatedAt   string
	ExpiresAt   string
	VisitCount  int
	MediaIds    []string
}

func createSharesTable(ctx context.Context) {
	schema := `CREATE TABLE shares (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		owner_user_id INTEGER NOT NULL,
		token TEXT NOT NULL,
		description TEXT,
		created_at TEXT NOT NULL,
		expires_at TEXT,
		visit_count INTEGER NOT NULL DEFAULT 0,
		FOREIGN KEY (owner_user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(token)
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_shares_owner_user_id", "shares", []string{"owner_user_id"}, false)
}

func createSharedMediaTable(ctx context.Context) {
	schema := `CREATE TABLE shared_media (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		share_id INTEGER NOT NULL,
		media_id TEXT NOT NULL,
		FOREIGN KEY (share_id) REFERENCES shares(id) ON DELETE CASCADE
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_shared_media_share_id", "shared_media", []string{"share_id"}, false)
}

func CreateShare(ctx context.Context, options CreateShareOptions) (int, error) {
	owner, err := GetUserByContext(ctx)
	if err != nil {
		return 0, err
	}

	createdDate := logic.GetCurrentTimeFormatted()
	token, err := logic.GenerateNewApiKey()
	if err != nil {
		return 0, fmt.Errorf("generating new token: %v", err)
	}

	if !options.ExpiresAt.IsZero() && logic.GetTimeFromString(createdDate).After(options.ExpiresAt) {
		return 0, fmt.Errorf("expires_at must be in the future")
	}

	expiresAtForDb := sql.NullString{}
	if !options.ExpiresAt.IsZero() {
		expiresAtForDb = sql.NullString{
			String: logic.FormatTimeAsString(options.ExpiresAt),
			Valid:  true,
		}
	}

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("starting transaction: %v", err)
	}

	shareQuery := `INSERT INTO shares (owner_user_id, token, description, created_at, expires_at)
		VALUES (?, ?, ?, ?, ?)`

	result, err := tx.ExecContext(ctx, shareQuery, owner.Id, token, options.Description, createdDate, expiresAtForDb)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("inserting share row: %v", err)
	}

	lastInserted, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("getting last inserted ID: %v", err)
	}

	mediaQuery := `INSERT INTO shared_media (share_id, media_id) VALUES (?, ?)`
	for _, mediaId := range options.MediaIds {
		_, err := tx.ExecContext(ctx, mediaQuery, lastInserted, mediaId)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("inserting shared media row: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("committing transaction: %v", err)
	}

	return int(lastInserted), nil
}

func DeleteShare(ctx context.Context, id int) error {
	owner, err := GetUserByContext(ctx)
	if err != nil {
		return err
	}

	query := `DELETE FROM shares WHERE id = ? AND owner_user_id = ?`
	result, err := DB.ExecContext(ctx, query, id, owner.Id)
	if err != nil {
		return fmt.Errorf("deleting share: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no share found with id %d for the current user", id)
	}

	return nil
}

func GetSharesByUser(ctx context.Context) ([]Share, error) {
	owner, err := GetUserByContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, owner_user_id, token, description, created_at, expires_at, visit_count FROM shares WHERE owner_user_id = ?`
	rows, err := DB.QueryContext(ctx, query, owner.Id)
	if err != nil {
		return nil, fmt.Errorf("querying shares: %v", err)
	}
	defer rows.Close()

	var shares []Share
	for rows.Next() {
		var share Share
		var expiresString sql.NullString
		if err := rows.Scan(&share.Id, &share.OwnerUserId, &share.Token, &share.Description, &share.CreatedAt, &expiresString, &share.VisitCount); err != nil {
			return nil, fmt.Errorf("scanning share: %v", err)
		}
		if expiresString.Valid {
			share.ExpiresAt = expiresString.String
		}
		shares = append(shares, share)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over shares: %v", err)
	}

	for index, share := range shares {
		mediaQuery := `SELECT media_id FROM shared_media WHERE share_id = ?`
		mediaRows, err := DB.QueryContext(ctx, mediaQuery, share.Id)
		if err != nil {
			if mediaRows != nil {
				mediaRows.Close()
			}
			return nil, fmt.Errorf("querying shared media: %v", err)
		}

		var mediaIds []string
		for mediaRows.Next() {
			var mediaId string
			if err := mediaRows.Scan(&mediaId); err != nil {
				mediaRows.Close()
				return nil, fmt.Errorf("scanning shared media: %v", err)
			}
			mediaIds = append(mediaIds, mediaId)
		}
		if err := mediaRows.Err(); err != nil {
			mediaRows.Close()
			return nil, fmt.Errorf("iterating over shared media: %v", err)
		}
		mediaRows.Close()

		shares[index].MediaIds = mediaIds
	}

	return shares, nil
}

func ClearExpiredShares(ctx context.Context) error {
	query := `DELETE FROM shares WHERE expires_at IS NOT NULL AND expires_at < ?`
	now := logic.GetCurrentTimeFormatted()
	_, err := DB.ExecContext(ctx, query, now)
	if err != nil {
		return err
	}
	return nil
}
