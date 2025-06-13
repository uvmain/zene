package database

import (
	"context"
	"log"
	"path/filepath"
	"sync"
	"zene/core/config"
	"zene/core/io"
	"zene/core/logic"

	"zombiezen.com/go/sqlite/sqlitex"
)

var dbFile = "sqlite.db"
var DbPool *sqlitex.Pool
var dbMutex sync.RWMutex

func Initialise(ctx context.Context) {
	openDatabase(ctx)
	createMetadataTable(ctx)
	createAlbumArtTable(ctx)
	createArtistArtTable(ctx)
	createFts(ctx)
	CreateSessionsTable(ctx)
	StartSessionCleanupRoutine(ctx)
}

func openDatabase(ctx context.Context) {
	dbFile := filepath.Join(config.DatabaseDirectory, dbFile)

	if io.FileExists(dbFile) {
		log.Println("Database already exists")
	} else {
		log.Println("Creating new database file")
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

	if err := logic.CheckContext(ctx); err != nil {
		CloseDatabase()
	}
}

func CloseDatabase() {
	if DbPool != nil {
		DbPool.Close()
	}
}
