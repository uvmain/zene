package database

import (
	"context"
	"fmt"
	"zene/core/logger"
)

func getUnendedMetadataWithPlaycountsSql(userId int) string {
	return fmt.Sprintf("SELECT m.*, IFNULL(up.play_count, 0) AS user_play_count, IFNULL(gp.global_play_count, 0) AS global_play_count FROM metadata m LEFT JOIN ( "+
		"SELECT musicbrainz_track_id, play_count FROM play_counts WHERE user_id = %d ) AS up ON m.musicbrainz_track_id = up.musicbrainz_track_id LEFT JOIN ( "+
		"SELECT musicbrainz_track_id, SUM(play_count) AS global_play_count FROM play_counts GROUP BY musicbrainz_track_id ) AS gp ON m.musicbrainz_track_id = gp.musicbrainz_track_id", userId)
}

func GetMediaFilePath(ctx context.Context, mediaId string) (string, error) {
	var filePath string
	query := `select file_path from (
		select file_path
		from metadata
		where musicbrainz_track_id = ?
		union ALL
		select file_path
		from podcast_episodes
		where guid = ?
		) limit 1;`
	err := DB.QueryRowContext(ctx, query, mediaId, mediaId).Scan(&filePath)
	if err != nil {
		logger.Printf("GetMediaFilePath query failed: %v", err)
		return "", err
	}
	return filePath, nil
}
