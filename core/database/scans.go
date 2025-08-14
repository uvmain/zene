package database

import (
	"context"
	"fmt"
)

type ScanRow struct {
	Count             int    `xml:"count,attr" json:"count"`
	FolderCount       int    `xml:"folderCount,attr" json:"folderCount"`
	ScanDate          string `xml:"scanDate,attr" json:"scanDate"`
	ScanType          string `xml:"scanType,attr" json:"scanType"`
	ScanCompletedDate string `xml:"scanCompletedDate,attr" json:"scanCompletedDate"`
}

func createScansTable(ctx context.Context) {
	schema := `CREATE TABLE scans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		count INTEGER NOT NULL,
		folder_count INTEGER NOT NULL,
		scan_date TEXT NOT NULL,
		scan_type TEXT NOT NULL,
		scan_completed_date TEXT,
	);`
	createTable(ctx, schema)
}

func GetScans(ctx context.Context, userId int64, metadataId string) ([]ScanRow, error) {
	query := `select count, folder_count, scan_date, scan_type, scan_completed_date from scans`

	rows, err := DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("getting scans: %v", err)
	}
	defer rows.Close()

	var scans []ScanRow
	for rows.Next() {
		var scan ScanRow
		if err := rows.Scan(&scan.Count, &scan.FolderCount, &scan.ScanDate, &scan.ScanType, &scan.ScanCompletedDate); err != nil {
			return nil, fmt.Errorf("scanning scan row: %v", err)
		}
		scans = append(scans, scan)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating scan rows: %v", err)
	}
	return scans, nil
}

func GetLatestScan(ctx context.Context) (ScanRow, error) {
	query := `SELECT count, folder_count, scan_date, scan_type, scan_completed_date FROM scans ORDER BY scan_date DESC LIMIT 1`
	row := DB.QueryRowContext(ctx, query)

	var scan ScanRow
	if err := row.Scan(&scan.Count, &scan.FolderCount, &scan.ScanDate, &scan.ScanType, &scan.ScanCompletedDate); err != nil {
		return ScanRow{}, fmt.Errorf("querying latest scan: %v", err)
	}
	return scan, nil
}

func InsertScan(ctx context.Context, scan ScanRow) error {
	query := `INSERT INTO scans (count, folder_count, scan_date, scan_type, scan_completed_date) VALUES (?, ?, ?, ?, ?)`
	if _, err := DB.ExecContext(ctx, query, scan.Count, scan.FolderCount, scan.ScanDate, scan.ScanType, scan.ScanCompletedDate); err != nil {
		return fmt.Errorf("inserting scan: %v", err)
	}
	return nil
}
