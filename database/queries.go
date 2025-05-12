package database

import (
	"fmt"
	"zene/types"
)

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
		row.Id = stmt.GetText("id")
		row.ScanDate = stmt.GetText("scan_date")
		row.FileCount = stmt.GetText("file_count")
		row.DateModified = stmt.GetText("date_modified")
		return row, nil
	}
}

func SelectAllFiles() ([]types.FilesRow, error) {
	stmt, err := Db.Prepare(`SELECT id, dir_path, filename, date_added, date_modified FROM files;`)

	var rows []types.FilesRow

	if err != nil {
		return rows, err
	}
	defer stmt.Finalize()

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return []types.FilesRow{}, err
		} else if !hasRow {
			break
		} else {
			var row types.FilesRow
			row.Id = stmt.GetText("id")
			row.DirPath = stmt.GetText("dir_path")
			row.Filename = stmt.GetText("filename")
			row.DateAdded = stmt.GetText("date_added")
			row.DateModified = stmt.GetText("date_modified")
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func SelectFileByFilename(filename string) (types.FilesRow, error) {
	stmt, err := Db.Prepare(`SELECT id, dir_path, filename, date_added, date_modified FROM files WHERE filename = $filename;`)
	if err != nil {
		return types.FilesRow{}, err
	}
	defer stmt.Finalize()

	stmt.SetText("$filename", filename)

	if hasRow, err := stmt.Step(); err != nil {
		return types.FilesRow{}, err
	} else if !hasRow {
		return types.FilesRow{}, nil
	} else {
		var row types.FilesRow
		row.Id = stmt.GetText("id")
		row.DirPath = stmt.GetText("dir_path")
		row.Filename = stmt.GetText("filename")
		row.DateAdded = stmt.GetText("date_added")
		row.DateModified = stmt.GetText("date_modified")
		return row, nil
	}
}

func InsertIntoFiles(dirPath string, fileName string, dateAdded string, dateModified string) error {
	stmt, err := Db.Prepare(`INSERT INTO files (dir_path, filename, date_added, date_modified)
		VALUES ($dir_path, $filename, $date_added, $date_modified);`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Finalize()

	stmt.SetText("$dir_path", dirPath)
	stmt.SetText("$filename", fileName)
	stmt.SetText("$date_added", dateAdded)
	stmt.SetText("$date_modified", dateModified)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to insert file: %v", err)
	}
	return nil
}

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
