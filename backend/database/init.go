package database

import (
	"log"
	"path/filepath"
	"sync"
	"zene/config"
	"zene/io"

	"zombiezen.com/go/sqlite/sqlitex"
)

var dbFile = "sqlite.db"
var DbPool *sqlitex.Pool
var dbMutex sync.RWMutex

func Initialise() {
	io.CreateDir(config.DatabaseDirectory)
	openDatabase()
	createScansTable()
	createFilesTable()
	createFilesTriggers()
	createMetadataTable()
	createMetadataTriggers()
	createAlbumArtTable()
	createArtistArtTable()
	createFts()
}

func openDatabase() {
	dbFile := filepath.Join(config.DatabaseDirectory, dbFile)

	if io.FileExists(dbFile) {
		log.Println("Database already exists")
	} else {
		log.Println("Creating database file")
	}

	poolOptions := sqlitex.PoolOptions{}
	poolOptions.Flags = 0
	poolOptions.PoolSize = 10

	var err error

	DbPool, err = sqlitex.NewPool(dbFile, poolOptions)
	if err != nil {
		log.Fatalf("Failed to open database pool: %v", err)
	} else {
		log.Println("Database pool opened")
	}
}

func CloseDatabase() {
	if DbPool != nil {
		DbPool.Close()
	}
}
