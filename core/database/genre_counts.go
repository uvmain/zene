package database

import (
	"context"
	"fmt"
	"zene/core/logger"
)

func migrateGenreCounts(ctx context.Context) {
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
