package art

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"zene/core/config"
	"zene/core/database"
	"zene/core/deezer"
	"zene/core/ffmpeg"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/musicbrainz"
)

func ImportArtForAlbum(ctx context.Context, musicBrainzAlbumId string, albumName string, artistName string) {
	trackMetadataRows, err := database.SelectTracksByAlbumId(ctx, musicBrainzAlbumId)
	if err != nil {
		logger.Printf("Error getting track data from database in ImportArtForAlbum: %v", err)
	}

	existingRow, err := database.SelectAlbumArtByMusicBrainzAlbumId(ctx, musicBrainzAlbumId)
	if err != nil && err != sql.ErrNoRows {
		logger.Printf("Error getting album art data from database in ImportArtForAlbum: %v", err)
	}

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
			rowTime, err := time.Parse(time.RFC3339Nano, existingRow.DateModified)
			if err != nil {
				logger.Printf("Error parsing existing row time in ImportArtForAlbum: %v", err)
			}
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
			// get art from tags if available
			art, err := ffmpeg.GetCoverArtFromTrack(ctx, trackMetadataRows[0].FilePath)
			if err != nil && len(art) > 0 {
				// save art from tags
				logger.Printf("Scan: Found album artwork in tags for %s, importing", albumName)
				getArtFromBytes(ctx, musicBrainzAlbumId, art)
			} else {
				// no local image, fallback to downloading from internet
				logger.Printf("Scan: No album artwork found for %s, downloading", albumName)
				getAlbumArtFromInternet(ctx, musicBrainzAlbumId, albumName, artistName)
			}
		}
	}
}

func getArtFromBytes(ctx context.Context, musicBrainzAlbumId string, artBytes []byte) {
	go resizeBytesAndSaveAsJPG(artBytes, filepath.Join(config.AlbumArtFolder, musicBrainzAlbumId), 512)
	err := database.UpsertAlbumArtRow(ctx, musicBrainzAlbumId)
	if err != nil {
		logger.Printf("Database: Error inserting album art row: %v", err)
	}
}

func getArtFromFolder(ctx context.Context, musicBrainzAlbumId string, imagePath string) {
	go resizeFileAndSaveAsJPG(imagePath, filepath.Join(config.AlbumArtFolder, musicBrainzAlbumId), 512)
	err := database.UpsertAlbumArtRow(ctx, musicBrainzAlbumId)
	if err != nil {
		logger.Printf("Database: Error inserting album art row: %v", err)
	}
}

func getAlbumArtFromInternet(ctx context.Context, musicBrainzAlbumId string, albumName string, artistName string) {
	logger.Printf("Fetching art for %s from deezer", musicBrainzAlbumId)
	albumArtUrl, err := deezer.GetAlbumArtUrlWithArtistNameAndAlbumName(ctx, artistName, albumName)
	if err != nil {
		logger.Printf("Failed to get album art url for %s from deezer: %v. Fetching from musicbrainz", musicBrainzAlbumId, err)
		albumArtUrl, err = musicbrainz.GetAlbumArtUrl(ctx, musicBrainzAlbumId)
		if err != nil {
			logger.Printf("Failed to get album art url for %s from musicbrainz: %v", musicBrainzAlbumId, err)
			return
		}
	}

	img, err := GetImageFromInternet(albumArtUrl)
	if err != nil {
		logger.Printf("Failed to get album art image for %s from %s: %v", musicBrainzAlbumId, albumArtUrl, err)
		return
	}
	go ResizeImageAndSaveAsJPG(img, filepath.Join(config.AlbumArtFolder, musicBrainzAlbumId), 512)

	err = database.UpsertAlbumArtRow(ctx, musicBrainzAlbumId)
	if err != nil {
		logger.Printf("Error inserting album art row: %v", err)
	}
}

func GetArtForTrack(ctx context.Context, musicBrainzTrackId string, size int) ([]byte, time.Time, error) {
	albumId, err := database.SelectAlbumIdByTrackId(ctx, musicBrainzTrackId)
	if err != nil {
		logger.Printf("Error getting album ID for track %s: %v", musicBrainzTrackId, err)
		return nil, time.Now(), fmt.Errorf("album not found: %s", musicBrainzTrackId)
	}
	return GetArtForAlbum(ctx, albumId, size)
}

func GetArtForAlbum(ctx context.Context, musicBrainzAlbumId string, size int) ([]byte, time.Time, error) {
	// prevent path traversal
	if strings.Contains(musicBrainzAlbumId, "/") || strings.Contains(musicBrainzAlbumId, "\\") || strings.Contains(musicBrainzAlbumId, "..") {
		return nil, time.Now(), fmt.Errorf("invalid album ID")
	}
	file_name := musicBrainzAlbumId + ".jpg"
	filePath, _ := filepath.Abs(filepath.Join(config.AlbumArtFolder, file_name))

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

type LocalArts struct {
	FolderArt   string `json:"folderArt"`
	EmbeddedArt string `json:"embeddedArt"`
}

func GetLocalArtAsBase64(ctx context.Context, musicBrainzAlbumId string) (LocalArts, error) {
	var localArts LocalArts

	folderArtBlob, _, err := GetArtForAlbum(ctx, musicBrainzAlbumId, 512)
	if err == nil {
		contentType := http.DetectContentType(folderArtBlob)
		localArts.FolderArt = "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(folderArtBlob)
	} else {
		localArts.FolderArt = ""
		logger.Printf("No folder art found for album %s: %v", musicBrainzAlbumId, err)
	}

	tracks, err := database.GetSongsForAlbum(ctx, musicBrainzAlbumId)

	trackArtBlob, err := ffmpeg.GetCoverArtFromTrack(ctx, tracks[0].Path)
	if err == nil {
		contentType := http.DetectContentType(trackArtBlob)
		localArts.EmbeddedArt = "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(trackArtBlob)
	} else {
		localArts.EmbeddedArt = ""
		logger.Printf("No embedded art found for album %s: %v", musicBrainzAlbumId, err)
	}

	return localArts, nil
}
