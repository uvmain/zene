package podcasts

import (
	"context"
	"fmt"
	"strconv"
	"zene/core/database"
	"zene/core/io"
	"zene/core/logger"
)

func DeletePodcastEpisodeById(ctx context.Context, episodeId int) error {
	var responseError string
	episode, err := database.GetPodcastEpisodeById(ctx, episodeId)
	if err != nil {
		logger.Printf("Error fetching podcast episode from database: %v", err)
		responseError += "Error fetching podcast episode from database, "
	}

	if episode.Id != "" && episode.Path != "" {
		err = io.DeleteFile(episode.Path)
		if err != nil {
			logger.Printf("Error deleting podcast episode file: %v", err)
			responseError += "Error deleting podcast episode file, "
		}
	}

	err = database.DeletePodcastEpisodeById(ctx, episode.Id)
	if err != nil {
		logger.Printf("Error deleting podcast episode from database: %v", err)
		responseError += "Error deleting podcast episode from database, "
	}

	if responseError != "" {
		responseError = responseError[:len(responseError)-2] // remove trailing comma and space
		return fmt.Errorf("%s", responseError)
	}

	return nil
}

func DeletePodcastChannelAndEpisodes(ctx context.Context, channelId int) error {
	var responseError string

	// fetch all episodes for the channel
	episodes, err := database.GetPodcastEpisodesByChannelId(ctx, channelId)
	if err != nil {
		logger.Printf("Error fetching podcast episodes from database: %v", err)
		responseError += "Error fetching podcast episodes from database, "
	}

	// delete each episode
	for _, episode := range episodes {
		episodeId, err := strconv.Atoi(episode.Id)
		if err != nil {
			logger.Printf("Error converting podcast episode id to int: %v", err)
			responseError += fmt.Sprintf("Error converting podcast episode id %s to int, ", episode.Id)
			continue
		}
		err = DeletePodcastEpisodeById(ctx, episodeId)
		if err != nil {
			logger.Printf("Error deleting podcast episode: %v", err)
			responseError += fmt.Sprintf("Error deleting podcast episode %d, ", episodeId)
		}
	}

	// delete the channel
	err = database.DeletePodcastChannel(ctx, channelId)
	if err != nil {
		logger.Printf("Error deleting podcast channel from database: %v", err)
		responseError += "Error deleting podcast channel from database, "
	}

	if responseError != "" {
		responseError = responseError[:len(responseError)-2] // remove trailing comma and space
		return fmt.Errorf("%s", responseError)
	}

	return nil
}
