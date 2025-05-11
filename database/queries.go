package database

import (
	"fmt"
	"zene/types"
)

func SelectAllFiles() ([]types.FilesRow, error) {
	stmt, _ := Db.Prepare(`SELECT dir_path, filename, date_added, mdate FROM files;`)

	var rows []types.FilesRow

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
			row.Mdate = stmt.GetText("mdate")
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func SelectFileByFilename(filename string) (types.FilesRow, error) {
	stmt, _ := Db.Prepare(`SELECT id, dir_path, filename, date_added, mdate FROM files WHERE filename = $filename;`)
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
		row.Mdate = stmt.GetText("mdate")
		return row, nil
	}
}

func InsertIntoFiles(dirPath string, fileName string, dateAdded string, dateModified string) error {
	stmt, err := Db.Prepare(`INSERT INTO files (dir_path, filename, date_added, mdate)
		VALUES ($dir_path, $filename, $date_added, $mdate);`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Finalize()

	stmt.SetText("$dir_path", dirPath)
	stmt.SetText("$filename", fileName)
	stmt.SetText("$date_added", dateAdded)
	stmt.SetText("$mdate", dateModified)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to insert file: %v", err)
	}
	return nil
}
