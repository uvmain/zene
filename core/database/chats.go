package database

import (
	"context"
	"fmt"
	"time"
	"zene/core/types"
)

func createChatsTable(ctx context.Context) {
	schema := `CREATE TABLE chats (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		message TEXT NOT NULL,
		timestamp INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_chats_user", "chats", []string{"user_id"}, false)
}

func InsertChat(ctx context.Context, userId int, message string) error {
	insertTimestampUnixSeconds := time.Now().UnixMilli()
	query := `INSERT INTO chats (user_id, message, timestamp)
		VALUES (?, ?, ?)`
	_, err := DB.ExecContext(ctx, query, userId, message, insertTimestampUnixSeconds)
	if err != nil {
		return fmt.Errorf("inserting chat: %v", err)
	}
	return nil
}

func GetChats(ctx context.Context, timeSince int) ([]types.Chat, error) {
	query := "SELECT u.username, c.message, c.timestamp FROM chats c join users u on c.user_id = u.id WHERE c.timestamp > ?"
	rows, err := DB.QueryContext(ctx, query, timeSince)
	if err != nil {
		return []types.Chat{}, fmt.Errorf("querying chats: %v", err)
	}
	defer rows.Close()

	var result []types.Chat
	for rows.Next() {
		var row types.Chat
		err := rows.Scan(&row.UserName, &row.Message, &row.Timestamp)
		if err != nil {
			return []types.Chat{}, fmt.Errorf("scanning chat row: %v", err)
		}
		result = append(result, row)
	}
	return result, nil
}
