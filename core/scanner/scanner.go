package scanner

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"
	"zene/core/art"
	"zene/core/config"
	"zene/core/database"
	"zene/core/ffprobe"
	"zene/core/globals"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/musicbrainz"
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
	changesMade := false
	start := time.Now()
	defer func() { logger.Printf("Scan completed in %s", time.Since(start)) }()
	defer func() { globals.IsScanning = false }()

	var scanResult types.ScanResponse
	for _, musicDir := range config.MusicDirs {
		changesMade, scanResult = scanMusicDir(ctx, musicDir)
		if !scanResult.Success {
			return scanResult
		}
	}

	musicbrainz.ClearMbCache()

	if changesMade {
		err := database.RepopulateGenreCountsTable(ctx)
		if err != nil {
			return scanError("Error repopulating genre counts table: %v", err)
		}
	}

	return types.ScanResponse{
		Success: true,
		Status:  "Scan run triggered",
	}
}

func scanMusicDir(ctx context.Context, musicDir string) (bool, types.ScanResponse) {
	changesMade := false

	logger.Printf("Starting scan of music dir %s", musicDir)

	// get a list of files from the filesystem
	logger.Printf("Scan: Getting list of audio files in the filesystem")
	audioFiles, err := getAudioFiles(ctx, musicDir)
	if err != nil {
		return false, scanError("Error scanning music directory for audio files: %v", err)
	}

	// get a current list of files from the metadata table
	logger.Printf("Scan: Getting list of metadata in the database")
	metadataFiles, err := database.SelectTrackFilesForScanner(ctx, musicDir)
	if err != nil {
		return false, scanError("Error scanning database for metadata files: %v", err)
	}

	// for each file found, either insert or update a metadata row
	audioFilesToInsert, err := getOutdatedOrMissing(audioFiles, metadataFiles)
	if len(audioFilesToInsert) > 0 {
		changesMade = true
	}
	if err != nil {
		return false, scanError("Error getting outdated or missing files: %v", err)
	}
	err = upsertMetadataForFiles(ctx, audioFilesToInsert)
	if err != nil {
		return false, scanError("Error upserting metadata rows: %v", err)
	}

	// for each metadata row that does not exist in the files list, delete that row
	logger.Printf("Scan: fetching orphaned metadata rows to delete")
	metadataRowsToDelete := logic.FilesInSliceOnceNotInSliceTwo(metadataFiles, audioFiles)
	fileCount := 0
	var filepaths []string
	for _, row := range metadataRowsToDelete {
		filepaths = append(filepaths, row.FilePathAbs)
	}

	if len(filepaths) == 0 {
		logger.Println("Scan: No orphaned metadata rows found")
	} else {
		logger.Printf("Scan: deleting orphaned metadata rows")
		err = database.DeleteMetadataRows(ctx, filepaths)
		if err != nil {
			scanError("Error deleting orphan metadata rows: %v", err)
		} else {
			fileCount += len(filepaths)
		}
		logger.Printf("Scan: %d orphaned metadata rows removed", fileCount)
	}

	err = getAlbumArtworkForMusicDir(ctx, musicDir)
	if err != nil {
		return changesMade, scanError(fmt.Sprintf("Error getting album artwork for music dir %s", musicDir), err)
	}

	err = getArtistArtworkForMusicDir(ctx, musicDir)
	if err != nil {
		return changesMade, scanError(fmt.Sprintf("Error getting artist artwork for music dir %s", musicDir), err)
	}

	return changesMade, types.ScanResponse{
		Success: true,
		Status:  "Scan run triggered",
	}
}

func getAudioFiles(ctx context.Context, musicDir string) ([]types.File, error) {
	audioFiles, err := io.GetFiles(ctx, musicDir, config.AudioFileTypes)
	if err != nil {
		return []types.File{}, fmt.Errorf("Error getting slice of audio files from the filesystem: %v", err)
	}
	return audioFiles, nil
}

// takes two slices of types.File
// returns entries from slice1 where either it does not exist in slice2, or the slice2 date is newer
func getOutdatedOrMissing(slice1, slice2 []types.File) ([]types.File, error) {
	slice2Map := make(map[string]string)
	for _, f := range slice2 {
		slice2Map[f.FilePathAbs] = f.DateModified
	}

	var result []types.File
	for _, file := range slice1 {
		date2Str, found := slice2Map[file.FilePathAbs]
		if !found {
			result = append(result, file)
			continue
		}

		date1, err1 := time.Parse(time.RFC3339Nano, file.DateModified)
		date2, err2 := time.Parse(time.RFC3339Nano, date2Str)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid date format: %v, %v", err1, err2)
		}

		if date1.After(date2) {
			result = append(result, file)
		}
	}

	return result, nil
}

func upsertMetadataForFiles(ctx context.Context, files []types.File) error {
	metadataSlice := make([]types.Metadata, 0, len(files))
	metadataMutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	bufferedChannel := make(chan struct{}, config.FfprobeConcurrentProcesses)

	logger.Printf("Scan: Fetching metadata tags for %d files", len(files))

	for _, file := range files {
		file := file // capture range variable
		wg.Add(1)
		bufferedChannel <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-bufferedChannel }()

			tags, err := ffprobe.GetTags(ctx, file.FilePathAbs)
			if err != nil {
				logger.Printf("Skipping %s: error retrieving tags: %v", file.FilePathAbs, err)
				return
			}

			metadata := types.Metadata{
				FilePath:            file.FilePathAbs,
				FileName:            filepath.Base(file.FilePathAbs),
				DateAdded:           logic.GetCurrentTimeFormatted(),
				DateModified:        file.DateModified,
				Format:              tags.Format,
				Duration:            tags.Duration,
				Size:                tags.Size,
				Bitrate:             tags.Bitrate,
				Title:               tags.Title,
				Artist:              tags.Artist,
				Album:               tags.Album,
				AlbumArtist:         tags.AlbumArtist,
				Genre:               tags.Genre,
				TrackNumber:         tags.TrackNumber,
				TotalTracks:         tags.TotalTracks,
				DiscNumber:          tags.DiscNumber,
				TotalDiscs:          tags.TotalDiscs,
				ReleaseDate:         tags.ReleaseDate,
				MusicBrainzArtistID: tags.MusicBrainzArtistID,
				MusicBrainzAlbumID:  tags.MusicBrainzAlbumID,
				MusicBrainzTrackID:  tags.MusicBrainzTrackID,
				Label:               tags.Label,
			}

			metadataMutex.Lock()
			metadataSlice = append(metadataSlice, metadata)
			metadataMutex.Unlock()
		}()
	}

	wg.Wait()

	logger.Printf("Scan: Upserting metadata for %d files", len(metadataSlice))
	err := database.UpsertMetadataRows(ctx, metadataSlice)
	if err != nil {
		return fmt.Errorf("Error upserting metadata rows: %v", err)
	}

	logger.Printf("Scan: metadata tags for %d files upserted", len(metadataSlice))
	return nil
}

func getAlbumArtworkForMusicDir(ctx context.Context, musicDir string) error {
	logger.Println("Getting album artwork")
	albums, err := database.SelectAllAlbumsForMusicDir(ctx, musicDir, "false", "", "")
	if err != nil {
		logger.Printf("Error fetching albums from database: %v", err)
		return err
	}
	for _, album := range albums {
		art.ImportArtForAlbum(ctx, album.MusicBrainzAlbumID, album.Album, album.Artist)
	}
	return nil
}

func getArtistArtworkForMusicDir(ctx context.Context, musicDir string) error {
	logger.Printf("Getting artist artwork for music dir %s", musicDir)

	albumArtists, err := database.SelectAlbumArtistsForMusicDir(ctx, musicDir, "", "", "", "", "", "")

	if err != nil {
		logger.Printf("Error fetching artists from database: %v", err)
		return err
	}
	for _, artist := range albumArtists {
		art.ImportArtForAlbumArtist(ctx, artist.MusicBrainzArtistID, artist.Artist)
	}

	return nil
}

func scanError(msg string, err error) types.ScanResponse {
	logger.Printf("%s: %v", msg, err)
	return types.ScanResponse{Success: false, Status: msg}
}
