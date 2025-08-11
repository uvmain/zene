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

	_, err := GetLatestVersion(ctx)
	if err == sql.ErrNoRows {
		newVersion := types.Version{
			ServerVersion:          "0.1.0",
			DatabaseVersion:        "1.0",
			SubsonicApiVersion:     "1.16.1",
			OpenSubsonicApiVersion: "1",
			Timestamp:              logic.GetCurrentTimeFormatted(),
		}
		logger.Printf("No versions found, inserting initial version: %v", newVersion)
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

func GetVersions(ctx context.Context) ([]types.Version, error) {
	query := "SELECT server_version, database_version, subsonic_api_version, open_subsonic_api_version, timestamp FROM versions order by id desc;"
	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return []types.Version{}, fmt.Errorf("querying versions in GetVersions: %v", err)
	}
	defer rows.Close()

	var result []types.Version
	for rows.Next() {
		var row types.Version
		err := rows.Scan(&row.ServerVersion, &row.DatabaseVersion, &row.SubsonicApiVersion, &row.OpenSubsonicApiVersion, &row.Timestamp)
		if err != nil {
			return []types.Version{}, fmt.Errorf("scanning version row: %v", err)
		}
		result = append(result, row)
	}
	return result, nil
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
