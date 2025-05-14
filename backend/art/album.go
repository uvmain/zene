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
	// database.InsertAlbumArtRow(musicBrainzAlbumId, time.Now().Format(time.RFC3339Nano))
}

func getArtFromInternet(musicBrainzAlbumId string) {
	// database.InsertAlbumArtRow(musicBrainzAlbumId, time.Now().Format(time.RFC3339Nano))
}
