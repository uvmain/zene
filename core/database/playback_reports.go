package database

import (
	"context"
	"database/sql"
	"time"
	"zene/core/logger"
)

func createPlaybackReportTable(ctx context.Context) {
	schema := `CREATE TABLE playback_reports (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		now_playing_id INTEGER NOT NULL,
		state TEXT NOT NULL,
		position_ms INTEGER NOT NULL,
		playback_rate TEXT NOT NULL,
		updated_at INTEGER NOT NULL,
		FOREIGN KEY (now_playing_id) REFERENCES now_playing(id) ON DELETE CASCADE,
		UNIQUE(now_playing_id)
	);`
	createTable(ctx, schema)
}

func UpsertPlaybackReport(ctx context.Context, nowPlayingId int, state string, positionMs int, playbackRate string, updatedAt int) error {
	query := `INSERT INTO playback_reports (now_playing_id, state, position_ms, playback_rate, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(now_playing_id) DO UPDATE SET state=excluded.state, position_ms=excluded.position_ms, playback_rate=excluded.playback_rate, updated_at=excluded.updated_at`
	_, err := DB.ExecContext(ctx, query, nowPlayingId, state, positionMs, playbackRate, updatedAt)
	if err != nil {
		return err
	}
	return nil
}

func GetPlaybackReportState(ctx context.Context, nowPlayingId int) (string, bool, error) {
	query := `SELECT state FROM playback_reports WHERE now_playing_id = ? LIMIT 1`
	var state string
	err := DB.QueryRowContext(ctx, query, nowPlayingId).Scan(&state)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}
	return state, true, nil
}

func CleanupPlaybackReports(ctx context.Context) error {
	thirtyMinutesAgo := time.Now().Add(-30 * time.Minute).UnixMilli()
	query := `DELETE FROM playback_reports WHERE updated_at < ?`
	_, err := DB.ExecContext(ctx, query, thirtyMinutesAgo)
	if err != nil {
		logger.Printf("Error cleaning up playback reports: %v", err)
		return err
	}
	return nil
}
