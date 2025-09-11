package database

import (
	"context"
	"fmt"
	"zene/core/logic"
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

func ClearPlayqueue(ctx context.Context) error {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("getting user by context in ClearPlayqueue")
	}

	query := `DELETE FROM playqueues WHERE user_id = ?;`

	_, err = DB.ExecContext(ctx, query, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func UpsertPlayqueue(ctx context.Context, trackIds []string, currentId string, position int, changedBy string) error {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("getting user by context in UpsertPlayqueue")
	}

	trackIdsString := ""
	for i, id := range trackIds {
		if i > 0 {
			trackIdsString += ","
		}
		trackIdsString += id
	}

	changed := logic.GetCurrentTimeFormatted()

	query := `INSERT INTO playqueues (user_id, changed, changed_by, position, track_ids, current_id)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET
			changed=excluded.changed,
			changed_by=excluded.changed_by,
			position=excluded.position,
			track_ids=excluded.track_ids,
			current_id=excluded.current_id
		WHERE playqueues.user_id=excluded.user_id;`

	_, err = DB.ExecContext(ctx, query, user.Id, changed, changedBy, position, trackIdsString, currentId)
	if err != nil {
		return err
	}
	return nil
}
