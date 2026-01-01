package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"
	"zene/core/config"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/version"

	_ "modernc.org/sqlite"
)

var dbFile = "sqlite.db"
var DbRead *sql.DB
var DbWrite *sql.DB
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
	migratePlayqueues(ctx)
	migratePodcasts(ctx)

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
		if err != nil {
			log.Fatalf("Error getting latest version after inserting default: %v", err)
		}
	}

	if existingVersion.DatabaseVersion != thisVersion.DatabaseVersion || existingVersion.ServerVersion != thisVersion.ServerVersion || existingVersion.SubsonicApiVersion != thisVersion.SubsonicApiVersion || existingVersion.OpenSubsonicApiVersion != thisVersion.OpenSubsonicApiVersion {
		logger.Printf("Version change detected, inserting new version record: old=%+v new=%+v", existingVersion, thisVersion)
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

	DbRead, err = sql.Open("sqlite", dbFilePath)
	if err != nil {
		log.Printf("Error opening read database file: %s", err)
	} else {
		log.Println("ReadOnly database connection opened")
	}

	_, err = DbRead.Exec("PRAGMA journal_mode=WAL;")
	_, err = DbRead.Exec("PRAGMA synchronous=NORMAL;")
	_, err = DbRead.Exec("PRAGMA busy_timeout=5000;")
	_, err = DbRead.Exec("PRAGMA temp_store=MEMORY;")
	_, err = DbRead.Exec("PRAGMA mmap_size=30000000000;")

	DbRead.SetMaxIdleConns(10)
	DbRead.SetMaxOpenConns(100)
	DbRead.SetConnMaxLifetime(time.Hour)

	DbWrite, err = sql.Open("sqlite", dbFilePath)
	if err != nil {
		log.Printf("Error opening write database file: %s", err)
	} else {
		log.Println("Write database connection opened")
	}

	_, err = DbWrite.Exec("PRAGMA journal_mode=WAL;")
	_, err = DbWrite.Exec("PRAGMA synchronous=NORMAL;")
	_, err = DbWrite.Exec("PRAGMA busy_timeout=5000;")
	_, err = DbWrite.Exec("PRAGMA temp_store=MEMORY;")
	_, err = DbWrite.Exec("PRAGMA mmap_size=30000000000;")

	DbWrite.SetMaxIdleConns(1)
	DbWrite.SetMaxOpenConns(1)
	DbWrite.SetConnMaxLifetime(time.Hour)
	_, err = DbWrite.Exec("PRAGMA locking_mode=IMMEDIATE;")

	if err := DbRead.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database with read connection: %v", err)
	}

	if err := DbWrite.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database with write connection: %v", err)
	}
}

func CleanShutdown() {
	if DbWrite != nil {
		log.Println("Closing database...")
		_, err = DbWrite.Exec("PRAGMA wal_checkpoint(FULL);")
		if err != nil {
			log.Printf("Error committing WAL checkpoint on shutdown: %v", err)
		}
		DbWrite.Close()
		DbRead.Close()

		// wait for data to be flushed to disk
		f, err := os.OpenFile(filepath.Join(config.DatabaseDirectory, dbFile), os.O_RDWR, 0660)
		if err == nil {
			f.Sync()
			f.Close()
		}
		log.Println("Database closed.")
	}
}
