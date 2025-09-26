package scheduler

import (
	"context"
	"path/filepath"
	"strconv"
	"zene/core/config"
	"zene/core/database"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/types"
)

func cleanupMissingPodcasts(ctx context.Context) {

	episodes, err := database.GetPodcastEpisodesUserless(ctx)
	if err != nil {
		logger.Printf("Error selecting stale audio cache entries: %v", err)
	}

	// remove files for missing episodes in DB

	podcastEpisodeFiles, err := io.GetFiles(ctx, config.PodcastDirectory, []string{})
	logger.Printf("Found %d podcast episode files in podcast directory", len(podcastEpisodeFiles))
	if err != nil {
		logger.Printf("Failed to get podcast episode files in cleanupMissingPodcasts: %v", err)
		return
	}
	for _, file := range podcastEpisodeFiles {
		found := false
		for _, episode := range episodes {
			if episode.Path == file.FilePath {
				// found matching episode
				found = true
				break
			}
		}
		if !found {
			// no matching episode was found, delete the file
			logger.Printf("Removing orphan podcast episode file: %s", file.FilePath)
			err := io.DeleteFile(file.FilePath)
			if err != nil {
				logger.Printf("Failed to delete orphan podcast episode file: %v", err)
			}
		}
	}

	// update episode DB status for missing episode files

	logger.Printf("Checking %d podcast episodes for missing files", len(episodes))

	for _, episode := range episodes {
		if episode.StreamId != "" && episode.Status == types.PodcastStatusCompleted {
			// check if the file exists
			fileName := episode.StreamId + "." + episode.Suffix
			filePath := filepath.Join(config.PodcastDirectory, fileName)
			exists := io.FileExists(filePath)
			if !exists {
				logger.Printf("Podcast episode file missing, updating DB status to not downloaded: %s", filePath)
				episodeIdInt, err := strconv.Atoi(episode.Id)
				if err != nil {
					logger.Printf("Failed to convert episode ID to int: %v", err)
					return
				}
				err = database.UpdatePodcastEpisodeStatus(ctx, episodeIdInt, string(types.PodcastStatusNew))
				if err != nil {
					logger.Printf("Failed to update podcast episode status: %v", err)
				}
			}
		}
	}
}
