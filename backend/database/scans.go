package database

import (
	"fmt"
	"zene/types"
)

func createScansTable() {
	tableName := "scans"
	schema := `CREATE TABLE IF NOT EXISTS scans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		scan_date TEXT NOT NULL,
		file_count INTEGER NOT NULL,
		date_modified TEXT NOT NULL
	);`
	createTable(tableName, schema)
}

func InsertScanRow(scanDate string, fileCount int, dateModified string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtInsertScanRow
	stmt.Reset()
	stmt.ClearBindings()
	stmt.SetText("$scan_date", scanDate)
	stmt.SetInt64("$file_count", int64(fileCount))
	stmt.SetText("$date_modified", dateModified)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to insert Scans row: %v", err)
	}
	return nil
}

func SelectLastScan() (types.ScanRow, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtSelectLastScan
	stmt.Reset()

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
