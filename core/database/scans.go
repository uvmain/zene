package database

import (
	"context"
	"database/sql"
	"fmt"
)

type ScanRow struct {
	Id            int64  `xml:"id,attr" json:"id"`
	Count         int64  `xml:"count,attr" json:"count"`
	FolderCount   int64  `xml:"folderCount,attr" json:"folderCount"`
	StartedDate   string `xml:"startedDate,attr" json:"startedDate"`
	Type          string `xml:"type,attr" json:"type"`
	CompletedDate string `xml:"completedDate,attr" json:"completedDate"`
}

func createScansTable(ctx context.Context) {
	schema := `CREATE TABLE scans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		count INTEGER NOT NULL,
		folder_count INTEGER NOT NULL,
		started_date TEXT NOT NULL,
		type TEXT NOT NULL,
		completed_date TEXT
	);`
	createTable(ctx, schema)
}

func GetScans(ctx context.Context, userId int64, metadataId string) ([]ScanRow, error) {
	query := `select count, folder_count, started_date, type, completed_date from scans`

	rows, err := DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("getting scans: %v", err)
	}
	defer rows.Close()

	var scans []ScanRow
	for rows.Next() {
		var scan ScanRow
		if err := rows.Scan(&scan.Count, &scan.FolderCount, &scan.StartedDate, &scan.Type, &scan.CompletedDate); err != nil {
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
	query := `SELECT id, count, folder_count, started_date, type, completed_date FROM scans ORDER BY id DESC LIMIT 1`
	row := DB.QueryRowContext(ctx, query)

	var scan ScanRow
	if err := row.Scan(&scan.Id, &scan.Count, &scan.FolderCount, &scan.StartedDate, &scan.Type, &scan.CompletedDate); err != nil {
		if err == sql.ErrNoRows {
			return ScanRow{}, err
		}
		return ScanRow{}, fmt.Errorf("querying latest scan: %v", err)
	}
	return scan, nil
}

func GetLatestCompletedScan(ctx context.Context) (ScanRow, error) {
	query := `SELECT id, count, folder_count, started_date, type, completed_date FROM scans WHERE completed_date IS NOT '' ORDER BY id DESC LIMIT 1`
	row := DB.QueryRowContext(ctx, query)

	var scan ScanRow
	if err := row.Scan(&scan.Id, &scan.Count, &scan.FolderCount, &scan.StartedDate, &scan.Type, &scan.CompletedDate); err != nil {
		if err == sql.ErrNoRows {
			return ScanRow{}, err
		}
		return ScanRow{}, fmt.Errorf("querying latest scan: %v", err)
	}
	return scan, nil
}

func InsertScan(ctx context.Context, scan ScanRow) (int64, error) {
	query := `INSERT INTO scans (count, folder_count, started_date, type, completed_date) VALUES (?, ?, ?, ?, ?)`
	result, err := DB.ExecContext(ctx, query, scan.Count, scan.FolderCount, scan.StartedDate, scan.Type, scan.CompletedDate)
	if err != nil {
		return 0, fmt.Errorf("inserting scan: %v", err)
	}
	return result.LastInsertId()
}

func UpdateScanProgress(ctx context.Context, scanId int64, scanRow ScanRow) error {
	query := `UPDATE scans SET count = ?, folder_count = ?, completed_date = ? WHERE id = ?`
	_, err := DB.ExecContext(ctx, query, scanRow.Count, scanRow.FolderCount, scanRow.CompletedDate, scanId)
	if err != nil {
		return fmt.Errorf("updating scan progress: %v", err)
	}
	return nil
}
