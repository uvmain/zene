package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"zene/core/logger"
	"zene/core/types"
)

func SelectAlbumArtistsForMusicDir(ctx context.Context, musicDir string, searchParam string, random string, recent string, chronological string, limit string, offset string) ([]types.ArtistResponse, error) {
	query := "select distinct m.album_artist, m.musicbrainz_artist_id FROM metadata m join music_folders mf on m.music_folder_id = mf.id and mf.name = ?"
	if searchParam != "" {
		query += " JOIN artists_fts f ON m.file_path = f.file_path"
	}
	query += " where m.album_artist = m.artist"

	if searchParam != "" {
		query += " and artists_fts MATCH ?"
	}

	if recent == "true" {
		query += " ORDER BY m.date_added desc"
	} else if random != "" {
		randomInt, err := strconv.Atoi(random)
		if err == nil {
			query += fmt.Sprintf(" ORDER BY ((m.rowid * %d) %% 1000000)", randomInt)
		}
	} else if chronological == "true" {
		query += " ORDER BY m.release_date desc"
	}

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err == nil {
			query += fmt.Sprintf(" limit %d", limitInt)
		}
	}
	if offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err == nil {
			query += fmt.Sprintf(" offset %d", offsetInt)
		}
	}

	query += ";"

	var rows *sql.Rows

	if searchParam != "" {
		rows, err = DB.QueryContext(ctx, query, musicDir, searchParam)
	} else {
		rows, err = DB.QueryContext(ctx, query, musicDir)
	}

	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.ArtistResponse{}, err
	}
	defer rows.Close()

	var results []types.ArtistResponse

	for rows.Next() {
		var result types.ArtistResponse
		if err := rows.Scan(&result.Artist, &result.MusicBrainzArtistID); err != nil {
			logger.Printf("Failed to scan row in SelectAlbumArtistsForMusicDir: %v", err)
			return []types.ArtistResponse{}, err
		}
		result.ImageURL = fmt.Sprintf("/api/artists/%s/art", result.MusicBrainzArtistID)
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return results, err
	}

	return results, nil
}
