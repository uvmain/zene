package database

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"
	"zene/core/config"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/logic"

	_ "modernc.org/sqlite"
)

var dbFile = "sqlite.db"
var DB *sql.DB
var err error

func Initialise(ctx context.Context) {
	openDatabase(ctx)
	createUsersTable(ctx)
	createSessionsTable(ctx)
	createMetadataTable(ctx)
	createPlayCountsTable(ctx)
	createLyricsTable(ctx)
	createAlbumArtTable(ctx)
	createArtistArtTable(ctx)
	createFtsTables(ctx)
	createGenresTable(ctx)
	createAudioCacheTable(ctx)
	createTemporaryTokensTable(ctx)
}

func openDatabase(ctx context.Context) {
	dbFile := filepath.Join(config.DatabaseDirectory, dbFile)

	if io.FileExists(dbFile) {
		logger.Println("Database already exists")
	} else {
		logger.Println("Creating new database file")
	}

	DB, err = sql.Open("sqlite", dbFile+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	if err := DB.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	logger.Println("Database opened with WAL mode enabled")

	if err := logic.CheckContext(ctx); err != nil {
		CloseDatabase()
	}
}

func CloseDatabase() {
	if DB != nil {
		DB.Close()
	}
}
