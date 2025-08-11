package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func createGenresTable(ctx context.Context) {
	schema := `CREATE TABLE track_genres (
		file_path TEXT NOT NULL,
		genre TEXT NOT NULL,
		FOREIGN KEY(file_path) REFERENCES metadata(file_path) ON DELETE CASCADE
	);`

	createTable(ctx, schema)
	createIndex(ctx, "idx_track_genres_genre", "track_genres", "genre", false)
	createGenresTriggers(ctx)

	var count int
	err = DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM track_genres").Scan(&count)
	if err != nil {
		log.Fatalf("error checking count of track_genres table: %v", err)
	}

	if count == 0 {
		log.Println("Database: track_genres table is empty, populating from metadata")
		err = populateGenresFromMetadata(ctx)
		if err != nil {
			log.Fatalf("error populating track_genres table from metadata: %v", err)
		} else {
			log.Println("Database: track_genres table populated from metadata")
		}
	}
}

func createGenresTriggers(ctx context.Context) {
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

func populateGenresFromMetadata(ctx context.Context) error {
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

func SelectDistinctGenres(ctx context.Context, limitParam string, searchParam string) ([]types.GenreResponse, error) {
	query := "select genre, count(file_path) as count from track_genres group by genre order by count desc"

	if limitParam != "" {
		limitInt, err := strconv.Atoi(limitParam)
		if err != nil {
			return []types.GenreResponse{}, fmt.Errorf("invalid limit value: %v", err)
		}
		query += fmt.Sprintf(" limit %d", limitInt)
	}

	query += ";"

	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.GenreResponse{}, err
	}
	defer rows.Close()

	var results []types.GenreResponse

	for rows.Next() {
		var result types.GenreResponse
		if err := rows.Scan(&result.Genre, &result.Count); err != nil {
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

func SelectTracksByGenres(ctx context.Context, genres []string, andOr string, limit int64, random string) ([]types.MetadataWithPlaycounts, error) {
	userId, _ := logic.GetUserIdFromContext(ctx)
	query := getMetadataWithGenresSql(userId, genres, andOr, limit, random)

	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.MetadataWithPlaycounts{}, err
	}
	defer rows.Close()

	var results []types.MetadataWithPlaycounts

	for rows.Next() {
		var result types.MetadataWithPlaycounts
		if err := rows.Scan(&result.FilePath, &result.DateAdded, &result.DateModified, &result.FileName, &result.Format, &result.Duration,
			&result.Size, &result.Bitrate, &result.Title, &result.Artist, &result.Album, &result.AlbumArtist, &result.Genre, &result.TrackNumber,
			&result.TotalTracks, &result.DiscNumber, &result.TotalDiscs, &result.ReleaseDate, &result.MusicBrainzArtistID, &result.MusicBrainzAlbumID,
			&result.MusicBrainzTrackID, &result.Label, &result.MusicFolderId, &result.UserPlayCount, &result.GlobalPlayCount); err != nil {
			logger.Printf("Failed to scan row in SelectTracksByGenres: %v", err)
			return []types.MetadataWithPlaycounts{}, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return results, err
	}

	return results, nil
}
