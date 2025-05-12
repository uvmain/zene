package scanner

import (
	"log"
	"os"
	"path/filepath"
	"time"
	"zene/config"
	"zene/database"
	"zene/io"
)

func ScanMusicDirectory() {
	lastScan, err := database.SelectLastScan()
	if err != nil {
		log.Printf("Failed to retrieve last scanned info: %v", err)
	}

	lastModified, err := time.Parse(time.RFC3339Nano, lastScan.DateModified)
	if err != nil {
		log.Printf("Error fetching lastModified from scans table: %v", err)
	}

	err = getFiles(lastModified)
	if err != nil {
		log.Printf("Error scanning files in music directory: %v", err)
	}

	err = cleanFiles()
	if err != nil {
		log.Printf("Error cleaning file rows: %v", err)
	}
}

func getFiles(lastModified time.Time) error {
	log.Printf("Scanning directory: %s", config.MusicDir)

	newModified := lastModified
	fileCount := 0

	scanResult := filepath.WalkDir(config.MusicDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error scanning directory %s: %v", path, err)
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			log.Printf("Error retrieving file info for %s: %v", path, err)
			return nil
		}
		modTime := info.ModTime()
		if modTime.After(lastModified) {
			if modTime.After(newModified) {
				newModified = modTime
			}
			fileCount += 1
			return database.InsertIntoFiles(filepath.Dir(path), info.Name(), time.Now().Format(time.RFC3339Nano), modTime.Format(time.RFC3339Nano))
		}
		return nil
	})

	database.InsertScanRow(time.Now().Format(time.RFC3339Nano), fileCount, newModified.Format(time.RFC3339Nano))
	log.Printf("Music directory scan completed, found %d new files", fileCount)

	return scanResult
}

func cleanFiles() error {
	files, err := database.SelectAllFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		filePath := filepath.Join(file.DirPath, file.Filename)
		if !io.FileExists(filePath) {
			log.Printf("Deleting files row %d for %s", file.Id, filePath)
			database.DeleteFileById(file.Id)
		}
	}
	return nil
}
