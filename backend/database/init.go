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
	createFilesTriggers()
	createAlbumArtTable()
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
