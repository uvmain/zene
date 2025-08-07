package database

import (
	"context"
	"database/sql"
	"fmt"
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
	createMusicFoldersTable(ctx)
	createUsersTable(ctx)
	createApiKeysTable(ctx)
	CreateAdminUserIfRequired(ctx)
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

	dataSource := fmt.Sprintf("file:%s?_journal_mode=WAL&_foreign_keys=on", dbFilePath)
	DB, err = sql.Open("sqlite3", dataSource)

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
