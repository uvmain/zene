package art

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"zene/core/config"
	"zene/core/database"
	"zene/core/io"
	"zene/core/musicbrainz"
)

func ImportArtForAlbum(musicBrainzAlbumId string, albumName string) {
	trackMetadataRows, err := database.SelectTracksByAlbumID(musicBrainzAlbumId)
	if err != nil {
		log.Printf("Error getting track data from database: %v", err)
	}

	existingRow, err := database.SelectAlbumArtByMusicBrainzAlbumId(musicBrainzAlbumId)
	if err != nil {
		log.Printf("Error getting album art data from database: %v", err)
	}
	rowTime, err := time.Parse(time.RFC3339Nano, existingRow.DateModified)

	directories := []string{}

	for _, trackMetadata := range trackMetadataRows {
		directory := filepath.Dir(trackMetadata.Filename)
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
				log.Printf("local album art for %s is newer, re-importing", albumName)
				getArtFromFolder(musicBrainzAlbumId, foundFile)
			}
		} else {
			// file hasn't been imported yet
			log.Printf("Found new album art for %s, importing", albumName)
			getArtFromFolder(musicBrainzAlbumId, foundFile)
		}
	} else {
		// we've already downloaded an image
		if rowExists {
			return
		} else {
			// no local image, download from internet
			log.Printf("No album artwork found for %s, downloading", albumName)
			getArtFromInternet(musicBrainzAlbumId)
		}
	}
}

func getArtFromFolder(musicBrainzAlbumId string, imagePath string) {
	go resizeFileAndSaveAsJPG(imagePath, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "sm"}, "_")), 64)
	go resizeFileAndSaveAsJPG(imagePath, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "md"}, "_")), 128)
	go resizeFileAndSaveAsJPG(imagePath, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "lg"}, "_")), 256)
	go resizeFileAndSaveAsJPG(imagePath, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "xl"}, "_")), 512)
	err := database.InsertAlbumArtRow(musicBrainzAlbumId, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		log.Printf("Error inserting album art row: %v", err)
	}
}

func getArtFromInternet(musicBrainzAlbumId string) {
	log.Printf("fetching art for %s from musicbrainz", musicBrainzAlbumId)
	albumArtUrl, err := musicbrainz.GetAlbumArtUrl(musicBrainzAlbumId)
	if err != nil {
		log.Printf("Failed to get album art url for %s from musicbrainz: %v", musicBrainzAlbumId, err)
		return
	}

	img, err := getImageFromInternet(albumArtUrl)
	if err != nil {
		log.Printf("Failed to get album art image for %s from %s: %v", musicBrainzAlbumId, albumArtUrl, err)
		return
	}
	go resizeImageAndSaveAsJPG(img, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "sm"}, "_")), 64)
	go resizeImageAndSaveAsJPG(img, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "md"}, "_")), 128)
	go resizeImageAndSaveAsJPG(img, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "lg"}, "_")), 256)
	go resizeImageAndSaveAsJPG(img, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "xl"}, "_")), 512)

	err = database.InsertAlbumArtRow(musicBrainzAlbumId, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		log.Printf("Error inserting album art row: %v", err)
	}
}

func GetArtForAlbum(musicBrainzAlbumId string, size string) ([]byte, error) {
	filename := strings.Join([]string{musicBrainzAlbumId, size}, "_")
	filename = strings.Join([]string{filename, "jpg"}, ".")
	filePath, _ := filepath.Abs(filepath.Join(config.AlbumArtFolder, filename))

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
