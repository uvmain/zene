package io

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"time"
	"zene/core/config"
	"zene/core/types"

	"github.com/djherbis/times"
)

func FileExists(absoluteFilePath string) bool {
	if _, err := os.Stat(absoluteFilePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDir(directoryPath string) {
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		err := os.MkdirAll(directoryPath, 0755)
		if err != nil {
			log.Printf("Error creating directory%s: %s", directoryPath, err)
		} else {
			log.Printf("Directory created: %s", directoryPath)
		}
	} else {
		log.Printf("Directory already exists: %s", directoryPath)
	}
}

func GetChangedTime(path string) (time.Time, error) {
	t, err := times.Stat(path)
	if err != nil {
		return time.Time{}, fmt.Errorf("Error retrieving file times for %s: %v", path, err)
	}

	modTime := t.ModTime()
	var changeTime time.Time

	if t.HasChangeTime() {
		changeTime = t.ChangeTime()
	} else {
		changeTime = t.ModTime()
	}

	if changeTime.After(modTime) {
		modTime = changeTime
	}
	return modTime, nil
}

func GetFileBlob(ctx context.Context, filePath string) ([]byte, error) {
	filePathAbs, _ := filepath.Abs(filePath)

	if _, err := os.Stat(filePathAbs); os.IsNotExist(err) {
		log.Printf("File does not exist: %s:  %s", filePathAbs, err)
		return nil, err
	}
	blob, err := os.ReadFile(filePathAbs)
	if err != nil {
		log.Printf("Error reading File for filepath %s: %s", filePathAbs, err)
		return nil, err
	}

	return blob, nil
}

func GetFiles(ctx context.Context, extensions []string) ([]types.File, error) {
	files := []types.File{}
	scanError := filepath.WalkDir(config.MusicDir, func(path string, d os.DirEntry, err error) error {

		if err != nil {
			log.Printf("Error scanning file path %s: %v", path, err)
			return err
		}
		if d.IsDir() {
			return nil
		}

		if !slices.Contains(extensions, filepath.Ext(path)) {
			return nil
		}

		var foundFile types.File
		foundFile.FilePathAbs, err = filepath.Abs(path)
		if err != nil {
			log.Printf("Error getting absolute file path for %s: %v", path, err)
			return err
		}
		modTime, err := GetChangedTime(path)
		if err != nil {
			log.Printf("Error getting changed time for %s: %v", path, err)
			return err
		}
		foundFile.DateModified = modTime.Format(time.RFC3339Nano)
		files = append(files, foundFile)
		return nil
	})

	if scanError != nil {
		log.Printf("Error scanning files: %v", scanError)
	}

	return files, nil
}
