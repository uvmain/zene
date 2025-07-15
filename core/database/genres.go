package database

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"zene/core/logic"
	"zene/core/types"
)

func createGenresTable(ctx context.Context) error {
	tableName := "genres"
	schema := `CREATE TABLE IF NOT EXISTS genres (
		file_path TEXT NOT NULL,
		genre TEXT NOT NULL,
		FOREIGN KEY(file_path) REFERENCES metadata(file_path) ON DELETE CASCADE
	);`
	
	err := createTable(ctx, tableName, schema)
	if err != nil {
		return err
	}
	
	// Create index on genre column for performance
	createIndex(ctx, "idx_genres_genre", "genres", "genre", false)
	
	// Check if table was just created (empty) and populate it
	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("taking a db conn from the pool in createGenresTable: %v", err)
	}
	defer DbPool.Put(conn)
	
	// Check if genres table is empty
	stmt := conn.Prep("SELECT COUNT(*) as count FROM genres;")
	defer stmt.Finalize()
	
	if hasRow, err := stmt.Step(); err != nil {
		return fmt.Errorf("checking genres table count: %v", err)
	} else if hasRow {
		count := stmt.GetInt64("count")
		if count == 0 {
			// Table is empty, populate from existing metadata
			err = populateGenresFromMetadata(ctx)
			if err != nil {
				return fmt.Errorf("populating genres table: %v", err)
			}
		}
	}
	
	// Create triggers for metadata table
	createGenresTriggers(ctx)
	
	return nil
}

func populateGenresFromMetadata(ctx context.Context) error {
	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("taking a db conn from the pool in populateGenresFromMetadata: %v", err)
	}
	defer DbPool.Put(conn)
	
	// Get all metadata with genres
	stmt := conn.Prep("SELECT file_path, genre FROM metadata WHERE genre IS NOT NULL AND genre != '';")
	defer stmt.Finalize()
	
	var metadataRows []struct {
		FilePath string
		Genre    string
	}
	
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return fmt.Errorf("reading metadata for genres population: %v", err)
		} else if !hasRow {
			break
		} else {
			metadataRows = append(metadataRows, struct {
				FilePath string
				Genre    string
			}{
				FilePath: stmt.GetText("file_path"),
				Genre:    stmt.GetText("genre"),
			})
		}
	}
	
	// Insert genres for each metadata row
	for _, row := range metadataRows {
		err = insertGenresForFilePath(ctx, row.FilePath, row.Genre)
		if err != nil {
			return fmt.Errorf("inserting genres for %s: %v", row.FilePath, err)
		}
	}
	
	return nil
}

func insertGenresForFilePath(ctx context.Context, filePath, genreString string) error {
	if genreString == "" {
		return nil
	}
	
	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("taking a db conn from the pool in insertGenresForFilePath: %v", err)
	}
	defer DbPool.Put(conn)
	
	// Split genres by semicolon
	genres := strings.Split(genreString, ";")
	
	// Insert each genre
	insertStmt := conn.Prep("INSERT INTO genres (file_path, genre) VALUES ($file_path, $genre);")
	defer insertStmt.Finalize()
	
	for _, genre := range genres {
		trimmedGenre := strings.TrimSpace(genre)
		if trimmedGenre != "" {
			insertStmt.SetText("$file_path", filePath)
			insertStmt.SetText("$genre", trimmedGenre)
			
			_, err = insertStmt.Step()
			if err != nil {
				return fmt.Errorf("inserting genre %s for %s: %v", trimmedGenre, filePath, err)
			}
			insertStmt.Reset()
		}
	}
	
	return nil
}

func deleteGenresForFilePath(ctx context.Context, filePath string) error {
	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("taking a db conn from the pool in deleteGenresForFilePath: %v", err)
	}
	defer DbPool.Put(conn)
	
	stmt := conn.Prep("DELETE FROM genres WHERE file_path = $file_path;")
	defer stmt.Finalize()
	stmt.SetText("$file_path", filePath)
	
	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("deleting genres for %s: %v", filePath, err)
	}
	
	return nil
}

func createGenresTriggers(ctx context.Context) {
	// Trigger for INSERT on metadata table
	insertTriggerSQL := `
		CREATE TRIGGER IF NOT EXISTS metadata_insert_genres
		AFTER INSERT ON metadata
		FOR EACH ROW
		BEGIN
			DELETE FROM genres WHERE file_path = NEW.file_path;
			INSERT INTO genres (file_path, genre)
			SELECT NEW.file_path, TRIM(value) 
			FROM (
				WITH RECURSIVE split(file_path, genre, str) AS (
					SELECT NEW.file_path, '', NEW.genre || ';'
					UNION ALL
					SELECT file_path,
						TRIM(SUBSTR(str, 0, INSTR(str, ';'))),
						SUBSTR(str, INSTR(str, ';') + 1)
					FROM split 
					WHERE str != ''
				)
				SELECT genre as value FROM split WHERE genre != ''
			);
		END;`
	
	// Trigger for UPDATE on metadata table
	updateTriggerSQL := `
		CREATE TRIGGER IF NOT EXISTS metadata_update_genres
		AFTER UPDATE ON metadata
		FOR EACH ROW
		WHEN OLD.genre != NEW.genre
		BEGIN
			DELETE FROM genres WHERE file_path = NEW.file_path;
			INSERT INTO genres (file_path, genre)
			SELECT NEW.file_path, TRIM(value)
			FROM (
				WITH RECURSIVE split(file_path, genre, str) AS (
					SELECT NEW.file_path, '', NEW.genre || ';'
					UNION ALL
					SELECT file_path,
						TRIM(SUBSTR(str, 0, INSTR(str, ';'))),
						SUBSTR(str, INSTR(str, ';') + 1)
					FROM split 
					WHERE str != ''
				)
				SELECT genre as value FROM split WHERE genre != ''
			);
		END;`
	
	// Trigger for DELETE on metadata table - handled by foreign key cascade
	
	createTrigger(ctx, "metadata_insert_genres", insertTriggerSQL)
	createTrigger(ctx, "metadata_update_genres", updateTriggerSQL)
}

func SelectDistinctGenres(ctx context.Context, limitParam string, searchParam string) ([]types.GenreResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.GenreResponse{}, fmt.Errorf("taking a db conn from the pool in SelectDistinctGenres: %v", err)
	}
	defer DbPool.Put(conn)

	stmtText := "SELECT genre, COUNT(file_path) as count FROM genres"
	
	// Add search filter if provided
	if searchParam != "" {
		stmtText = fmt.Sprintf("%s WHERE genre LIKE '%%%s%%'", stmtText, strings.ReplaceAll(searchParam, "'", "''"))
	}
	
	// Group by genre and order by count
	stmtText = fmt.Sprintf("%s GROUP BY genre ORDER BY count DESC", stmtText)

	if limitParam != "" {
		limitInt, err := strconv.Atoi(limitParam)
		if err != nil {
			return []types.GenreResponse{}, fmt.Errorf("invalid limit value: %v", err)
		}
		stmtText = fmt.Sprintf("%s LIMIT %d", stmtText, limitInt)
	}

	stmtText = fmt.Sprintf("%s;", stmtText)

	stmt := conn.Prep(stmtText)
	defer stmt.Finalize()

	var genres []types.GenreResponse

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return []types.GenreResponse{}, err
		} else if !hasRow {
			break
		} else {
			genre := types.GenreResponse{
				Genre: stmt.GetText("genre"),
				Count: int(stmt.GetInt64("count")),
			}
			genres = append(genres, genre)
		}
	}

	if genres == nil {
		genres = []types.GenreResponse{}
	}
	
	return genres, nil
}

func SelectTracksByGenres(ctx context.Context, genres []string, andOr string, limit int64, random string) ([]types.MetadataWithPlaycounts, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.MetadataWithPlaycounts{}, fmt.Errorf("taking a db conn from the pool in SelectTracksByAlbumId: %v", err)
	}
	defer DbPool.Put(conn)

	userId, _ := logic.GetUserIdFromContext(ctx)
	stmtText := getMetadataWithGenresSql(userId, genres, andOr, limit, random)

	stmt := conn.Prep(stmtText)
	defer stmt.Finalize()

	var rows []types.MetadataWithPlaycounts

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return []types.MetadataWithPlaycounts{}, err
		} else if !hasRow {
			break
		} else {
			row := types.MetadataWithPlaycounts{
				FilePath:            stmt.GetText("file_path"),
				DateAdded:           stmt.GetText("date_added"),
				DateModified:        stmt.GetText("date_modified"),
				FileName:            stmt.GetText("file_name"),
				Format:              stmt.GetText("format"),
				Duration:            stmt.GetText("duration"),
				Size:                stmt.GetText("size"),
				Bitrate:             stmt.GetText("bitrate"),
				Title:               stmt.GetText("title"),
				Artist:              stmt.GetText("artist"),
				Album:               stmt.GetText("album"),
				AlbumArtist:         stmt.GetText("album_artist"),
				Genre:               stmt.GetText("genre"),
				TrackNumber:         stmt.GetText("track_number"),
				TotalTracks:         stmt.GetText("total_tracks"),
				DiscNumber:          stmt.GetText("disc_number"),
				TotalDiscs:          stmt.GetText("total_discs"),
				ReleaseDate:         stmt.GetText("release_date"),
				MusicBrainzArtistID: stmt.GetText("musicbrainz_artist_id"),
				MusicBrainzAlbumID:  stmt.GetText("musicbrainz_album_id"),
				MusicBrainzTrackID:  stmt.GetText("musicbrainz_track_id"),
				Label:               stmt.GetText("label"),
				UserPlayCount:       stmt.GetInt64("user_play_count"),
				GlobalPlayCount:     stmt.GetInt64("global_play_count"),
			}
			rows = append(rows, row)
		}
	}
	if rows == nil {
		rows = []types.MetadataWithPlaycounts{}
	}
	return rows, nil
}
