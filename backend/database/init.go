package database

import (
	"log"
	"path/filepath"
	"sync"
	"zene/config"
	"zene/io"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var dbFile = "sqlite.db"
var DbRW *sqlite.Conn
var DbReadOnly *sqlite.Conn
var dbMutex sync.Mutex

func Initialise() {
	io.CreateDir(config.DatabaseDirectory)
	openDatabase()
	prepareStatements()
	createScansTable()
	createFilesTable()
	createFilesTriggers()
	createMetadataTable()
	createMetadataTriggers()
	createAlbumArtTable()
}

func openDatabase() {
	dbFile := filepath.Join(config.DatabaseDirectory, "sqlite.db")

	if io.FileExists(dbFile) {
		log.Println("Database already exists")
	} else {
		log.Println("Creating database file")
	}

	var err error

	DbRW, err = sqlite.OpenConn(dbFile, sqlite.OpenReadWrite|sqlite.OpenCreate)
	if err != nil {
		log.Fatalf("Failed to open CRUD database connection: %v", err)
	}

	DbReadOnly, err = sqlite.OpenConn(dbFile, sqlite.OpenReadOnly)
	if err != nil {
		log.Fatalf("Failed to open read-only database connection: %v", err)
	}

	err = sqlitex.ExecuteTransient(DbRW, "PRAGMA journal_mode=WAL;", nil)
	if err != nil {
		log.Fatalf("Failed to set WAL mode: %v", err)
	} else {
		log.Printf("Database is in WAL mode")
	}
}

func CloseDatabase() {
	if DbRW != nil {
		DbRW.Close()
	}
	if DbReadOnly != nil {
		DbReadOnly.Close()
	}
}
