package database

import (
	"context"
	"database/sql"
	"fmt"
	"zene/core/logger"
	"zene/core/types"
)

func createGenreCountsTable(ctx context.Context) {
	schema := `CREATE TABLE genre_counts (
		genre TEXT NOT NULL,
		song_count INTEGER NOT NULL,
		album_count INTEGER NOT NULL
	);`

	createTable(ctx, schema)
	createIndex(ctx, "idx_genre_counts_genre", "genre_counts", []string{"genre"}, false)
	createIndex(ctx, "idx_genre_counts_song_order", "genre_counts", []string{"song_count desc"}, false)
}

func RepopulateGenreCountsTable(ctx context.Context) error {
	logger.Printf("Repopulating genre_counts table")
	// clear the existing genre_counts contents
	_, err := DB.ExecContext(ctx, "DELETE FROM genre_counts")
	if err != nil {
		return fmt.Errorf("clearing genre_counts table: %v", err)
	}

	var stmt = `INSERT INTO genre_counts (genre, song_count, album_count)
		select genre, count(musicbrainz_track_id) as song_count, count(distinct musicbrainz_album_id) as album_count
		from (
		WITH RECURSIVE split_genre(file_path, musicbrainz_album_id, musicbrainz_track_id, genre, rest) AS (
		SELECT 
			file_path, musicbrainz_album_id, musicbrainz_track_id,
			'', 
			genre || ';'  -- add trailing semicolon for parsing
		FROM metadata
		UNION ALL
		SELECT
			file_path, musicbrainz_album_id, musicbrainz_track_id,
			substr(rest, 0, instr(rest, ';')),               -- get one genre
			substr(rest, instr(rest, ';') + 1)               -- remaining string
		FROM split_genre
		WHERE rest <> ''
		)
		SELECT file_path, musicbrainz_album_id, musicbrainz_track_id, genre
		FROM split_genre
		WHERE genre <> ''
		)
		group by genre;`

	_, err = DB.ExecContext(ctx, stmt)

	if err != nil {
		return fmt.Errorf("inserting data into genre_counts table: %v", err)
	}

	return nil
}

func SelectGenreCounts(ctx context.Context) ([]types.Genre, error) {
	query := `select genre, song_count, album_count
		from genre_counts
		order by song_count desc`

	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.Genre{}, err
	}
	defer rows.Close()

	var results []types.Genre

	for rows.Next() {
		var result types.Genre
		if err := rows.Scan(&result.Value, &result.SongCount, &result.AlbumCount); err != nil {
			logger.Printf("Failed to scan row in SelectGenreCounts: %v", err)
			return nil, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return results, err
	}

	return results, nil
}

func SelectGenreCount(ctx context.Context, genre string) (types.Genre, error) {
	query := "SELECT genre, song_count, album_count FROM genre_counts WHERE genre = ?"
	var result types.Genre
	err := DB.QueryRowContext(ctx, query, genre).Scan(&result.Value, &result.SongCount, &result.AlbumCount)
	if err == sql.ErrNoRows {
		logger.Printf("No genre found for %s", genre)
		return types.Genre{}, nil
	} else if err != nil {
		logger.Printf("Error querying genre for %s: %v", genre, err)
		return types.Genre{}, err
	}

	return result, nil
}
