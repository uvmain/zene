package database

import (
	"context"
	"zene/core/logger"
)

func createFtsTables(ctx context.Context) {
	createFtsMetadataTable(ctx)
	createFtsMetadataTriggers(ctx)
	insertFtsMetadataData(ctx)
	createFtsArtistsTable(ctx)
	createFtsArtistsTriggers(ctx)
	insertFtsArtistsData(ctx)
}

func createFtsMetadataTable(ctx context.Context) error {
	tableName := "metadata_fts"
	schema := `CREATE VIRTUAL TABLE IF NOT EXISTS metadata_fts USING fts5(file_path, file_name, title, artist, album, album_artist, genre, release_date, label, tokenize="trigram remove_diacritics 1");`
	err := createTable(ctx, tableName, schema)
	return err
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

	_, err := DB.ExecContext(ctx, query)
	if err != nil {
		logger.Printf("Database: error inserting data into metadata_fts table: %v", err)
	} else {
		logger.Println("Database: data inserted into metadata_fts table")
	}
}

func createFtsArtistsTable(ctx context.Context) error {
	tableName := "artists_fts"
	schema := `CREATE VIRTUAL TABLE IF NOT EXISTS artists_fts USING fts5(file_path, artist, tokenize="trigram remove_diacritics 1");`
	err := createTable(ctx, tableName, schema)
	return err
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
        UPDATE artists_fts SET artist = new.artist WHERE file_path = old.file_path;
    END;`)
}

func insertFtsArtistsData(ctx context.Context) {
	query := `INSERT INTO artists_fts (file_path, artist)
		SELECT file_path, artist FROM metadata;`

	_, err := DB.ExecContext(ctx, query)
	if err != nil {
		logger.Printf("Database: error inserting data into artists_fts table: %v", err)
	} else {
		logger.Println("Database: data inserted into artists_fts table")
	}
}
