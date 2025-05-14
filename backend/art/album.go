package art

import (
	"log"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"zene/database"
	"zene/io"
)

func GetArtForAlbum(musicBrainzAlbumId string, albumName string) {
	trackMetadataRows, err := database.SelectMetadataByAlbumID(musicBrainzAlbumId)
	if err != nil {
		log.Printf("Error getting track data from database: %v", err)
	}

	existingRow, err := database.SelectAlbumArtByMusicBrainzAlbumId(musicBrainzAlbumId)
	if err != nil {
		log.Printf("Error getting album art data from database: %v", err)
	}
	existingTime, err := time.Parse(time.RFC3339, existingRow.DateModified)

	directories := []string{}

	for _, trackMetadata := range trackMetadataRows {
		directory := filepath.Dir(trackMetadata.Filename)
		if !slices.Contains(directories, directory) {
			directories = append(directories, directory)
		}
	}
	directories = slices.Compact(directories)

	var foundFile string
	var foundTime time.Time

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

	if foundFile != "" {
		foundTime = io.GetChangedTime(foundFile)
		log.Printf("Found %s for %s, %v", foundFile, albumName, foundTime)
	} else {
		log.Printf("No art found %s for %s, %v", foundFile, albumName, foundTime)
	}

	if existingTime.Before(foundTime) {
		log.Printf("album art for %s is newer", albumName)
	}
}
