package art

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	currentStatus, err := database.SelectAlbumArtByMusicBrainzAlbumId(ctx, musicBrainzAlbumId)
	if err != nil && err != sql.ErrNoRows {
		logger.Printf("Error getting album art data from database in ImportArtForAlbum: %v", err)
	}

	directory := filepath.Dir(currentStatus.FilePath)

	var foundFile string
	var fileTime time.Time

	folderFilePath := filepath.Join(directory, "folder.jpg")
	albumFileName := strings.Join([]string{albumName, "jpg"}, ".")
	albumFilePath := filepath.Join(directory, albumFileName)

	if io.FileExists(folderFilePath) {
		foundFile = folderFilePath
	} else if io.FileExists(albumFilePath) {
		foundFile = albumFilePath
	}

	artFileExists := (foundFile != "")
	if artFileExists {
		fileTime, err = io.GetChangedTime(foundFile)
		if err != nil {
			artFileExists = false
		}
	}

	audioFileChangedTime, err := io.GetChangedTime(currentStatus.FilePath)
	if err != nil {
		logger.Printf("Error getting audio file changed time in ImportArtForAlbum: %v", err)
	}

	artRowExists := (currentStatus.DateModified != "")

	if artRowExists {
		rowTime, err := time.Parse(time.RFC3339Nano, currentStatus.DateModified)
		if err != nil {
			logger.Printf("Error parsing existing row time in ImportArtForAlbum: %v", err)
		}
		// if currentStatus.dateModified is after audioFileChangedTime and fileTime, and the albumArt actually exists, break out of function
		if rowTime.After(audioFileChangedTime) && rowTime.After(fileTime) && io.FileExists(filepath.Join(config.AlbumArtFolder, musicBrainzAlbumId+".jpg")) {
			return
		} else {
			// re-import art
			logger.Printf("Scan: Album art for %s is outdated, re-importing", albumName)
			_ = io.DeleteFile(filepath.Join(config.AlbumArtFolder, musicBrainzAlbumId+".jpg"))
			if fileTime.After(rowTime) {
				getArtFromFolder(ctx, musicBrainzAlbumId, foundFile)
			} else if audioFileChangedTime.After(rowTime) {
				err = getEmbeddedArtFromTrack(ctx, musicBrainzAlbumId, currentStatus.FilePath)
			}
		}
		return
	}

	// if currentStatus.DateModified == "", we haven't imported any art yet
	if artFileExists {
		logger.Printf("Scan: Found local album art for %s, importing", albumName)
		getArtFromFolder(ctx, musicBrainzAlbumId, foundFile)
		return
	}

	// get embedded art if available
	err = getEmbeddedArtFromTrack(ctx, musicBrainzAlbumId, currentStatus.FilePath)
	if err == nil {
		logger.Printf("Scan: Found embedded album art for %s, importing", albumName)
		return
	}

	// no local image, fallback to downloading from internet
	getAlbumArtFromInternet(ctx, musicBrainzAlbumId, albumName, artistName)
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

func getEmbeddedArtFromTrack(ctx context.Context, musicBrainzAlbumId string, audioFilePath string) error {
	embeddedArt, err := ffmpeg.GetCoverArtFromTrack(ctx, audioFilePath)
	if err != nil {
		return fmt.Errorf("Error getting cover art from track in ImportArtForAlbum: %v", err)
	} else {
		getArtFromBytes(ctx, musicBrainzAlbumId, embeddedArt)
	}
	return nil
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
	tracks, err := database.GetSongsForAlbum(ctx, musicBrainzAlbumId)
	directory := filepath.Dir(tracks[0].Path)

	var foundFile string

	folderFilePath := filepath.Join(directory, "folder.jpg")
	albumFileName := strings.Join([]string{tracks[0].Album, "jpg"}, ".")
	albumFilePath := filepath.Join(directory, albumFileName)

	if io.FileExists(folderFilePath) {
		foundFile = folderFilePath
	} else if io.FileExists(albumFilePath) {
		foundFile = albumFilePath
	}

	if foundFile != "" {
		folderArtBytes, err := getBytesFromFilePath(foundFile)
		if err == nil {
			contentType := http.DetectContentType(folderArtBytes)
			localArts.FolderArt = "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(folderArtBytes)
		} else {
			localArts.FolderArt = ""
			logger.Printf("No folder art found for album %s: %v", musicBrainzAlbumId, err)
		}
	}

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
