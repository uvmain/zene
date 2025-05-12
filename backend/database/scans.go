package database

import (
	"fmt"
	"zene/types"
)

func InsertScanRow(scanDate string, fileCount int, dateModified string) error {
	stmt, err := Db.Prepare(`INSERT INTO scans (scan_date, file_count, date_modified)
		VALUES ($scan_date, $file_count, $date_modified);`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
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

func SelectLastScan() (types.ScanRow, error) {
	stmt, err := Db.Prepare(`SELECT id, scan_date, file_count, date_modified from scans order by id desc limit 1;`)

	if err != nil {
		return types.ScanRow{}, err
	}
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
