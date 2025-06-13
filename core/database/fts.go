package database

import (
	"context"
	"log"

	"zombiezen.com/go/sqlite/sqlitex"
)

func createFtsTables(ctx context.Context) {
	createFtsMetadataTable(ctx)
	createFtsMetadataTriggers(ctx)
	insertFtsMetadataData(ctx)
	createFtsArtistsTable(ctx)
	createFtsArtistsTriggers(ctx)
	insertFtsArtistsData(ctx)
}

func createFtsMetadataTable(ctx context.Context) {
	tableName := "metadata_fts"
	schema := `CREATE VIRTUAL TABLE IF NOT EXISTS metadata_fts USING fts5(file_path, file_name, title, artist, album, album_artist, genre, release_date, label, tokenize="trigram remove_diacritics 1");`
	createTable(ctx, tableName, schema)
}

func createFtsMetadataTriggers(ctx context.Context) {
	createTrigger(ctx, "tr_metadata_insert_fts", `CREATE TRIGGER tr_metadata_insert_fts AFTER INSERT ON metadata
        BEGIN
            INSERT INTO metadata_fts (file_path, file_name, title, artist, album, album_artist, genre, release_date, label)
            VALUES (new.file_path, new.file_name, new.title, new.artist, new.album, new.album_artist, new.genre, new.release_date, new.label);
        END;`)

	createTrigger(ctx, "tr_metadata_delete_fts", `CREATE TRIGGER tr_metadata_delete_fts AFTER DELETE ON metadata
    BEGIN
        DELETE FROM metadata_fts WHERE file_path = old.file_path;
    END;`)

	createTrigger(ctx, "tr_metadata_update_fts", `CREATE TRIGGER tr_metadata_update_fts AFTER UPDATE ON metadata
    BEGIN
        UPDATE metadata_fts SET
						file_name = new.file_name,
						title = new.title,
						artist = new.artist,
						album = new.album,
						album_artist = new.album_artist,
						genre = new.genre,
						release_date = new.release_date,
						label = new.label
        WHERE file_path = new.file_path;
    END;`)
}

func insertFtsMetadataData(ctx context.Context) {
	const query = `
		INSERT INTO metadata_fts (
			file_path, file_name, title, artist, album, album_artist, genre, release_date, label
		)
		SELECT 
			file_path, file_name, title, artist, album, album_artist, genre, release_date, label 
		FROM metadata;`

	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in insertFtsMetadataData: %v", err)
		return
	}
	defer DbPool.Put(conn)

	err = sqlitex.ExecuteTransient(conn, query, nil)
	if err != nil {
		log.Printf("Error inserting data into metadata_fts table: %s", err)
	} else {
		log.Println("Data inserted into metadata_fts table")
	}
}

func createFtsArtistsTable(ctx context.Context) {
	tableName := "artists_fts"
	schema := `CREATE VIRTUAL TABLE IF NOT EXISTS artists_fts USING fts5(file_path, artist, tokenize="trigram remove_diacritics 1");`
	createTable(ctx, tableName, schema)
}

func createFtsArtistsTriggers(ctx context.Context) {
	createTrigger(ctx, "tr_metadata_insert_artists_fts", `CREATE TRIGGER tr_metadata_insert_artists_fts AFTER INSERT ON metadata
        BEGIN
            INSERT INTO artists_fts (file_path, artist) VALUES (new.file_path, new.artist);
        END;`)

	createTrigger(ctx, "tr_metadata_delete_artists_fts", `CREATE TRIGGER tr_metadata_delete_artists_fts AFTER DELETE ON metadata
    BEGIN
        DELETE FROM artists_fts WHERE file_path = old.file_path;
    END;`)

	createTrigger(ctx, "tr_metadata_update_artists_fts", `CREATE TRIGGER tr_metadata_update_artists_fts AFTER UPDATE ON metadata
    BEGIN
        UPDATE artists_fts SET metadata_id = new.file_path, artist = new.artist WHERE metadata_id = old.file_path;
    END;`)
}

func insertFtsArtistsData(ctx context.Context) {
	query := `INSERT INTO artists_fts (metadata_id, artist)
		SELECT file_path, artist FROM metadata;`

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
