package scanner

import (
	"context"
	"fmt"
	"path/filepath"
	"time"
	"zene/core/art"
	"zene/core/config"
	"zene/core/database"
	"zene/core/ffprobe"
	"zene/core/globals"
	"zene/core/io"
	"zene/core/logger"
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
	logger.Printf("Starting scan of music dir")
	start := time.Now()
	defer func() { logger.Printf("Scan completed in %s", time.Since(start)) }()
	defer func() { globals.IsScanning = false }()

	// get a list of files from the filesystem
	logger.Printf("Scan: Getting list of audio files in the filesystem")
	audioFiles, err := getAudioFiles(ctx)
	if err != nil {
		return scanError("Error scanning music directory for audio files: %v", err)
	}

	// get a current list of files from the metadata table
	logger.Printf("Scan: Getting list of metadata in the database")
	metadataFiles, err := database.SelectTrackFilesForScanner(ctx)
	if err != nil {
		return scanError("Error scanning database for metadata files: %v", err)

	}

	// for each file found, either insert or update a metadata row
	logger.Printf("Scan: Upserting metadata into database")
	audioFilesToInsert, err := getOutdatedOrMissing(audioFiles, metadataFiles)
	if err != nil {
		return scanError("Error getting outdated or missing files: %v", err)
	}
	fileCount := 0
	for _, audioFile := range audioFilesToInsert {
		err = upsertMetadataForFile(ctx, audioFile)
		if err != nil {
			scanError("Error inserting or updating metadata row: %v", err)
		} else {
			fileCount += 1
		}
	}
	logger.Printf("Scan: %d metadata rows upserted", fileCount)

	// for each metadata row that does not exist in the files list, delete that row
	logger.Printf("Scan: Cleaning orphaned metadata rows")
	metadataRowsToDelete := filesInSliceOnceNotInSliceTwo(metadataFiles, audioFiles)
	fileCount = 0
	for _, metadataRow := range metadataRowsToDelete {
		err = database.DeleteMetadataRow(ctx, metadataRow.FilePathAbs)
		if err != nil {
			scanError("Error deleting orphan metadata row: %v", err)
		} else {
			fileCount += 1
		}
	}
	logger.Printf("Scan: %d orphaned metadata rows removed", fileCount)

	err = getAlbumArtwork(ctx)
	if err != nil {
		return scanError("Error getting album artwork: %v", err)
	}

	err = getArtistArtwork(ctx)
	if err != nil {
		return scanError("Error getting artist artwork: %v", err)
	}

	musicbrainz.ClearMbCache()

	return types.ScanResponse{
		Success: true,
		Status:  "Scan run triggered",
	}
}

func getAudioFiles(ctx context.Context) ([]types.File, error) {
	audioFiles, err := io.GetFiles(ctx, config.AudioFileTypes)
	if err != nil {
		return []types.File{}, fmt.Errorf("Error getting slice of audio files from the filesystem: %v", err)
	}
	return audioFiles, nil
}

func filesInSliceOnceNotInSliceTwo(slice1, slice2 []types.File) []types.File {
	slice2Map := make(map[string]bool)
	for _, f := range slice2 {
		slice2Map[f.FilePathAbs] = true
	}

	var diff []types.File
	for _, f := range slice1 {
		if !slice2Map[f.FilePathAbs] {
			diff = append(diff, f)
		}
	}

	return diff
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

func upsertMetadataForFile(ctx context.Context, file types.File) error {
	metadata := types.Metadata{}
	tags, err := ffprobe.GetTags(ctx, file.FilePathAbs)
	if err != nil {
		return fmt.Errorf("Error retrieving tags for %s: %v", file.FilePathAbs, err)
	}

	metadata.FilePath = file.FilePathAbs
	metadata.FileName = filepath.Base(file.FilePathAbs)
	metadata.DateAdded = time.Now().Format(time.RFC3339Nano)
	metadata.DateModified = file.DateModified
	metadata.Format = tags.Format
	metadata.Duration = tags.Duration
	metadata.Size = tags.Size
	metadata.Bitrate = tags.Bitrate
	metadata.Title = tags.Title
	metadata.Artist = tags.Artist
	metadata.Album = tags.Album
	metadata.AlbumArtist = tags.AlbumArtist
	metadata.Genre = tags.Genre
	metadata.TrackNumber = tags.TrackNumber
	metadata.TotalTracks = tags.TotalTracks
	metadata.DiscNumber = tags.DiscNumber
	metadata.TotalDiscs = tags.TotalDiscs
	metadata.ReleaseDate = tags.ReleaseDate
	metadata.MusicBrainzArtistID = tags.MusicBrainzArtistID
	metadata.MusicBrainzAlbumID = tags.MusicBrainzAlbumID
	metadata.MusicBrainzTrackID = tags.MusicBrainzTrackID
	metadata.Label = tags.Label

	err = database.InsertMetadataRow(ctx, metadata)
	if err != nil {
		return fmt.Errorf("Error inserting metadata for %s: %v", file.FilePathAbs, err)
	}
	return nil
}

func getAlbumArtwork(ctx context.Context) error {
	logger.Println("Getting album artwork")
	albums, err := database.SelectAllAlbums(ctx, "false", "", "")
	if err != nil {
		logger.Printf("Error fetching albums from database: %v", err)
		return err
	}
	for _, album := range albums {
		art.ImportArtForAlbum(ctx, album.MusicBrainzAlbumID, album.Album)
	}
	return nil
}

func getArtistArtwork(ctx context.Context) error {
	logger.Println("Getting artist artwork")

	albumArtists, err := database.SelectAlbumArtists(ctx, "", "", "", "", "", "")

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
