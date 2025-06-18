package database

import "fmt"

func getMetadataWithPlaycountsSql(userId int64) string {
	return fmt.Sprintf("SELECT "+
		"m.*, "+
		"IFNULL(up.play_count, 0) AS user_play_count, "+
		"IFNULL(gp.global_playcount, 0) AS global_playcount "+
		"FROM metadata m "+
		"LEFT JOIN ( "+
		"SELECT musicbrainz_track_id, play_count "+
		"FROM play_counts "+
		"WHERE user_id = %d "+
		") AS up ON m.musicbrainz_track_id = up.musicbrainz_track_id "+
		"LEFT JOIN ( "+
		"SELECT musicbrainz_track_id, SUM(play_count) AS global_playcount "+
		"FROM play_counts "+
		"GROUP BY musicbrainz_track_id "+
		") AS gp ON m.musicbrainz_track_id = gp.musicbrainz_track_id", userId)
}
