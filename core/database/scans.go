package database

import (
	"context"
	"fmt"
	"log"
	"zene/core/types"
)

func createScansTable(ctx context.Context) {
	tableName := "scans"
	schema := `CREATE TABLE IF NOT EXISTS scans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		scan_date TEXT NOT NULL,
		file_count INTEGER NOT NULL,
		date_modified TEXT NOT NULL
	);`
	createTable(ctx, tableName, schema)
}

func InsertScanRow(ctx context.Context, scanDate string, fileCount int, dateModified string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in InsertScanRow: %v", err)
		return err
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`INSERT INTO scans (scan_date, file_count, date_modified) VALUES ($scan_date, $file_count, $date_modified);`)
	defer stmt.Finalize()
	stmt.SetText("$scan_date", scanDate)
	stmt.SetInt64("$file_count", int64(fileCount))
	stmt.SetText("$date_modified", dateModified)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to insert Scans row: %v", err)
	}
	return nil
}

func SelectLastScan(ctx context.Context) (types.ScanRow, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in SelectLastScan: %v", err)
		return types.ScanRow{}, err
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT id, scan_date, file_count, date_modified from scans order by id desc limit 1;`)
	defer stmt.Finalize()

	hasRow, err := stmt.Step()

	if err != nil {
		return types.ScanRow{}, err
	} else if !hasRow {
		return types.ScanRow{}, nil
	} else {
		var row types.ScanRow
		row.Id = int(stmt.GetInt64("id"))
		row.ScanDate = stmt.GetText("scan_date")
		row.FileCount = stmt.GetText("file_count")
		row.DateModified = stmt.GetText("date_modified")
		return row, nil
	}
}
