package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
	"zene/core/version"
)

func createVersionsTable(ctx context.Context) {
	schema := `CREATE TABLE versions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		server_version TEXT NOT NULL,
		database_version TEXT NOT NULL,
		subsonic_api_version TEXT NOT NULL,
		open_subsonic_api_version TEXT NOT NULL,
		timestamp TEXT NOT NULL
	);`
	createTable(ctx, schema)

	newVersion := types.Version{
		ServerVersion:          version.Version.ServerVersion,
		DatabaseVersion:        version.Version.DatabaseVersion,
		SubsonicApiVersion:     version.Version.SubsonicApiVersion,
		OpenSubsonicApiVersion: version.Version.OpenSubsonicApiVersion,
		Timestamp:              logic.GetCurrentTimeFormatted(),
	}

	latestVersion, err := GetLatestVersion(ctx)
	if err == sql.ErrNoRows {
		logger.Printf("No versions found, inserting initial version: %v", newVersion)
		err := InsertVersion(ctx, newVersion)
		if err != nil {
			log.Fatalf("Error inserting version: %v", err)
		}
	}

	if latestVersion.ServerVersion != newVersion.ServerVersion {
		logger.Printf("updating versions table: %v", newVersion)
		err := InsertVersion(ctx, newVersion)
		if err != nil {
			log.Fatalf("Error inserting version: %v", err)
		}
	}
}

func InsertVersion(ctx context.Context, version types.Version) error {
	insertTimestampUnixSeconds := time.Now().Unix()
	query := `INSERT INTO versions (server_version, database_version, subsonic_api_version, open_subsonic_api_version, timestamp)
		VALUES (?, ?, ?, ?, ?);`
	_, err := DB.ExecContext(ctx, query,
		version.ServerVersion,
		version.DatabaseVersion,
		version.SubsonicApiVersion,
		version.OpenSubsonicApiVersion,
		insertTimestampUnixSeconds)
	if err != nil {
		return fmt.Errorf("inserting version in InsertVersion: %v", err)
	}
	return nil
}

func GetLatestVersion(ctx context.Context) (types.Version, error) {
	query := "SELECT server_version, database_version, subsonic_api_version, open_subsonic_api_version, timestamp FROM versions ORDER BY id desc limit 1;"

	var result types.Version
	err := DB.QueryRowContext(ctx, query).Scan(&result.ServerVersion, &result.DatabaseVersion, &result.SubsonicApiVersion, &result.OpenSubsonicApiVersion, &result.Timestamp)
	if err == sql.ErrNoRows {
		return types.Version{}, err
	} else if err != nil {
		logger.Printf("Error querying version: %v", err)
		return types.Version{}, err
	}
	return result, nil
}
