package scanner

import (
	"context"
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

func RunScan(ctx context.Context) types.ScanResponse {
	if globals.IsScanning == true {
		return types.ScanResponse{
			Success: false,
			Status:  "Scan already in progress",
		}
	}

	globals.IsScanning = true
	log.Printf("Starting scan of music dir")

	lastScan, err := database.SelectLastScan(ctx)
	if err != nil {
		log.Printf("Failed to retrieve last scanned info: %v", err)
		return types.ScanResponse{
			Success: false,
			Status:  "Error retrieving last scanned info",
		}
	}

	lastModified, err := time.Parse(time.RFC3339Nano, lastScan.DateModified)
	if err != nil {
		lastModified = time.Now().Add(-(time.Hour * 24 * 365 * 10))
		log.Printf("Falling back to lastModified of %v: %v", lastModified, err)
	}

	files, err := getFiles(ctx, lastModified)
	if err != nil {
		log.Printf("Error scanning music directory: %v", err)
		return types.ScanResponse{
			Success: false,
			Status:  "Error scanning music directory",
		}
	}

	err = cleanFiles(ctx, files)
	if err != nil {
		log.Printf("Error cleaning file rows: %v", err)
		return types.ScanResponse{
			Success: false,
			Status:  "Error cleaning file rows",
		}
	}

	err = getAlbumArtwork(ctx)
	if err != nil {
		log.Printf("Error getting album artwork: %v", err)
		return types.ScanResponse{
			Success: false,
			Status:  "Error getting album artwork",
		}
	}

	err = getArtistArtwork(ctx)
	if err != nil {
		log.Printf("Error getting artist artwork: %v", err)
		return types.ScanResponse{
			Success: false,
			Status:  "Error getting artist artwork",
		}
	}

	globals.IsScanning = false

	log.Println("Scanner run complete")
	return types.ScanResponse{
		Success: true,
		Status:  "Scan complete",
	}
}

func getFiles(ctx context.Context, lastModified time.Time) (map[string]struct{}, error) {
	log.Printf("Scanning directory: %s", config.MusicDir)
	filesystemFilePaths := make(map[string]struct{})

	dbFileModTimes, err := database.SelectAllFilePathsAndModTimes(ctx)
	if err != nil {
		log.Printf("Failed to retrieve file modification times from database: %v", err)
		return filesystemFilePaths, err
	}

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
		// Add all encountered file paths to the map.
		filesystemFilePaths[path] = struct{}{}

		info, err := d.Info()
		if err != nil {
			log.Printf("Error retrieving file info for %s: %v", path, err)
			return err
		}

		fsModTime := io.GetChangedTime(path)
		dbModTimeStr, fileExistsInMap := dbFileModTimes[path]

		needsProcessing := false
		if !fileExistsInMap {
			needsProcessing = true
		} else {
			dbModTime, err := time.Parse(time.RFC3339Nano, dbModTimeStr)
			if err != nil {
				// Consider errors in parsing as a need to process, as the stored time is invalid.
				log.Printf("Error parsing stored modification time for %s: %v. Reprocessing.", path, err)
				needsProcessing = true
			} else if fsModTime.After(dbModTime) {
				needsProcessing = true
			}
		}

		// Also consider if the file was modified after the last scan time.
		if fsModTime.After(lastModified) {
			needsProcessing = true
		}

		if needsProcessing {
			if fsModTime.After(newModified) {
				newModified = fsModTime
			}

			fileCount += 1
			fileRowId, err := database.InsertIntoFiles(ctx, filepath.Dir(path), info.Name(), path, time.Now().Format(time.RFC3339Nano), fsModTime.Format(time.RFC3339Nano))
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

				err = database.InsertTrackMetadataRow(ctx, fileRowId, trackMetadata)
				if err != nil {
					log.Printf("Error inserting metadata for %s: %v", path, err)
					return err
				}
			}
		}
		return nil
	})

	database.InsertScanRow(ctx, time.Now().Format(time.RFC3339Nano), fileCount, newModified.Format(time.RFC3339Nano))
	log.Printf("Music directory scan completed, found %d new files", fileCount)

	if scanError != nil {
		log.Printf("Error scanning files in music directory: %v", scanError)
	}
	return filesystemFilePaths, scanError
}

func cleanFiles(ctx context.Context, filesystemFilePaths map[string]struct{}) error {
	log.Println("Cleaning orphan files")

	if filesystemFilePaths == nil || len(filesystemFilePaths) == 0 {
		log.Println("Skipping file cleaning as the list of filesystem paths is empty or nil.")
		return nil
	}

	filesFromDB, err := database.SelectAllFiles(ctx)

	if err != nil {
		return err
	}
	for _, fileFromDB := range filesFromDB {
		if _, exists := filesystemFilePaths[fileFromDB.FilePath]; !exists {
			log.Printf("Deleting files row %d for %s (not found on filesystem)", fileFromDB.Id, fileFromDB.FilePath)
			database.DeleteFileById(ctx, fileFromDB.Id)
		}
	}

	database.CleanTrackMetadata(ctx)

	return nil
}

func getAlbumArtwork(ctx context.Context) error {
	log.Println("Getting album artwork")
	albums, err := database.SelectAllAlbums(ctx, "false", "", "")
	if err != nil {
		log.Printf("Error fetching albums from database: %v", err)
		return err
	}
	for _, album := range albums {
		art.ImportArtForAlbum(ctx, album.MusicBrainzAlbumID, album.Album)
	}
	return nil
}

func getArtistArtwork(ctx context.Context) error {
	log.Println("Getting artist artwork")

	albumArtists, err := database.SelectAlbumArtists(ctx, "", "false", "", "", "")

	if err != nil {
		log.Printf("Error fetching artists from database: %v", err)
		return err
	}
	for _, artist := range albumArtists {
		art.ImportArtForAlbumArtist(ctx, artist.MusicBrainzArtistID, artist.Artist)
	}

	return nil
}
