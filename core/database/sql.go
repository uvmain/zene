package database

import (
	"fmt"
	"strconv"
)

func getUnendedMetadataWithPlaycountsSql(userId int) string {
	return fmt.Sprintf("SELECT m.*, IFNULL(up.play_count, 0) AS user_play_count, IFNULL(gp.global_play_count, 0) AS global_play_count FROM metadata m LEFT JOIN ( "+
		"SELECT musicbrainz_track_id, play_count FROM play_counts WHERE user_id = %d ) AS up ON m.musicbrainz_track_id = up.musicbrainz_track_id LEFT JOIN ( "+
		"SELECT musicbrainz_track_id, SUM(play_count) AS global_play_count FROM play_counts GROUP BY musicbrainz_track_id ) AS gp ON m.musicbrainz_track_id = gp.musicbrainz_track_id", userId)
}

func getMetadataWithGenresSql(userId int, genres []string, condition string, limit int, random string) string {
	stmt := getUnendedMetadataWithPlaycountsSql(userId)
	stmt += " JOIN track_genres tg ON m.file_path = tg.file_path"
	for index, genre := range genres {
		if genre == "" {
			continue
		}
		if index != 0 {
			if condition == "or" {
				stmt += " OR "
			} else {
				stmt += " AND "
			}
		} else {
			stmt += " WHERE "
		}
		stmt += fmt.Sprintf("(tg.genre = '%s')", genre)
	}
	if random != "" {
		randomInteger, err := strconv.Atoi(random)
		if err == nil {
			stmt += fmt.Sprintf(" ORDER BY ((m.rowid * %d) %% 1000000)", randomInteger)
		}
	}
	if limit > 0 {
		return fmt.Sprintf("%s limit %d;", stmt, limit)
	} else {
		return fmt.Sprintf("%s;", stmt)
	}
}
