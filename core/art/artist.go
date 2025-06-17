package art

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"zene/core/config"
	"zene/core/database"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/musicbrainz"
	"zene/core/types"
)

func ImportArtForArtists(ctx context.Context, artists []types.ArtistResponse) {
	for _, artist := range artists {
		if artist.MusicBrainzArtistID == "" {
			logger.Printf("Skipping artist with empty musicbrainz ID: %s", artist.Artist)
			continue
		}
		logger.Printf("Importing art for artist: %s (%s)", artist.Artist, artist.MusicBrainzArtistID)
		getArtistArtFromInternet(ctx, artist.MusicBrainzArtistID)
	}
	logger.Println("Finished importing art for artists")
}

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
			getArtistArtFromInternet(ctx, musicBrainzArtistId)
		}
	}
}

func getArtistArtFromFolder(ctx context.Context, musicBrainzArtistId string, imagePath string) {
	go resizeFileAndSaveAsJPG(imagePath, filepath.Join(config.ArtistArtFolder, musicBrainzArtistId), 512)
	err := database.InsertArtistArtRow(ctx, musicBrainzArtistId, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		logger.Printf("Error inserting artist art row: %v", err)
	}
}

func getArtistArtFromInternet(ctx context.Context, musicBrainzArtistId string) {
	logger.Printf("fetching art for %s from musicbrainz", musicBrainzArtistId)
	artistArtUrl, err := musicbrainz.GetArtistArtUrl(ctx, musicBrainzArtistId)
	if err != nil {
		logger.Printf("Failed to get artist art url for %s from musicbrainz: %v", musicBrainzArtistId, err)
		return
	}

	img, err := getImageFromInternet(artistArtUrl)
	if err != nil {
		logger.Printf("Failed to get artist art image for %s from %s: %v", musicBrainzArtistId, artistArtUrl, err)
		return
	}
	go resizeImageAndSaveAsJPG(img, filepath.Join(config.ArtistArtFolder, musicBrainzArtistId), 512)

	err = database.InsertArtistArtRow(ctx, musicBrainzArtistId, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		logger.Printf("Error inserting artist art row: %v", err)
	}
}

func GetArtForArtist(ctx context.Context, musicBrainzArtistId string) ([]byte, error) {
	file_name := strings.Join([]string{musicBrainzArtistId, "jpg"}, ".")
	filePath, _ := filepath.Abs(filepath.Join(config.ArtistArtFolder, file_name))

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.Printf("Image file does not exist: %s:  %s", filePath, err)
		return nil, err
	}
	blob, err := os.ReadFile(filePath)
	if err != nil {
		logger.Printf("Error reading image for file_name %s: %s", file_name, err)
		return nil, err
	}
	return blob, nil
}
