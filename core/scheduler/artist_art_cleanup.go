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

func cleanupArtistArt(ctx context.Context) {

	artistIds, err := database.SelectArtistArtIds(ctx)
	if err != nil {
		logger.Printf("Error selecting artist art IDs: %v", err)
		return
	}

	artistArtFiles, err := io.GetFiles(ctx, config.ArtistArtFolder, []string{".jpg"})
	if err != nil {
		logger.Printf("Error getting artist art files: %v", err)
		return
	}

	if len(artistArtFiles) == 0 && len(artistIds) == 0 {
		return
	}

	files := make([]string, len(artistArtFiles))

	if len(artistArtFiles) > 0 {
		for i, file := range artistArtFiles {
			files[i] = filepath.Base(file.FilePath)
		}
	}

	// clean up orphaned artist_art rows
	artistArtRowsDeleted := 0
	for _, artistId := range artistIds {
		artistIdJpg := artistId + ".jpg"
		if !slices.Contains(files, artistIdJpg) {
			err = database.DeleteArtistArtRow(ctx, artistId)
			if err != nil {
				logger.Printf("Error deleting artist art for artist ID %s: %v", artistId, err)
			} else {
				artistArtRowsDeleted++
			}
		}
	}

	if artistArtRowsDeleted > 0 {
		logger.Printf("Artist art cleanup: deleted %d orphaned artist_art rows.", artistArtRowsDeleted)
	}

	// clean up orphaned artist art files
	artistArtFilesDeleted := 0
	for _, artFile := range files {
		artfileId := strings.TrimSuffix(artFile, filepath.Ext(artFile))
		if !slices.Contains(artistIds, artfileId) {
			logger.Printf("Deleting orphaned artist art file: %s", artFile)
			err = io.DeleteFile(filepath.Join(config.ArtistArtFolder, artFile))
			if err != nil {
				logger.Printf("Error deleting orphaned artist art file %s: %v", artFile, err)
			} else {
				artistArtFilesDeleted++
			}
		}
	}

	if artistArtFilesDeleted > 0 {
		logger.Printf("Artist art cleanup: deleted %d orphaned artist art files.", artistArtFilesDeleted)
	}
}
