package database

import (
	"context"
	"log"
	"path/filepath"
	"sync"
	"zene/core/config"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/logic"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var dbFile = "sqlite.db"
var DbPool *sqlitex.Pool
var dbMutex sync.RWMutex

func Initialise(ctx context.Context) {
	openDatabase(ctx)
	createUsersTable(ctx)
	createSessionsTable(ctx)
	createMetadataTable(ctx)
	createPlayCountsTable(ctx)
	createAlbumArtTable(ctx)
	createArtistArtTable(ctx)
	createFtsTables(ctx)
	createAudioCacheTable(ctx)
}

func openDatabase(ctx context.Context) {
	dbFile := filepath.Join(config.DatabaseDirectory, dbFile)

	if io.FileExists(dbFile) {
		logger.Println("Database already exists")
	} else {
		logger.Println("Creating new database file")
	}

	poolOptions := sqlitex.PoolOptions{
		PrepareConn: func(conn *sqlite.Conn) error {
			err := sqlitex.ExecuteTransient(conn, `PRAGMA foreign_keys = on;`, nil)
			if err != nil {
				logger.Printf("Prepare internal error: %v", err)
			}
			return err
		},
	}
	poolOptions.Flags = 0
	poolOptions.PoolSize = 10

	var err error

	DbPool, err = sqlitex.NewPool(dbFile, poolOptions)
	if err != nil {
		log.Fatalf("Failed to open database pool: %v", err)
	} else {
		logger.Println("Database pool opened")
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
