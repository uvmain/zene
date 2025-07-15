package database

import (
	"fmt"
	"strconv"
	"strings"
)

func getUnendedMetadataWithPlaycountsSql(userId int64) string {
	return fmt.Sprintf("SELECT m.*, IFNULL(up.play_count, 0) AS user_play_count, IFNULL(gp.global_play_count, 0) AS global_play_count FROM metadata m LEFT JOIN ( "+
		"SELECT musicbrainz_track_id, play_count FROM play_counts WHERE user_id = %d ) AS up ON m.musicbrainz_track_id = up.musicbrainz_track_id LEFT JOIN ( "+
		"SELECT musicbrainz_track_id, SUM(play_count) AS global_play_count FROM play_counts GROUP BY musicbrainz_track_id ) AS gp ON m.musicbrainz_track_id = gp.musicbrainz_track_id", userId)
}

func getMetadataWithGenresSql(userId int64, genres []string, condition string, limit int64, random string) string {
	stmt := getUnendedMetadataWithPlaycountsSql(userId)
	
	if len(genres) > 0 {
		stmt += " WHERE m.file_path IN ("
		stmt += "SELECT DISTINCT g.file_path FROM genres g WHERE "
		
		var genreConditions []string
		for _, genre := range genres {
			if genre != "" {
				genreConditions = append(genreConditions, fmt.Sprintf("g.genre = '%s'", strings.ReplaceAll(genre, "'", "''")))
			}
		}
		
		if len(genreConditions) > 0 {
			if condition == "or" {
				stmt += strings.Join(genreConditions, " OR ")
			} else {
				// For "and" condition, we need to ensure all genres are present for the same file_path
				stmt += strings.Join(genreConditions, " OR ")
				stmt += fmt.Sprintf(" GROUP BY g.file_path HAVING COUNT(DISTINCT g.genre) = %d", len(genreConditions))
			}
		}
		stmt += ")"
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
