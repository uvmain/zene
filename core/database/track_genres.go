package database

import (
	"context"
	"fmt"
	"log"
	"zene/core/logger"
	"zene/core/types"
)

func createTrackGenresTable(ctx context.Context) {
	schema := `CREATE TABLE track_genres (
		file_path TEXT NOT NULL,
		genre TEXT NOT NULL,
		FOREIGN KEY(file_path) REFERENCES metadata(file_path) ON DELETE CASCADE
	);`

	createTable(ctx, schema)
	createIndex(ctx, "idx_track_genres_file_path", "track_genres", []string{"file_path"}, false)
	createIndex(ctx, "idx_track_genres_genre", "track_genres", []string{"genre"}, false)
	createIndex(ctx, "idx_track_genres_genre_file_path", "track_genres", []string{"genre", "file_path"}, false)
	createIndex(ctx, "idx_track_genres_genre_file_path", "track_genres", []string{"genre", "file_path"}, false)

	createTrackGenresTriggers(ctx)

	var count int
	err = DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM track_genres").Scan(&count)
	if err != nil {
		log.Fatalf("error checking count of track_genres table: %v", err)
	}

	if count == 0 {
		var count int
		err = DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM metadata").Scan(&count)
		if err != nil {
			log.Fatalf("error checking count of metadata table: %v", err)
		}
		if count > 0 {
			log.Println("Database: track_genres table is empty, populating from metadata")
			err = populateTrackGenresFromMetadata(ctx)
			if err != nil {
				log.Fatalf("error populating track_genres table from metadata: %v", err)
			} else {
				log.Println("Database: track_genres table populated from metadata")
			}
		}
	}
}

func createTrackGenresTriggers(ctx context.Context) {
	createTrigger(ctx, `CREATE TRIGGER tr_metadata_insert_genres AFTER INSERT ON metadata
	BEGIN
		INSERT INTO track_genres (file_path, genre)
		WITH RECURSIVE split_genre(file_path, genre, rest) AS (
		SELECT 
			file_path,
			'', 
			genre || ';'  -- add trailing semicolon for parsing
		FROM metadata where file_path = new.file_path
		UNION ALL
		SELECT
			file_path,
			substr(rest, 0, instr(rest, ';')),               -- get one genre
			substr(rest, instr(rest, ';') + 1)               -- remaining string
		FROM split_genre
		WHERE rest <> ''
		)
		SELECT file_path, genre
		FROM split_genre
		WHERE genre <> '';
	END;`)

	createTrigger(ctx, `CREATE TRIGGER tr_metadata_delete_genres AFTER DELETE ON metadata
    BEGIN
        DELETE FROM track_genres WHERE file_path = old.file_path;
    END;`)

	createTrigger(ctx, `CREATE TRIGGER tr_metadata_update_genres AFTER UPDATE ON metadata
    BEGIN
		DELETE FROM track_genres WHERE file_path = old.file_path;
        INSERT INTO track_genres (file_path, genre)
		WITH RECURSIVE split_genre(file_path, genre, rest) AS (
		SELECT 
			file_path,
			'', 
			genre || ';'  -- add trailing semicolon for parsing
		FROM metadata where file_path = new.file_path

		UNION ALL

		SELECT
			file_path,
			substr(rest, 0, instr(rest, ';')),               -- get one genre
			substr(rest, instr(rest, ';') + 1)               -- remaining string
		FROM split_genre
		WHERE rest <> ''
		)
		SELECT file_path, genre
		FROM split_genre
		WHERE genre <> '';
    END;`)
}

func populateTrackGenresFromMetadata(ctx context.Context) error {
	var stmt = `INSERT INTO track_genres (file_path, genre)
		WITH RECURSIVE split_genre(file_path, genre, rest) AS (
		SELECT 
			file_path,
			'', 
			genre || ';'  -- add trailing semicolon for parsing
		FROM metadata

		UNION ALL

		SELECT
			file_path,
			substr(rest, 0, instr(rest, ';')), -- get one genre
			substr(rest, instr(rest, ';') + 1) -- remaining string
		FROM split_genre
		WHERE rest <> ''
		)
		SELECT file_path, genre
		FROM split_genre
		WHERE genre <> '';`

	_, err := DB.ExecContext(ctx, stmt)

	if err != nil {
		return fmt.Errorf("inserting data into track_genres table: %v", err)
	}

	return nil
}

func SelectDistinctGenres(ctx context.Context) ([]types.Genre, error) {
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
			logger.Printf("Failed to scan row in SelectDistinctGenres: %v", err)
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
