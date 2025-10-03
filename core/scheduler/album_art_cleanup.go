package scheduler

import (
	"context"
	"path/filepath"
	"slices"
	"strings"
	"zene/core/config"
	"zene/core/database"
	"zene/core/io"
	"zene/core/logger"
)

func cleanupAlbumArt(ctx context.Context) {

	albumIds, err := database.SelectAlbumArtIds(ctx)
	if err != nil {
		logger.Printf("Error selecting album art IDs: %v", err)
		return
	}

	albumArtFiles, err := io.GetFiles(ctx, config.AlbumArtFolder, []string{".jpg"})
	if err != nil {
		logger.Printf("Error getting album art files: %v", err)
		return
	}

	if len(albumArtFiles) == 0 && len(albumIds) == 0 {
		return
	}

	files := make([]string, len(albumArtFiles))

	if len(albumArtFiles) > 0 {
		for i, file := range albumArtFiles {
			files[i] = filepath.Base(file.FilePath)
		}
	}

	// clean up orphaned album_art rows
	albumArtRowsDeleted := 0
	for _, albumId := range albumIds {
		albumIdJpg := albumId + ".jpg"
		if !slices.Contains(files, albumIdJpg) {
			err = database.DeleteAlbumArtRow(ctx, albumId)
			if err != nil {
				logger.Printf("Error deleting album art for album ID %s: %v", albumId, err)
			} else {
				albumArtRowsDeleted++
			}
		}
	}

	if albumArtRowsDeleted > 0 {
		logger.Printf("Album art cleanup: deleted %d orphaned album_art rows.", albumArtRowsDeleted)
	}

	// clean up orphaned album art files
	albumArtFilesDeleted := 0
	for _, artFile := range files {
		artfileId := strings.TrimSuffix(artFile, filepath.Ext(artFile))
		if !slices.Contains(albumIds, artfileId) {
			logger.Printf("Deleting orphaned album art file: %s", artFile)
			err = io.DeleteFile(filepath.Join(config.AlbumArtFolder, artFile))
			if err != nil {
				logger.Printf("Error deleting orphaned album art file %s: %v", artFile, err)
			} else {
				albumArtFilesDeleted++
			}
		}
	}

	if albumArtFilesDeleted > 0 {
		logger.Printf("Album art cleanup: deleted %d orphaned album art files.", albumArtFilesDeleted)
	}
}
