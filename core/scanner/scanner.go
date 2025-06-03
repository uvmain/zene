package scanner

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"time"
	"zene/core/art"
	"zene/core/config"
	"zene/core/database"
	"zene/core/ffprobe"
	"zene/core/globals"
	"zene/core/io"
	"zene/core/types"
)

func RunScan() types.ScanResponse {
	if globals.Syncing == true {
		return types.ScanResponse{
			Success: false,
			Status:  "Scan already in progress",
		}
	}

	globals.Syncing = true
	log.Printf("Starting scan of music dir")

	lastScan, err := database.SelectLastScan()
	if err != nil {
		log.Printf("Failed to retrieve last scanned info: %v", err)
		return types.ScanResponse{
			Success: false,
			Status:  "Error retrieving last scanned info",
		}
	}

	lastModified, err := time.Parse(time.RFC3339Nano, lastScan.DateModified)
	if err != nil {
		log.Printf("Error fetching lastModified from scans table: %v", err)
	}

	err = getFiles(lastModified)
	if err != nil {
		log.Printf("Error scanning music directory: %v", err)
		return types.ScanResponse{
			Success: false,
			Status:  "Error scanning music directory",
		}
	}

	err = cleanFiles()
	if err != nil {
		log.Printf("Error cleaning file rows: %v", err)
		return types.ScanResponse{
			Success: false,
			Status:  "Error cleaning file rows",
		}
	}

	err = getAlbumArtwork()
	if err != nil {
		log.Printf("Error getting album artwork: %v", err)
		return types.ScanResponse{
			Success: false,
			Status:  "Error getting album artwork",
		}
	}

	err = getArtistArtwork()
	if err != nil {
		log.Printf("Error getting artist artwork: %v", err)
		return types.ScanResponse{
			Success: false,
			Status:  "Error getting artist artwork",
		}
	}

	globals.Syncing = false

	log.Println("Scanner run complete")
	return types.ScanResponse{
		Success: true,
		Status:  "Scan complete",
	}
}

func getFiles(lastModified time.Time) error {
	log.Printf("Scanning directory: %s", config.MusicDir)

	newModified := lastModified
	fileCount := 0

	scanError := filepath.WalkDir(config.MusicDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error scanning directory %s: %v", path, err)
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			log.Printf("Error retrieving file info for %s: %v", path, err)
			return err
		}

		modTime := io.GetChangedTime(path)

		row, err := database.SelectFileByFilePath(path)
		if err != nil {
			log.Printf("Error selecting file by path %s: %v", path, err)
		}
		rowExists := false
		if row.Id != 0 {
			rowExists = true
		}

		if modTime.After(lastModified) || !rowExists {
			if modTime.After(newModified) {
				newModified = modTime
			}

			fileCount += 1
			fileRowId, err := database.InsertIntoFiles(filepath.Dir(path), info.Name(), path, time.Now().Format(time.RFC3339Nano), modTime.Format(time.RFC3339Nano))
			if err != nil {
				log.Printf("Error inserting files row for %s: %v", path, err)
				return err
			}

			if slices.Contains(config.AudioFileTypes, filepath.Ext(path)) {
				trackMetadata, err := ffprobe.GetTags(path)
				if err != nil {
					log.Printf("Error retrieving tags for %s: %v", path, err)
					return err
				}

				err = database.InsertTrackMetadataRow(fileRowId, trackMetadata)
				if err != nil {
					log.Printf("Error inserting metadata for %s: %v", path, err)
					return err
				}
			}
		}
		return nil
	})

	database.InsertScanRow(time.Now().Format(time.RFC3339Nano), fileCount, newModified.Format(time.RFC3339Nano))
	log.Printf("Music directory scan completed, found %d new files", fileCount)

	if scanError != nil {
		log.Printf("Error scanning files in music directory: %v", scanError)
	}
	return scanError
}

func cleanFiles() error {
	log.Println("Cleaning orphan files")
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

func getAlbumArtwork() error {
	log.Println("Getting album artwork")
	albums, err := database.SelectAllAlbums("false", "", "")
	if err != nil {
		log.Printf("Error fetching albums from database: %v", err)
		return err
	}
	for _, album := range albums {
		art.ImportArtForAlbum(album.MusicBrainzAlbumID, album.Album)
	}
	return nil
}

func getArtistArtwork() error {
	log.Println("Getting artist artwork")

	albumArtists, err := database.SelectAlbumArtists("", "false", "", "", "")

	if err != nil {
		log.Printf("Error fetching artists from database: %v", err)
		return err
	}
	for _, artist := range albumArtists {
		art.ImportArtForAlbumArtist(artist.MusicBrainzArtistID, artist.Artist)
	}

	return nil
}
