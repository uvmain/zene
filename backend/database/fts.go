package database

import (
	"log"

	"zombiezen.com/go/sqlite/sqlitex"
)

func createFts() {
	createFtsTable()
	createFtsTriggers()
	insertFtsData()
}

func createFtsTable() {
	tableName := "track_metadata_fts"
	schema := `CREATE VIRTUAL TABLE IF NOT EXISTS track_metadata_fts USING fts5(file_id, filename, title, artist, album, album_artist, genre, release_date, label, tokenize="trigram");`
	createTable(tableName, schema)
}

func insertFtsData() {
	query := `INSERT INTO track_metadata_fts (file_id, filename, title, artist, album, album_artist, genre, release_date, label)
		SELECT file_id, filename, title, artist, album, album_artist, genre, release_date, label FROM track_metadata;`

	err := sqlitex.ExecuteTransient(DbRW, query, nil)
	if err != nil {
		log.Printf("Error inserting data into fts table: %s", err)
	} else {
		log.Println("Data inserted into fts table")
	}
}

func createFtsTriggers() {
	createTriggerIfNotExists("after_insert_fts", `CREATE TRIGGER after_insert_fts AFTER INSERT ON track_metadata
        BEGIN
            INSERT INTO track_metadata_fts (file_id, filename, title, artist, album, album_artist, genre, release_date, label)
            VALUES (new.file_id, new.filename, new.title, new.artist, new.album, new.album_artist, new.genre, new.release_date, new.label);
        END;`)

	createTriggerIfNotExists("after_delete_fts", `CREATE TRIGGER after_delete_fts AFTER DELETE ON track_metadata
    BEGIN
        DELETE FROM track_metadata_fts WHERE file_id = old.file_id;
    END;`)

	createTriggerIfNotExists("after_update_fts", `CREATE TRIGGER after_update_fts AFTER UPDATE ON track_metadata
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
