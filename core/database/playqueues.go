package database

import (
	"context"
	"fmt"
	"zene/core/logic"
	"zene/core/types"
)

func migratePlayqueues(ctx context.Context) {
	schema := `CREATE TABLE playqueues (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		changed TEXT NOT NULL,
		changed_by TEXT NOT NULL,
		position INTEGER NOT NULL,
		track_ids TEXT NOT NULL,
		current_index INTEGER DEFAULT 0,
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

	_, err = DbWrite.ExecContext(ctx, query, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func UpsertPlayqueue(ctx context.Context, trackIds []string, currentIndex int, position int, changedBy string) error {
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

	query := `INSERT INTO playqueues (user_id, changed, changed_by, position, track_ids, current_index)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET
			changed=excluded.changed,
			changed_by=excluded.changed_by,
			position=excluded.position,
			track_ids=excluded.track_ids,
			current_index=excluded.current_index
		WHERE playqueues.user_id=excluded.user_id;`

	_, err = DbWrite.ExecContext(ctx, query, user.Id, changed, changedBy, position, trackIdsString, currentIndex)
	if err != nil {
		return err
	}
	return nil
}

func GetPlayqueue(ctx context.Context) (types.PlayqueueRowParsed, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return types.PlayqueueRowParsed{}, fmt.Errorf("getting user by context in GetPlayqueue")
	}

	query := `select u.username, p.changed as changed_date,
		p.changed_by as changed_by,
		p.position,
		p.track_ids,
		p.current_index
		from playqueues p
		join users u on u.id = p.user_id`

	row := DbRead.QueryRowContext(ctx, query+" where p.user_id = ?;", user.Id)

	var username, changed, changedBy, trackIds string
	var position, currentIndex int
	err = row.Scan(&username, &changed, &changedBy, &position, &trackIds, &currentIndex)
	if err != nil {
		return types.PlayqueueRowParsed{}, err
	}

	trackIdArray := []string{}
	if trackIds != "" {
		trackIdArray = append(trackIdArray, logic.StringToArray(trackIds, ",")...)
	}

	return types.PlayqueueRowParsed{
		Username:     username,
		Changed:      changed,
		ChangedBy:    changedBy,
		Position:     position,
		TrackIds:     trackIdArray,
		CurrentIndex: currentIndex,
	}, nil
}
