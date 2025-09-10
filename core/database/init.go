package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"zene/core/config"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/version"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var dbFile = "sqlite.db"
var DB *sql.DB
var err error

func Initialise(ctx context.Context) {
	openDatabase(ctx)
	migrateVersions(ctx)
	migrateMusicFolders(ctx)
	migrateUsers(ctx)
	migrateApiKeys(ctx)
	CreateAdminUserIfRequired(ctx)
	migrateMetadata(ctx)
	migratePlayCounts(ctx)
	migrateChats(ctx)
	migrateLyrics(ctx)
	migrateArt(ctx)
	migrateTrackGenres(ctx)
	migrateGenreCounts(ctx)
	migrateAudioCache(ctx)
	migrateUserStars(ctx)
	migrateUserRatings(ctx)
	migrateScans(ctx)
	migrateNowPlaying(ctx)
	migrateSimilarArtists(ctx)
	migrateTopSongs(ctx)
	migratePlaylists(ctx)
	migrateInternetRadio(ctx)
	migrateBookmarks(ctx)

	checkVersion(ctx)
}

func checkVersion(ctx context.Context) {
	thisVersion := version.Version

	existingVersion, err := GetLatestVersion(ctx)
	if err != nil {
		logger.Printf("Error getting latest version, defaulting to DB version %s: %v", thisVersion.DatabaseVersion, err)
		existingVersion = version.Version
		existingVersion.DatabaseVersion = thisVersion.DatabaseVersion
		InsertVersion(ctx, existingVersion)
		existingVersion, err = GetLatestVersion(ctx)
	}

	if existingVersion.DatabaseVersion != thisVersion.DatabaseVersion {
		logger.Printf("Database version change detected, migrating from %v to %v", existingVersion.DatabaseVersion, thisVersion.DatabaseVersion)
		InsertVersion(ctx, thisVersion)
	} else if existingVersion.ServerVersion != thisVersion.ServerVersion {
		logger.Printf("Server version change detected, migrating from %v to %v", existingVersion.ServerVersion, thisVersion.ServerVersion)
		InsertVersion(ctx, thisVersion)
	} else if existingVersion.SubsonicApiVersion != thisVersion.SubsonicApiVersion {
		logger.Printf("Subsonic API version change detected, migrating from %v to %v", existingVersion.SubsonicApiVersion, thisVersion.SubsonicApiVersion)
		InsertVersion(ctx, thisVersion)
	} else if existingVersion.OpenSubsonicApiVersion != thisVersion.OpenSubsonicApiVersion {
		logger.Printf("OpenSubsonic API version change detected, migrating from %v to %v", existingVersion.OpenSubsonicApiVersion, thisVersion.OpenSubsonicApiVersion)
		InsertVersion(ctx, thisVersion)
	}
}

func openDatabase(ctx context.Context) {
	dbFilePath := filepath.Join(config.DatabaseDirectory, dbFile)

	if io.FileExists(dbFilePath) {
		logger.Println("Database already exists")
	} else {
		logger.Println("Creating new database file")
	}

	dataSource := fmt.Sprintf("file:%s?_journal_mode=WAL&_foreign_keys=on", dbFilePath)
	var err error
	DB, err = sql.Open("sqlite3", dataSource)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	DB.SetMaxIdleConns(5)

	if err := DB.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
}

func CleanShutdown() {
	if DB != nil {
		log.Println("Closing database...")
		DB.Exec("PRAGMA wal_checkpoint(FULL);")
		DB.Close()

		// wait for data to be flushed to disk
		f, err := os.OpenFile(filepath.Join(config.DatabaseDirectory, dbFile), os.O_RDWR, 0660)
		if err == nil {
			f.Sync()
			f.Close()
		}
		log.Println("Database closed.")
	}
}
