package database

import (
	"context"
	"fmt"
	"zene/core/logic"
	"zene/core/types"
)

func createPlayCountsTable(ctx context.Context) {
	schema := `CREATE TABLE play_counts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		musicbrainz_track_id TEXT NOT NULL,
		play_count INTEGER NOT NULL DEFAULT 0,
		last_played TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE (user_id, musicbrainz_track_id)
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_playcounts_user_track ", "play_counts", []string{"user_id", "musicbrainz_track_id"}, true)
	createIndex(ctx, "idx_playcounts_track ", "play_counts", []string{"musicbrainz_track_id"}, false)
	createIndex(ctx, "idx_play_counts_user", "play_counts", []string{"user_id"}, false)
}

func UpsertPlayCount(ctx context.Context, userId int, musicbrainzTrackId string) error {
	query := `INSERT INTO play_counts (user_id, musicbrainz_track_id, play_count, last_played)
		VALUES (?, ?, 1, ?)
		ON CONFLICT(user_id, musicbrainz_track_id)
		DO UPDATE SET play_count = play_count + 1, last_played = excluded.last_played`

	_, err := DB.ExecContext(ctx, query, userId, musicbrainzTrackId, logic.GetCurrentTimeFormatted())
	if err != nil {
		return fmt.Errorf("upserting playcount: %v", err)
	}
	return nil
}


