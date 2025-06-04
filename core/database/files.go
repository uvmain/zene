package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"zene/core/types"
)

func createFilesTable() {
	tableName := "files"
	schema := `CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		dir_path TEXT NOT NULL,
		file_path TEXT NOT NULL UNIQUE,
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
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT id, dir_path, filename, file_path, date_added, date_modified FROM files;`)
	defer stmt.Finalize()

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
			row.FilePath = stmt.GetText("file_path")
			row.Filename = stmt.GetText("filename")
			row.DateAdded = stmt.GetText("date_added")
			row.DateModified = stmt.GetText("date_modified")
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func SelectFileByFileId(fileId string) (types.FilesRow, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT id, dir_path, filename, file_path, date_added, date_modified FROM files WHERE id = $fileid;`)
	defer stmt.Finalize()
	stmt.SetText("$fileid", fileId)

	if hasRow, err := stmt.Step(); err != nil {
		return types.FilesRow{}, err
	} else if !hasRow {
		return types.FilesRow{}, nil
	} else {
		var row types.FilesRow
		row.Id = int(stmt.GetInt64("id"))
		row.DirPath = stmt.GetText("dir_path")
		row.FilePath = stmt.GetText("file_path")
		row.Filename = stmt.GetText("filename")
		row.DateAdded = stmt.GetText("date_added")
		row.DateModified = stmt.GetText("date_modified")
		return row, nil
	}
}

func GetFileBlob(fileId string) ([]byte, error) {
	row, err := SelectFileByFileId(fileId)
	if err != nil {
		return []byte{}, err
	}
	filePath, _ := filepath.Abs(row.FilePath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("File does not exist: %s:  %s", filePath, err)
		return nil, err
	}
	blob, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading File for filepath %s: %s", filePath, err)
		return nil, err
	}
	return blob, nil
}

func SelectFileByFilePath(filePath string) (types.FilesRow, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT id, dir_path, file_path, filename, date_added, date_modified FROM files WHERE file_path = $file_path;`)
	defer stmt.Finalize()
	stmt.SetText("$file_path", filePath)

	if hasRow, err := stmt.Step(); err != nil {
		return types.FilesRow{}, err
	} else if !hasRow {
		return types.FilesRow{}, nil
	} else {
		var row types.FilesRow
		row.Id = int(stmt.GetInt64("id"))
		row.DirPath = stmt.GetText("dir_path")
		row.Filename = stmt.GetText("filename")
		row.FilePath = stmt.GetText("file_path")
		row.DateAdded = stmt.GetText("date_added")
		row.DateModified = stmt.GetText("date_modified")
		return row, nil
	}
}

func DeleteFileById(id int) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`delete FROM files WHERE id = $id;`)
	defer stmt.Finalize()
	stmt.SetInt64("$id", int64(id))

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to delete files row for id %d: %v", id, err)
	}
	return nil
}

func InsertIntoFiles(dirPath string, fileName string, filePath string, dateAdded string, dateModified string) (int, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`INSERT INTO files (dir_path, file_path, filename, date_added, date_modified)
		VALUES ($dir_path, $file_path, $filename, $date_added, $date_modified)
		ON CONFLICT(file_path) DO UPDATE SET date_modified=excluded.date_modified
	 	WHERE excluded.date_modified>files.date_modified;`)
	defer stmt.Finalize()
	stmt.SetText("$dir_path", dirPath)
	stmt.SetText("$filename", fileName)
	stmt.SetText("$file_path", filePath)
	stmt.SetText("$date_added", dateAdded)
	stmt.SetText("$date_modified", dateModified)

	_, err = stmt.Step()
	if err != nil {
		return 0, fmt.Errorf("failed to insert file: %v", err)
	}

	rowId := int(conn.LastInsertRowID())
	return rowId, nil
}

func SelectAllFilePathsAndModTimes() (map[string]string, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
		return nil, err
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT file_path, date_modified FROM files;`)
	defer stmt.Finalize()

	fileModTimes := make(map[string]string)

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return nil, err
		} else if !hasRow {
			break
		} else {
			filePath := stmt.GetText("file_path")
			dateModified := stmt.GetText("date_modified")
			fileModTimes[filePath] = dateModified
		}
	}
	return fileModTimes, nil
}
