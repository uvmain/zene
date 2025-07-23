package art

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"zene/core/config"
	"zene/core/database"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/musicbrainz"
)

func ImportArtForAlbum(ctx context.Context, musicBrainzAlbumId string, albumName string) {
	trackMetadataRows, err := database.SelectTracksByAlbumId(ctx, musicBrainzAlbumId)
	if err != nil {
		logger.Printf("Error getting track data from database: %v", err)
	}

	existingRow, err := database.SelectAlbumArtByMusicBrainzAlbumId(ctx, musicBrainzAlbumId)
	if err != nil {
		logger.Printf("Error getting album art data from database: %v", err)
	}
	rowTime, err := time.Parse(time.RFC3339Nano, existingRow.DateModified)

	directories := []string{}

	for _, trackMetadata := range trackMetadataRows {
		directory := filepath.Dir(trackMetadata.FilePath)
		if !slices.Contains(directories, directory) {
			directories = append(directories, directory)
		}
	}
	directories = slices.Compact(directories)

	var foundFile string
	var fileTime time.Time

	for _, directory := range directories {
		folderFilePath := filepath.Join(directory, "folder.jpg")
		albumFileName := strings.Join([]string{albumName, "jpg"}, ".")
		albumFilePath := filepath.Join(directory, albumFileName)
		if io.FileExists(folderFilePath) {
			foundFile = folderFilePath
			break
		} else if io.FileExists(albumFilePath) {
			foundFile = albumFilePath
			break
		}
	}

	fileExists := (foundFile != "")
	rowExists := (existingRow.MusicbrainzAlbumId != "")

	// if file exists
	if fileExists {
		// if row exists
		if rowExists {
			// if row is newer, do nothing
			if rowTime.After(fileTime) {
				return
			} else {
				// if row is older, getArtFromFolder()
				logger.Printf("Scan: local album art for %s is newer, re-importing", albumName)
				getArtFromFolder(ctx, musicBrainzAlbumId, foundFile)
			}
		} else {
			// file hasn't been imported yet
			logger.Printf("Scan: Found new album art for %s, importing", albumName)
			getArtFromFolder(ctx, musicBrainzAlbumId, foundFile)
		}
	} else {
		// we've already downloaded an image
		if rowExists {
			return
		} else {
			// no local image, download from internet
			logger.Printf("Scan: No album artwork found for %s, downloading", albumName)
			getArtFromInternet(ctx, musicBrainzAlbumId)
		}
	}
}

func getArtFromFolder(ctx context.Context, musicBrainzAlbumId string, imagePath string) {
	go resizeFileAndSaveAsJPG(imagePath, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "sm"}, "_")), 64)
	go resizeFileAndSaveAsJPG(imagePath, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "md"}, "_")), 128)
	go resizeFileAndSaveAsJPG(imagePath, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "lg"}, "_")), 256)
	go resizeFileAndSaveAsJPG(imagePath, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "xl"}, "_")), 512)
	err := database.InsertAlbumArtRow(ctx, musicBrainzAlbumId, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		logger.Printf("Database: Error inserting album art row: %v", err)
	}
}

func getArtFromInternet(ctx context.Context, musicBrainzAlbumId string) {
	logger.Printf("Fetching art for %s from musicbrainz", musicBrainzAlbumId)
	albumArtUrl, err := musicbrainz.GetAlbumArtUrl(ctx, musicBrainzAlbumId)
	if err != nil {
		logger.Printf("Failed to get album art url for %s from musicbrainz: %v", musicBrainzAlbumId, err)
		return
	}

	img, err := getImageFromInternet(albumArtUrl)
	if err != nil {
		logger.Printf("Failed to get album art image for %s from %s: %v", musicBrainzAlbumId, albumArtUrl, err)
		return
	}
	go resizeImageAndSaveAsJPG(img, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "sm"}, "_")), 64)
	go resizeImageAndSaveAsJPG(img, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "md"}, "_")), 128)
	go resizeImageAndSaveAsJPG(img, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "lg"}, "_")), 256)
	go resizeImageAndSaveAsJPG(img, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "xl"}, "_")), 512)

	err = database.InsertAlbumArtRow(ctx, musicBrainzAlbumId, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		logger.Printf("Error inserting album art row: %v", err)
	}
}

func GetArtForAlbum(ctx context.Context, musicBrainzAlbumId string, size string) ([]byte, time.Time, error) {
	file_name := strings.Join([]string{musicBrainzAlbumId, size}, "_")
	file_name = strings.Join([]string{file_name, "jpg"}, ".")
	filePath, _ := filepath.Abs(filepath.Join(config.AlbumArtFolder, file_name))

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, time.Now(), fmt.Errorf("file does not exist: %s:  %s", filePath, err)
	}

	modTime := info.ModTime()

	blob, err := os.ReadFile(filePath)
	if err != nil {
		return nil, time.Now(), fmt.Errorf("error reading image for filename %s: %s", filePath, err)
	}
	return blob, modTime, nil
}
