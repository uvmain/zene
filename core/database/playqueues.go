package database

import (
	"context"
)

func migratePlayqueues(ctx context.Context) {
	schema := `CREATE TABLE playqueues (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		changed TEXT NOT NULL,
		changed_by TEXT NOT NULL,
		position INTEGER NOT NULL,
		track_ids TEXT NOT NULL,
		current_id TEXT,
		UNIQUE(user_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_playqueues_user", "playqueues", []string{"user_id"}, false)
}
