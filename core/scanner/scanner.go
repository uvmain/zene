package scanner

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"
	"time"
	"zene/core/art"
	"zene/core/config"
	"zene/core/database"
	"zene/core/ffprobe"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/musicbrainz"
	"zene/core/types"
)

func RunScan(ctx context.Context) (types.ScanStatus, error) {
	latestScan, err := database.GetLatestScan(ctx)
	if err != nil && err != sql.ErrNoRows {
		logger.Printf("Error getting latest scan: %v", err)
		return types.ScanStatus{
			Scanning:    false,
			Count:       0,
			FolderCount: 0,
		}, err
	}

	if latestScan.Id > 0 && latestScan.CompletedDate == "" {
		startedTime := logic.GetStringTimeFormatted(latestScan.StartedDate)
		bootTime := logic.GetBootTime()
		if latestScan.Id > 0 && startedTime.Before(bootTime) {
			// orphaned scan, set it to completed and continue to run a new one
			logger.Printf("Orphaned scan detected. Setting scan %d to completed.", latestScan.Id)
			fileAndFolderCount, err := database.GetFileAndFolderCounts(ctx)
			if err != nil {
				logger.Printf("Error getting file and folder counts in RunScan: %v", err)
			}
			database.UpdateScanProgress(ctx, latestScan.Id, database.ScanRow{
				Count:         cmp.Or(fileAndFolderCount.FileCount, latestScan.Count),
				FolderCount:   cmp.Or(fileAndFolderCount.FolderCount, latestScan.FolderCount),
				CompletedDate: logic.GetCurrentTimeFormatted(),
			})
		} else {
			return types.ScanStatus{
				Scanning:      true,
				Count:         latestScan.Count,
				FolderCount:   latestScan.FolderCount,
				StartedDate:   latestScan.StartedDate,
				Type:          latestScan.Type,
				CompletedDate: latestScan.CompletedDate,
			}, fmt.Errorf("A scan is already in progress. Please wait for it to complete before starting a new one.")
		}
	}
	newScan := database.ScanRow{
		Count:       0,
		FolderCount: 0,
		StartedDate: logic.GetCurrentTimeFormatted(),
		Type:        "full",
	}

	scanId, err := database.InsertScan(ctx, newScan)
	if err != nil {
		return types.ScanStatus{}, err
	}

	go scanMusicDirs(ctx, int(scanId))

	return types.ScanStatus{
		Scanning:    true,
		Count:       newScan.Count,
		FolderCount: newScan.FolderCount,
		StartedDate: newScan.StartedDate,
		Type:        newScan.Type,
	}, nil
}

func scanMusicDirs(ctx context.Context, scanId int) error {
	scanUpdate := database.ScanRow{
		Count:         0,
		FolderCount:   0,
		CompletedDate: logic.GetCurrentTimeFormatted(),
	}

	if ctx.Err() != nil {
		database.UpdateScanProgress(ctx, scanId, scanUpdate)
		return fmt.Errorf("scan was cancelled, context error: %v", ctx.Err())
	}

	start := time.Now()
	changesMade := false
	var err error

	defer func() { logger.Printf("Scan completed in %s", time.Since(start)) }()

	for _, musicDir := range config.MusicDirs {
		changesMade, err = scanMusicDir(ctx, musicDir)
		if err != nil {
			return fmt.Errorf("scanning music directory %s: %v", musicDir, err)
		}
	}

	if changesMade {
		musicbrainz.ClearMbCache()

		err := database.RepopulateGenreCountsTable(ctx)
		if err != nil {
			return fmt.Errorf("repopulating genre counts table: %v", err)
		}

		err = PopulateSimilarArtistsTable(ctx)
		if err != nil {
			return fmt.Errorf("populating similar artists table: %v", err)
		}

		err = PopulateTopSongsTable(ctx)
		if err != nil {
			return fmt.Errorf("populating top songs table: %v", err)
		}
	}

	fileAndFolderCount, err := database.GetFileAndFolderCounts(ctx)
	if err != nil {
		return fmt.Errorf("getting file and folder counts: %v", err)
	}

	scanUpdate = database.ScanRow{
		Count:         fileAndFolderCount.FileCount,
		FolderCount:   fileAndFolderCount.FolderCount,
		CompletedDate: logic.GetCurrentTimeFormatted(),
	}

	database.UpdateScanProgress(ctx, scanId, scanUpdate)

	return nil
}

func scanMusicDir(ctx context.Context, musicDir string) (bool, error) {
	changesMade := false

	logger.Printf("Starting scan of music dir %s", musicDir)

	// get a list of files from the filesystem
	logger.Printf("Scan: Getting list of audio files in the filesystem")
	audioFiles, err := getAudioFiles(ctx, musicDir)
	if err != nil {
		return false, fmt.Errorf("scanning music directory for audio files: %v", err)
	}

	// get a current list of files from the metadata table
	logger.Printf("Scan: Getting list of metadata in the database")
	metadataFiles, err := database.SelectTrackFilesForScanner(ctx, musicDir)
	if err != nil {
		return false, fmt.Errorf("scanning database for metadata files: %v", err)
	}

	// for each file found, either insert or update a metadata row

	existingMetadataToUpdate := []types.File{}
	newMetadataToInsert := []types.File{}

	for _, audioFile := range audioFiles {
		matchingMetadata, err := logic.FilterArray(metadataFiles, func(mf types.File) (bool, error) {
			return mf.FilePathAbs == audioFile.FilePathAbs, nil
		})
		if err != nil {
			logger.Printf("Error filtering metadata files for audio file %s: %v", audioFile.FilePathAbs, err)
			return false, fmt.Errorf("filtering metadata files: %v", err)
		}
		if len(matchingMetadata) > 0 {
			// Update existing metadata..
			metadataDateModified := matchingMetadata[0].DateModified
			if logic.GetStringTimeFormatted(metadataDateModified).Before(logic.GetStringTimeFormatted(audioFile.DateModified)) {
				// if the file's modified date is more recent than in the database
				existingMetadataToUpdate = append(existingMetadataToUpdate, audioFile)
			}
		} else {
			// Insert new metadata
			newMetadataToInsert = append(newMetadataToInsert, audioFile)
		}
	}

	if len(existingMetadataToUpdate) > 0 {
		err = upsertMetadataForFiles(ctx, existingMetadataToUpdate)
		if err != nil {
			return false, fmt.Errorf("upserting metadata for existing files: %v", err)
		}
		changesMade = true
	}

	if len(newMetadataToInsert) > 0 {
		err = upsertMetadataForFiles(ctx, newMetadataToInsert)
		if err != nil {
			return false, fmt.Errorf("upserting metadata for new files: %v", err)
		}
		changesMade = true
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
			return false, fmt.Errorf("deleting orphan metadata rows: %v", err)
		} else {
			fileCount += len(filepaths)
		}
		if fileCount > 0 {
			changesMade = true
			logger.Printf("Scan: %d orphaned metadata rows removed", fileCount)
		}
	}

	err = getAlbumArtworkForMusicDir(ctx, musicDir)
	if err != nil {
		return changesMade, fmt.Errorf("getting album artwork for music dir %s: %v", musicDir, err)
	}

	err = getArtistArtworkForMusicDir(ctx, musicDir)
	if err != nil {
		return changesMade, fmt.Errorf("getting artist artwork for music dir %s: %v", musicDir, err)
	}

	return changesMade, nil
}

func getAudioFiles(ctx context.Context, musicDir string) ([]types.File, error) {
	audioFiles, err := io.GetFiles(ctx, musicDir, config.AudioFileTypes)
	if err != nil {
		return []types.File{}, fmt.Errorf("getting slice of audio files from the filesystem: %v", err)
	}
	return audioFiles, nil
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

			fileMetadata, err := ffprobe.GetMetadata(ctx, file.FilePathAbs)
			if err != nil {
				logger.Printf("Skipping %s: error retrieving metadata from file: %v", file.FilePathAbs, err)
				return
			}

			metadata := types.Metadata{
				FilePath:            file.FilePathAbs,
				FileName:            filepath.Base(file.FilePathAbs),
				DateAdded:           logic.GetCurrentTimeFormatted(),
				DateModified:        file.DateModified,
				Format:              fileMetadata.Format,
				Duration:            fileMetadata.Duration,
				Size:                fileMetadata.Size,
				Bitrate:             fileMetadata.Bitrate,
				Title:               fileMetadata.Title,
				Artist:              fileMetadata.Artist,
				Album:               fileMetadata.Album,
				AlbumArtist:         fileMetadata.AlbumArtist,
				Genre:               fileMetadata.Genre,
				TrackNumber:         fileMetadata.TrackNumber,
				TotalTracks:         fileMetadata.TotalTracks,
				DiscNumber:          fileMetadata.DiscNumber,
				TotalDiscs:          fileMetadata.TotalDiscs,
				ReleaseDate:         fileMetadata.ReleaseDate,
				MusicBrainzArtistID: fileMetadata.MusicBrainzArtistID,
				MusicBrainzAlbumID:  fileMetadata.MusicBrainzAlbumID,
				MusicBrainzTrackID:  fileMetadata.MusicBrainzTrackID,
				Label:               fileMetadata.Label,
				Codec:               fileMetadata.Codec,
				BitDepth:            fileMetadata.BitDepth,
				SampleRate:          fileMetadata.SampleRate,
				Channels:            fileMetadata.Channels,
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
		return fmt.Errorf("upserting metadata rows: %v", err)
	}

	if len(metadataSlice) > 0 {
		logger.Printf("Scan: metadata tags for %d files upserted", len(metadataSlice))
	} else {
		logger.Printf("Scan: no metadata tags found for files")
	}
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

	albumArtists, err := database.SelectArtistsForMusicDir(ctx, musicDir)

	if err != nil {
		logger.Printf("Error fetching artists from database: %v", err)
		return err
	}
	for _, artist := range albumArtists {
		art.ImportArtForArtist(ctx, artist.MusicBrainzArtistID, artist.Artist, artist.IsAlbumArtist)
	}

	return nil
}

type ArtistsToCheck struct {
	MusicBrainzId string
	ArtistName    string
}
