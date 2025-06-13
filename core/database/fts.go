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
	schema := `CREATE VIRTUAL TABLE IF NOT EXISTS track_metadata_fts USING fts5(metadata_id, file_name, title, artist, album, album_artist, genre, release_date, label, tokenize="trigram remove_diacritics 1");`
	createTable(ctx, tableName, schema)
}

func createFtsMetadataTriggers(ctx context.Context) {
	createTrigger(ctx, "after_insert_fts", `CREATE TRIGGER after_insert_fts AFTER INSERT ON track_metadata
        BEGIN
            INSERT INTO track_metadata_fts (metadata_id, file_name, title, artist, album, album_artist, genre, release_date, label)
            VALUES (new.id, new.file_name, new.title, new.artist, new.album, new.album_artist, new.genre, new.release_date, new.label);
        END;`)

	createTrigger(ctx, "after_delete_fts", `CREATE TRIGGER after_delete_fts AFTER DELETE ON track_metadata
    BEGIN
        DELETE FROM track_metadata_fts WHERE metadata_id = old.id;
    END;`)

	createTrigger(ctx, "after_update_fts", `CREATE TRIGGER after_update_fts AFTER UPDATE ON track_metadata
    BEGIN
        UPDATE track_metadata_fts SET
						metadata_id = new.id,
            file_name = new.file_name,
						title = new.title,
						artist = new.artist,
						album = new.album,
						album_artist = new.album_artist,
						genre = new.genre,
						release_date = new.release_date,
						label = new.label
        WHERE filmetadata_ide_id = old.metadata_id;
    END;`)
}

func insertFtsMetadataData(ctx context.Context) {
	const query = `
		INSERT INTO track_metadata_fts (
			metadata_id, file_name, title, artist, album, album_artist, genre, release_date, label
		)
		SELECT 
			id, file_name, title, artist, album, album_artist, genre, release_date, label 
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
	schema := `CREATE VIRTUAL TABLE IF NOT EXISTS artists_fts USING fts5(metadata_id, artist, tokenize="trigram remove_diacritics 1");`
	createTable(ctx, tableName, schema)
}

func createFtsArtistsTriggers(ctx context.Context) {
	createTrigger(ctx, "after_insert_artists_fts", `CREATE TRIGGER after_insert_artists_fts AFTER INSERT ON track_metadata
        BEGIN
            INSERT INTO artists_fts (metadata_id, artist) VALUES (new.id, new.artist);
        END;`)

	createTrigger(ctx, "after_delete_artists_fts", `CREATE TRIGGER after_delete_artists_fts AFTER DELETE ON track_metadata
    BEGIN
        DELETE FROM artists_fts WHERE metadata_id = old.id;
    END;`)

	createTrigger(ctx, "after_update_artists_fts", `CREATE TRIGGER after_update_artists_fts AFTER UPDATE ON track_metadata
    BEGIN
        UPDATE artists_fts SET metadata_id = new.id, artist = new.artist WHERE metadata_id = old.id;
    END;`)
}

func insertFtsArtistsData(ctx context.Context) {
	query := `INSERT INTO artists_fts (metadata_id, artist)
		SELECT id, artist FROM track_metadata;`

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
