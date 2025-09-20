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
	"zene/core/deezer"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/musicbrainz"
)

func ImportArtForAlbumArtist(ctx context.Context, musicBrainzArtistId string, artistName string) {
	albumDirectories, err := database.SelectArtistSubDirectories(ctx, musicBrainzArtistId)
	if err != nil {
		logger.Printf("Error getting artist subdirectories from database: %v", err)
	}

	existingRow, err := database.SelectArtistArtByMusicBrainzArtistId(ctx, musicBrainzArtistId)
	if err != nil {
		logger.Printf("Error getting artist art data from database: %v", err)
	}
	rowTime, err := time.Parse(time.RFC3339Nano, existingRow.DateModified)

	directories := []string{}

	for _, albumDirectory := range albumDirectories {
		if err := logic.CheckContext(ctx); err != nil {
			return
		}
		directory := filepath.Dir(albumDirectory)
		if !slices.Contains(directories, directory) {
			directories = append(directories, directory)
		}
	}
	directories = slices.Compact(directories)

	var foundFile string
	var fileTime time.Time

	for _, directory := range directories {
		if err := logic.CheckContext(ctx); err != nil {
			return
		}
		folderFilePath := filepath.Join(directory, "artist.jpg")
		artistFileName := strings.Join([]string{artistName, "jpg"}, ".")
		artistFilePath := filepath.Join(directory, artistFileName)
		if io.FileExists(folderFilePath) {
			foundFile = folderFilePath
			break
		} else if io.FileExists(artistFilePath) {
			foundFile = artistFilePath
			break
		}
	}

	fileExists := (foundFile != "")
	rowExists := (existingRow.MusicbrainzArtistId != "")

	// if file exists
	if fileExists {
		// if row exists
		if rowExists {
			// if row is newer, do nothing
			if rowTime.After(fileTime) {
				return
			} else {
				// if row is older, getArtFromFolder()
				logger.Printf("local artist art for %s is newer, re-importing", artistName)
				getArtistArtFromFolder(ctx, musicBrainzArtistId, foundFile)
			}
		} else {
			// file hasn't been imported yet
			logger.Printf("Found new artist art for %s, importing", artistName)
			getArtistArtFromFolder(ctx, musicBrainzArtistId, foundFile)
		}
	} else {
		// we've already downloaded an image
		if rowExists {
			return
		} else {
			// no local image, download from internet
			logger.Printf("No artist artwork found for %s, downloading", artistName)
			getArtistArtFromInternet(ctx, musicBrainzArtistId, artistName)
		}
	}
}

func getArtistArtFromFolder(ctx context.Context, musicBrainzArtistId string, imagePath string) {
	go resizeFileAndSaveAsJPG(imagePath, filepath.Join(config.ArtistArtFolder, musicBrainzArtistId), 512)
	err := database.InsertArtistArtRow(ctx, musicBrainzArtistId, logic.GetCurrentTimeFormatted())
	if err != nil {
		logger.Printf("Error inserting artist art row: %v", err)
	}
}

func getArtistArtFromInternet(ctx context.Context, musicBrainzArtistId string, artistName string) {
	logger.Printf("fetching artist art for %s from deezer", artistName)
	artistArtUrl, err := deezer.GetArtistArtUrlWithArtistName(ctx, artistName)

	if err != nil {
		logger.Printf("failed to get artist art url for %s from deezer: %v", musicBrainzArtistId, err)
		logger.Printf("fetching art for %s from musicbrainz", musicBrainzArtistId)
		artistArtUrl, err = musicbrainz.GetArtistArtUrl(ctx, musicBrainzArtistId)
		if err != nil {
			logger.Printf("failed to get artist art url for %s from musicbrainz: %v", musicBrainzArtistId, err)
			return
		}
	}

	img, err := GetImageFromInternet(artistArtUrl)
	if err != nil {
		logger.Printf("Failed to get artist art image for %s from %s: %v", musicBrainzArtistId, artistArtUrl, err)
		return
	}
	go ResizeImageAndSaveAsJPG(img, filepath.Join(config.ArtistArtFolder, musicBrainzArtistId), 512)

	err = database.InsertArtistArtRow(ctx, musicBrainzArtistId, logic.GetCurrentTimeFormatted())
	if err != nil {
		logger.Printf("Error inserting artist art row: %v", err)
	}
}

func GetArtForArtist(ctx context.Context, musicBrainzArtistId string, size int) ([]byte, time.Time, error) {
	file_name := fmt.Sprintf("%s.jpg", musicBrainzArtistId)
	filePath, _ := filepath.Abs(filepath.Join(config.ArtistArtFolder, file_name))

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, time.Now(), fmt.Errorf("file does not exist: %s:  %s", filePath, err)
	}

	modTime := info.ModTime()

	blob, err := logic.ResizeJpegImage(ctx, filePath, size, 90)
	if err != nil {
		return nil, time.Now(), fmt.Errorf("error reading image for filepath %s: %s", filePath, err)
	}
	return blob, modTime, nil
}
