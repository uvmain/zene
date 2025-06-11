package database

import (
	"context"
	"log"

	"zombiezen.com/go/sqlite/sqlitex"
)

func createFts(ctx context.Context) {
	createFtsMetadataTable(ctx)
	createFtsMetadataTriggers(ctx)
	insertFtsMetadataData(ctx)
	createFtsArtistsTable(ctx)
	createFtsArtistsTriggers(ctx)
	insertFtsArtistsData(ctx)
}

func createFtsMetadataTable(ctx context.Context) {
	tableName := "track_metadata_fts"
	schema := `CREATE VIRTUAL TABLE IF NOT EXISTS track_metadata_fts USING fts5(file_id, filename, title, artist, album, album_artist, genre, release_date, label, tokenize="trigram remove_diacritics 1");`
	createTable(ctx, tableName, schema)
}

func createFtsMetadataTriggers(ctx context.Context) {
	createTriggerIfNotExists(ctx, "after_insert_fts", `CREATE TRIGGER after_insert_fts AFTER INSERT ON track_metadata
        BEGIN
            INSERT INTO track_metadata_fts (file_id, filename, title, artist, album, album_artist, genre, release_date, label)
            VALUES (new.file_id, new.filename, new.title, new.artist, new.album, new.album_artist, new.genre, new.release_date, new.label);
        END;`)

	createTriggerIfNotExists(ctx, "after_delete_fts", `CREATE TRIGGER after_delete_fts AFTER DELETE ON track_metadata
    BEGIN
        DELETE FROM track_metadata_fts WHERE file_id = old.file_id;
    END;`)

	createTriggerIfNotExists(ctx, "after_update_fts", `CREATE TRIGGER after_update_fts AFTER UPDATE ON track_metadata
    BEGIN
        UPDATE track_metadata_fts SET
						file_id = new.file_id,
            filename = new.filename,
						title = new.title,
						artist = new.artist,
						album = new.album,
						album_artist = new.album_artist,
						genre = new.genre,
						release_date = new.release_date,
						label = new.label
        WHERE file_id = old.file_id;
    END;`)
}

func insertFtsMetadataData(ctx context.Context) {
	const query = `
		INSERT INTO track_metadata_fts (
			file_id, filename, title, artist, album, album_artist, genre, release_date, label
		)
		SELECT 
			file_id, filename, title, artist, album, album_artist, genre, release_date, label 
		FROM track_metadata;`

	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in insertFtsMetadataData: %v", err)
		return
	}
	defer DbPool.Put(conn)

	err = sqlitex.ExecuteTransient(conn, query, nil)
	if err != nil {
		log.Printf("Error inserting data into track_metadata_fts table: %s", err)
	} else {
		log.Println("Data inserted into track_metadata_fts table")
	}
}

func createFtsArtistsTable(ctx context.Context) {
	tableName := "artists_fts"
	schema := `CREATE VIRTUAL TABLE IF NOT EXISTS artists_fts USING fts5(file_id, artist, tokenize="trigram remove_diacritics 1");`
	createTable(ctx, tableName, schema)
}

func createFtsArtistsTriggers(ctx context.Context) {
	createTriggerIfNotExists(ctx, "after_insert_artists_fts", `CREATE TRIGGER after_insert_artists_fts AFTER INSERT ON track_metadata
        BEGIN
            INSERT INTO artists_fts (file_id, artist) VALUES (new.file_id, new.artist);
        END;`)

	createTriggerIfNotExists(ctx, "after_delete_artists_fts", `CREATE TRIGGER after_delete_artists_fts AFTER DELETE ON track_metadata
    BEGIN
        DELETE FROM artists_fts WHERE file_id = old.file_id;
    END;`)

	createTriggerIfNotExists(ctx, "after_update_artists_fts", `CREATE TRIGGER after_update_artists_fts AFTER UPDATE ON track_metadata
    BEGIN
        UPDATE artists_fts SET file_id = new.file_id, artist = new.artist WHERE file_id = old.file_id;
    END;`)
}

func insertFtsArtistsData(ctx context.Context) {
	query := `INSERT INTO artists_fts (file_id, artist)
		SELECT file_id, artist FROM track_metadata;`

	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in insertFtsArtistsData: %v", err)
		return
	}
	defer DbPool.Put(conn)
	err = sqlitex.ExecuteTransient(conn, query, nil)
	if err != nil {
		log.Printf("Error inserting data into artists_fts table: %s", err)
	} else {
		log.Println("Data inserted into artists_fts table")
	}
}
