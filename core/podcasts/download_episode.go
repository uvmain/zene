package podcasts

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"zene/core/config"
	"zene/core/database"
	"zene/core/ffprobe"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func DownloadPodcastEpisode(ctx context.Context, episodeId string) error {
	episodeIdInt, err := strconv.Atoi(episodeId)
	if err != nil {
		return fmt.Errorf("invalid episode ID: %v", err)
	}

	episode, err := database.GetPodcastEpisodeById(ctx, episodeIdInt)

	if err != nil {
		return fmt.Errorf("getting podcast episode by ID: %v", err)
	}
	if episode.Id == "0" {
		return fmt.Errorf("podcast episode not found")
	}

	targetFilePath := filepath.Join(config.PodcastDirectory, fmt.Sprintf("%s.%s", episode.StreamId, episode.Suffix))
	fileExists := io.FileExists(targetFilePath)

	if episode.Status == types.PodcastStatusDownloading && fileExists {
		return fmt.Errorf("podcast episode is already downloading")
	}
	if episode.Status == types.PodcastStatusCompleted && fileExists {
		return fmt.Errorf("podcast episode is already completed")
	}

	go downloadEpisodeInBackground(episode)
	return nil
}

func downloadEpisodeInBackground(episode types.PodcastEpisode) {
	ctx := context.Background()

	episodeIdInt, err := strconv.Atoi(episode.Id)
	if err != nil {
		logger.Printf("Invalid episode ID: %v", err)
		return
	}

	channelIdInt, err := strconv.Atoi(episode.ChannelId)
	if err != nil {
		logger.Printf("Invalid channel ID: %v", err)
		return
	}

	err = database.UpdatePodcastEpisodeStatus(ctx, episodeIdInt, string(types.PodcastStatusDownloading))
	if err != nil {
		logger.Printf("Error updating podcast episode status: %v", err)
		return
	}

	err = database.UpdatePodcastChannelStatus(ctx, channelIdInt, string(types.PodcastStatusDownloading))
	if err != nil {
		logger.Printf("Error updating podcast channel status: %v", err)
		return
	}

	// download file from episode.SourceUrl and save to config.podcastDirectory
	targetFilePath := filepath.Join(config.PodcastDirectory, fmt.Sprintf("%s.%s", episode.StreamId, episode.Suffix))
	err = net.DownloadBinaryFile(episode.SourceUrl, targetFilePath)
	if err != nil {
		logger.Printf("Error downloading podcast episode: %v", err)
		database.UpdatePodcastEpisodeStatus(ctx, episodeIdInt, string(types.PodcastStatusError))
		database.UpdatePodcastChannelStatus(ctx, channelIdInt, string(types.PodcastStatusError))
		return
	}

	// get file size and media duration from file
	fileInfo, err := os.Stat(targetFilePath)
	if err != nil {
		logger.Printf("Error getting file info: %v", err)
		database.UpdatePodcastEpisodeStatus(ctx, episodeIdInt, string(types.PodcastStatusError))
		database.UpdatePodcastChannelStatus(ctx, channelIdInt, string(types.PodcastStatusError))
		return
	}

	ffprobeDuration, ffprobeBitrate, err := ffprobe.GetDurationAndBitrate(ctx, targetFilePath)
	if err != nil {
		logger.Printf("Error getting ffprobe metadata: %v", err)
		database.UpdatePodcastEpisodeStatus(ctx, episodeIdInt, string(types.PodcastStatusError))
		database.UpdatePodcastChannelStatus(ctx, channelIdInt, string(types.PodcastStatusError))
		return
	}
	logger.Printf("FFprobe details for episode %s: Duration: %s, Bitrate: %s", episode.Title, ffprobeDuration, ffprobeBitrate)

	err = database.AddFileDetailsToEpisode(ctx, episodeIdInt, targetFilePath, fileInfo.Size(), episode.ContentType, ffprobeDuration, ffprobeBitrate)
	if err != nil {
		logger.Printf("Error adding file details to podcast episode: %v", err)
		return
	}

	err = database.UpdatePodcastEpisodeStatus(ctx, episodeIdInt, string(types.PodcastStatusCompleted))
	if err != nil {
		logger.Printf("Error updating podcast episode status: %v", err)
		return
	}

	err = database.UpdatePodcastChannelStatus(ctx, channelIdInt, string(types.PodcastStatusCompleted))
	if err != nil {
		logger.Printf("Error updating podcast channel status: %v", err)
		return
	}

	logger.Printf("Successfully downloaded podcast episode: %s", episode.Title)
}
