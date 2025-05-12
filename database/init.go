package database

import (
	"log"
	"path/filepath"
	"zene/config"
	"zene/io"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var dbFile = "sqlite.db"
var Db *sqlite.Conn

func Initialise() {
	io.CreateDir(config.DatabaseDirectory)
	openDatabase()
	createScansTable()
	createFilesTable()
	createMetadataTable()
}

func openDatabase() {
	dbFile := filepath.Join(config.DatabaseDirectory, "sqlite.db")

	if io.FileExists(dbFile) {
		log.Println("Database already exists")
	} else {
		log.Println("Creating database file")
	}

	conn, err := sqlite.OpenConn(dbFile, sqlite.OpenReadWrite|sqlite.OpenCreate)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	Db = conn

	err = sqlitex.ExecuteTransient(conn, "PRAGMA journal_mode=WAL;", nil)
	if err != nil {
		log.Fatalf("Failed to set WAL mode: %v", err)
	} else {
		log.Printf("Database is in WAL mode")
	}
}

func CloseDatabase() {
	if Db != nil {
		Db.Close()
	}
}

func createScansTable() {
	tableName := "scans"
	schema := `CREATE TABLE IF NOT EXISTS scans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		scan_date TEXT NOT NULL,
		file_count INTEGER NOT NULL,
		date_modified TEXT NOT NULL
	);`

	createTable(tableName, schema)
}

func createFilesTable() {
	tableName := "files"
	schema := `CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		dir_path TEXT NOT NULL,
		filename TEXT NOT NULL,
		date_added TEXT NOT NULL,
		date_modified TEXT NOT NULL
	);`
	createTable(tableName, schema)
}

func createMetadataTable() {
	tableName := "metadata"
	schema := `CREATE TABLE IF NOT EXISTS track_metadata (
		musicbrainz_track_id TEXT PRIMARY KEY,
		filename TEXT,
		format TEXT,
		duration TEXT,
		size TEXT,
		bitrate TEXT,
		title TEXT,
		artist TEXT,
		album TEXT,
		album_artist TEXT,
		genre TEXT,
		track_number TEXT,
		total_tracks TEXT,
		disc_number TEXT,
		total_discs TEXT,
		release_date TEXT,
		musicbrainz_artist_id TEXT,
		musicbrainz_album_id TEXT,
		label TEXT
	);`
	createTable(tableName, schema)
}
