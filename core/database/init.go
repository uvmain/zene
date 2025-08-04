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

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var dbFile = "sqlite.db"
var DB *sql.DB
var err error

func Initialise(ctx context.Context) {
	openDatabase(ctx)
	createUsersTable(ctx)
	createApiKeysTable(ctx)
	createMetadataTable(ctx)
	createPlayCountsTable(ctx)
	createLyricsTable(ctx)
	createAlbumArtTable(ctx)
	createArtistArtTable(ctx)
	createFtsTables(ctx)
	createGenresTable(ctx)
	createAudioCacheTable(ctx)
}

func openDatabase(ctx context.Context) {
	dbFilePath := filepath.Join(config.DatabaseDirectory, dbFile)

	if io.FileExists(dbFilePath) {
		logger.Println("Database already exists")
	} else {
		logger.Println("Creating new database file")
	}

	DB, err = sql.Open("sqlite3", dbFilePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	log.Printf("Database file opened: %s", dbFilePath)

	_, err = DB.Exec("pragma foreign_keys = 1;")
	if err != nil {
		log.Printf("Error enabling foreign keys: %s", err)
	} else {
		log.Println("Foreign keys enabled")
	}

	_, err = DB.Exec("pragma journal_mode = wal;")
	if err != nil {
		log.Printf("Error entering WAL mode: %s", err)
	} else {
		log.Println("Database is in WAL mode")
	}

	// DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	if err := DB.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	if err := logic.CheckContext(ctx); err != nil {
		CloseDatabase()
	}
}

func CloseDatabase() {
	if DB != nil {
		DB.Close()
	}
}
