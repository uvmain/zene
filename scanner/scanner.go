package scanner

import (
	"log"
	"os"
	"path/filepath"
	"time"
	"zene/config"
	"zene/database"
)

func ScanMusicDirectory() error {
	log.Printf("Scanning directory: %s", config.MusicDir)
	lastScan, err := database.SelectLastScan()
	if err != nil {
		log.Printf("Failed to retrieve last scanned info: %v", err)
	}

	lastModified, err := time.Parse(time.RFC3339Nano, lastScan.DateModified)
	if err != nil {
		log.Printf("Error fetching lastModified from scans table: %v", err)
	}

	newModified := lastModified
	fileCount := 0

	scanResult := filepath.Walk(config.MusicDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error scanning directory %s: %v", path, err)
			return nil
		}
		if info.IsDir() {
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
	log.Printf("Music directory scan completed, new modified time is %s", newModified.Format(time.RFC3339Nano))

	return scanResult
}
