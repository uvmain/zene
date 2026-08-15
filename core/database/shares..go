package database

import (
	"context"
	"fmt"
	"zene/core/logic"
)

type CreateShareOptions struct {
	Token       string
	Description string
	ExpiresAt   string
	ApiKeyId    string
	MediaIds    []string
}

func createSharesTable(ctx context.Context) {
	schema := `CREATE TABLE shares (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		owner_user_id INTEGER NOT NULL,
		token TEXT NOT NULL,
		description TEXT,
		created_at TEXT NOT NULL,
		expires_at TEXT NOT NULL,
		api_key_id TEXT NOT NULL,
		visit_count INTEGER NOT NULL DEFAULT 0,
		FOREIGN KEY (owner_user_id) REFERENCES users(id) ON DELETE CASCADE
		FOREIGN KEY (api_key_id) REFERENCES api_keys(id) ON DELETE CASCADE
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

func UpsertShare(ctx context.Context, options CreateShareOptions) error {
	owner, err := GetUserByContext(ctx)
	if err != nil {
		return err
	}
	createdDate := logic.GetCurrentTimeFormatted()

	shareQuery := `INSERT INTO shares (owner_user_id, token, description, created_at, expires_at, api_key_id)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(token) DO UPDATE SET description=excluded.description, expires_at=excluded.expires_at
		WHERE excluded.expires_at>shares.expires_at`

	result, err := DB.ExecContext(ctx, shareQuery, owner.Id, options.Token, options.Description, createdDate, options.ExpiresAt, options.ApiKeyId)
	if err != nil {
		return fmt.Errorf("inserting share row: %v", err)
	}
	lastInserted, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last inserted ID: %v", err)
	}

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %v", err)
	}

	mediaQuery := `INSERT INTO shared_media (share_id, media_id) VALUES (?, ?)`
	for _, mediaId := range options.MediaIds {
		_, err := tx.ExecContext(ctx, mediaQuery, lastInserted, mediaId)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("inserting shared media row: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("committing transaction: %v", err)
	}
	return nil
}

func ClearExpiredShares(ctx context.Context) error {
	query := `DELETE FROM shares WHERE expires_at < ?`
	now := logic.GetCurrentTimeFormatted()
	_, err := DB.ExecContext(ctx, query, now)
	if err != nil {
		return err
	}
	return nil
}
