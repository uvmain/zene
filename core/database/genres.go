package database

import (
	"database/sql"
	"context"
	"fmt"
	"log"
	"strconv"
	"zene/core/logic"
	"zene/core/types"

)

func createGenresTable(ctx context.Context) error {
	tableName := "track_genres"
	schema := `CREATE TABLE IF NOT EXISTS track_genres (
		file_path TEXT NOT NULL,
		genre TEXT NOT NULL,
		FOREIGN KEY(file_path) REFERENCES metadata(file_path) ON DELETE CASCADE
	);`

	err := createTable(ctx, tableName, schema)
	if err != nil {
		return err
	}

	createIndex(ctx, "idx_track_genres_genre", "track_genres", "genre", false)
	createGenresTriggers(ctx)

	// get count of rows in track_genres table
	var query = "SELECT COUNT(*) as count FROM track_genres;")
	
	if // TODO: Query single row; err != nil {
		return fmt.Errorf("error checking count of track_genres table: %v", err)
	} else if hasRow {
		count := stmt.GetInt64("count")
		if count == 0 {
			log.Println("track_genres table is empty, populating from metadata")
			err = populateGenresFromMetadata(ctx)
			if err != nil {
				return fmt.Errorf("error populating track_genres table from metadata: %v", err)
			} else {
				log.Println("track_genres table populated from metadata")
			}
		}
	}

	return nil
}

func createGenresTriggers(ctx context.Context) {
	createTrigger(ctx, "tr_metadata_insert_genres", `CREATE TRIGGER tr_metadata_insert_genres AFTER INSERT ON metadata
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

	createTrigger(ctx, "tr_metadata_delete_genres", `CREATE TRIGGER tr_metadata_delete_genres AFTER DELETE ON metadata
    BEGIN
        DELETE FROM track_genres WHERE file_path = old.file_path;
    END;`)

	createTrigger(ctx, "tr_metadata_update_genres", `CREATE TRIGGER tr_metadata_update_genres AFTER UPDATE ON metadata
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
	var stmtText = `INSERT INTO track_genres (file_path, genre)
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


	err = sqlitex.ExecuteTransient(conn, stmtText, nil)
	if err != nil {
		return fmt.Errorf("error inserting data into track_genres table: %v", err)
	} else {
		return nil
	}
}

func SelectDistinctGenres(ctx context.Context, limitParam string, searchParam string) ([]types.GenreResponse, error) {


	stmtText := "select genre, count(file_path) as count from track_genres group by genre order by count desc"

	if limitParam != "" {
		limitInt, err := strconv.Atoi(limitParam)
		if err != nil {
			return []types.GenreResponse{}, fmt.Errorf("invalid limit value: %v", err)
		}
		stmtText = fmt.Sprintf("%s limit %d", stmtText, limitInt)
	}

	stmtText = fmt.Sprintf("%s;", stmtText)

	var query = stmtText)
	

	var rows []types.GenreResponse
	for {
		// TODO: Query single row
		if err != nil {
			return nil, err
		}
		if !hasRow {
			break
		}
		row := types.GenreResponse{
			Genre: stmt.GetText("genre"),
			Count: int(stmt.GetInt64("count")),
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func SelectTracksByGenres(ctx context.Context, genres []string, andOr string, limit int64, random string) ([]types.MetadataWithPlaycounts, error) {


	userId, _ := logic.GetUserIdFromContext(ctx)
	stmtText := getMetadataWithGenresSql(userId, genres, andOr, limit, random)

	var query = stmtText)
	

	var rows []types.MetadataWithPlaycounts

	for {
		if // TODO: Query single row; err != nil {
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
