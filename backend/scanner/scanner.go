package scanner

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"time"
	"zene/art"
	"zene/config"
	"zene/database"
	"zene/ffprobe"
	"zene/globals"
	"zene/io"

	"github.com/djherbis/times"
)

func ScanMusicDirectory() {
	globals.Syncing = true

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

	globals.Syncing = false
}

func getFiles(lastModified time.Time) error {
	log.Printf("Scanning directory: %s", config.MusicDir)

	newModified := lastModified
	fileCount := 0
	var dirModTime time.Time
	var checkedAlbums []string

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

		t, err := times.Stat(path)
		if err != nil {
			log.Printf("Error retrieving file times for %s: %v", path, err)
		}

		modTime := t.ModTime()
		changeTime := t.ChangeTime()
		if changeTime.After(modTime) {
			modTime = changeTime
		}

		row, err := database.SelectFileByFilePath(path)
		rowExists := false
		if row.Id != 0 {
			rowExists = true
		}

		if modTime.After(lastModified) || !rowExists {
			if modTime.After(newModified) {
				newModified = modTime
			}
			if dirModTime.After(newModified) {
				newModified = dirModTime
			}
			fileCount += 1
			fileRowId, err := database.InsertIntoFiles(filepath.Dir(path), info.Name(), time.Now().Format(time.RFC3339Nano), modTime.Format(time.RFC3339Nano))
			if err != nil {
				log.Printf("Error inserting files row for %s: %v", path, err)
				return nil
			}

			if slices.Contains(config.AudioFileTypes, filepath.Ext(path)) {
				trackMetadata, err := ffprobe.GetTags(path)
				if err != nil {
					log.Printf("Error retrieving tags for %s: %v", path, err)
					return nil
				}

				err = database.InsertTrackMetadataRow(fileRowId, trackMetadata)
				if err != nil {
					log.Printf("Error inserting metadata for %s: %v", path, err)
					return nil
				}

				if !slices.Contains(checkedAlbums, trackMetadata.MusicBrainzAlbumID) {
					checkedAlbums = append(checkedAlbums, trackMetadata.MusicBrainzAlbumID)
					art.GetArtForAlbum(trackMetadata.MusicBrainzAlbumID)
				}
			}
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
