package art

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"zene/config"
	"zene/database"
	"zene/io"
	"zene/musicbrainz"
)

func ImportArtForArtist(musicBrainzArtistId string, artistName string) {
	albumDirectories, err := database.SelectArtistSubDirectories(musicBrainzArtistId)
	if err != nil {
		log.Printf("Error getting artist subdirectories from database: %v", err)
	}

	existingRow, err := database.SelectArtistArtByMusicBrainzArtistId(musicBrainzArtistId)
	if err != nil {
		log.Printf("Error getting artist art data from database: %v", err)
	}
	rowTime, err := time.Parse(time.RFC3339Nano, existingRow.DateModified)

	directories := []string{}

	for _, albumDirectory := range albumDirectories {
		directory := filepath.Dir(albumDirectory)
		if !slices.Contains(directories, directory) {
			directories = append(directories, directory)
		}
	}
	directories = slices.Compact(directories)

	var foundFile string
	var fileTime time.Time

	for _, directory := range directories {
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
				log.Printf("local artist art for %s is newer, re-importing", artistName)
				getArtistArtFromFolder(musicBrainzArtistId, foundFile)
			}
		} else {
			// file hasn't been imported yet
			log.Printf("Found new artist art for %s, importing", artistName)
			getArtistArtFromFolder(musicBrainzArtistId, foundFile)
		}
	} else {
		// we've already downloaded an image
		if rowExists {
			return
		} else {
			// no local image, download from internet
			log.Printf("No artist artwork found for %s, downloading", artistName)
			getArtistArtFromInternet(musicBrainzArtistId)
		}
	}
}

func getArtistArtFromFolder(musicBrainzArtistId string, imagePath string) {
	go resizeFileAndSaveAsJPG(imagePath, filepath.Join(config.ArtistArtFolder, musicBrainzArtistId), 512)
	err := database.InsertArtistArtRow(musicBrainzArtistId, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		log.Printf("Error inserting artist art row: %v", err)
	}
}

func getArtistArtFromInternet(musicBrainzArtistId string) {
	log.Printf("fetching art for %s from musicbrainz", musicBrainzArtistId)
	artistArtUrl, err := musicbrainz.GetArtistArtUrl(musicBrainzArtistId)
	if err != nil {
		log.Printf("Failed to get artist art url for %s from musicbrainz: %v", musicBrainzArtistId, err)
		return
	}

	img, err := getImageFromInternet(artistArtUrl)
	if err != nil {
		log.Printf("Failed to get artist art image for %s from %s: %v", musicBrainzArtistId, artistArtUrl, err)
		return
	}
	go resizeImageAndSaveAsJPG(img, filepath.Join(config.ArtistArtFolder, musicBrainzArtistId), 512)

	err = database.InsertArtistArtRow(musicBrainzArtistId, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		log.Printf("Error inserting artist art row: %v", err)
	}
}

func GetArtForArtist(musicBrainzArtistId string) ([]byte, error) {
	filename := strings.Join([]string{musicBrainzArtistId, "jpg"}, ".")
	filePath, _ := filepath.Abs(filepath.Join(config.ArtistArtFolder, filename))

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Image file does not exist: %s:  %s", filePath, err)
		return nil, err
	}
	blob, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading image for filename %s: %s", filename, err)
		return nil, err
	}
	return blob, nil
}
