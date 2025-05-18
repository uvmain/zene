package database

import (
	"fmt"
	"path/filepath"
	"zene/types"
)

func createFilesTable() {
	tableName := "files"
	schema := `CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		dir_path TEXT NOT NULL,
		filename TEXT NOT NULL,
		date_added TEXT NOT NULL,
		date_modified TEXT NOT NULL
	);`
	createTable(tableName, schema)
}

func createFilesTriggers() {
	createTriggerIfNotExists("files_after_delete_track_metadata", `CREATE TRIGGER files_after_delete_track_metadata AFTER DELETE ON files
	BEGIN
			DELETE FROM track_metadata WHERE file_id = old.id;
	END;`)
}

func SelectAllFiles() ([]types.FilesRow, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtSelectAllFiles
	stmt.Reset()

	var rows []types.FilesRow

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return []types.FilesRow{}, err
		} else if !hasRow {
			break
		} else {
			var row types.FilesRow
			row.Id = int(stmt.GetInt64("id"))
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
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtSelectFileByFilename
	stmt.Reset()
	stmt.ClearBindings()
	stmt.SetText("$filename", filename)

	if hasRow, err := stmt.Step(); err != nil {
		return types.FilesRow{}, err
	} else if !hasRow {
		return types.FilesRow{}, nil
	} else {
		var row types.FilesRow
		row.Id = int(stmt.GetInt64("id"))
		row.DirPath = stmt.GetText("dir_path")
		row.Filename = stmt.GetText("filename")
		row.DateAdded = stmt.GetText("date_added")
		row.DateModified = stmt.GetText("date_modified")
		return row, nil
	}
}

func SelectFileByFilePath(filePath string) (types.FilesRow, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtSelectFileByFilePath
	stmt.Reset()
	stmt.ClearBindings()
	stmt.SetText("$dir_path", filepath.Dir(filePath))
	stmt.SetText("$filename", filepath.Base(filePath))

	if hasRow, err := stmt.Step(); err != nil {
		return types.FilesRow{}, err
	} else if !hasRow {
		return types.FilesRow{}, nil
	} else {
		var row types.FilesRow
		row.Id = int(stmt.GetInt64("id"))
		row.DirPath = stmt.GetText("dir_path")
		row.Filename = stmt.GetText("filename")
		row.DateAdded = stmt.GetText("date_added")
		row.DateModified = stmt.GetText("date_modified")
		return row, nil
	}
}

func DeleteFileById(id int) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtDeleteFileById
	stmt.Reset()
	stmt.ClearBindings()
	stmt.SetInt64("$id", int64(id))

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to delete files row for id %d: %v", id, err)
	}
	return nil
}

func InsertIntoFiles(dirPath string, fileName string, dateAdded string, dateModified string) (int, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtInsertIntoFiles
	stmt.Reset()
	stmt.ClearBindings()
	stmt.SetText("$dir_path", dirPath)
	stmt.SetText("$filename", fileName)
	stmt.SetText("$date_added", dateAdded)
	stmt.SetText("$date_modified", dateModified)

	_, err = stmt.Step()
	if err != nil {
		return 0, fmt.Errorf("failed to insert file: %v", err)
	}

	rowId := int(DbRW.LastInsertRowID())
	return rowId, nil
}
