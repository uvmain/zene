package database

import (
	"context"
	"fmt"
	"log"
)

func GetFileBlob(ctx context.Context, fileId string) ([]byte, error) {
	// row, err := SelectFileByFileId(ctx, fileId)
	// if err != nil {
	// 	return []byte{}, err
	// }
	// filePath, _ := filepath.Abs(row.FilePath)

	// if _, err := os.Stat(filePath); os.IsNotExist(err) {
	// 	log.Printf("File does not exist: %s:  %s", filePath, err)
	// 	return nil, err
	// }
	// blob, err := os.ReadFile(filePath)
	// if err != nil {
	// 	log.Printf("Error reading File for filepath %s: %s", filePath, err)
	// 	return nil, err
	// }

	// return blob, nil
	return []byte{}, nil
}

func DeleteFileById(ctx context.Context, id int) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in DeleteFileById: %v", err)
		return err
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

func CleanTrackMetadata(ctx context.Context) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in DeleteFileById: %v", err)
		return err
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`DELETE FROM track_metadata
		WHERE NOT EXISTS (
			SELECT 1
			FROM files
			WHERE files.id = track_metadata.file_id
		);`)
	defer stmt.Finalize()

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to delete track_metadata during clean: %v", err)
	}
	return nil
}

func InsertIntoFiles(ctx context.Context, dirPath string, fileName string, filePath string, dateAdded string, dateModified string) (int, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in InsertIntoFiles: %v", err)
		return 0, err
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`INSERT INTO files (dir_path, file_path, file_name, date_added, date_modified)
		VALUES ($dir_path, $file_path, $file_name, $date_added, $date_modified)
		ON CONFLICT(file_path) DO UPDATE SET date_modified=excluded.date_modified
	 	WHERE excluded.date_modified>files.date_modified;`)
	defer stmt.Finalize()
	stmt.SetText("$dir_path", dirPath)
	stmt.SetText("$file_name", fileName)
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

func SelectAllFilePathsAndModTimes(ctx context.Context) (map[string]string, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in SelectAllFilePathsAndModTimes: %v", err)
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
