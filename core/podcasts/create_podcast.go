package podcasts

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"zene/core/art"
	"zene/core/config"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"

	"github.com/mmcdole/gofeed"
)

func CreateNewPodcastFromFeedUrl(ctx context.Context, feedUrl string) error {
	user, err := database.GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("checking user with context: %v", err)
	}
	if !user.PodcastRole {
		return fmt.Errorf("user does not have permissions to create podcasts")
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedUrl)
	if err != nil {
		return fmt.Errorf("parsing feed URL: %v", err)
	}

	coverArt, err := SavePodcastImage(ctx, feed.Image.URL)
	if err != nil {
		return fmt.Errorf("saving podcast channel image: %v", err)
	}

	podcastId, err := database.CreatePodcastChannel(ctx, feedUrl, feed.Title, feed.Description, feed.Image.URL, coverArt, "", feed.Categories)
	if err != nil {
		return fmt.Errorf("creating podcast channel: %v", err)
	}

	logger.Printf("Created podcast channel '%s' with ID %d, fetching episodes...", feed.Title, podcastId)

	go createPodcastEpisodesForFeed(ctx, feed, podcastId)

	return nil
}

func createPodcastEpisodesForFeed(ctx context.Context, feed *gofeed.Feed, podcastId int) error {
	existingEpisodes, err := database.GetPodcastEpisodesByChannelId(ctx, podcastId)
	if err != nil {
		return fmt.Errorf("getting existing podcast episodes: %v", err)
	}

	existingGuids := []string{}
	for _, episode := range existingEpisodes {
		existingGuids = append(existingGuids, episode.StreamId)
	}

	newEpisodes := []*gofeed.Item{}
	for _, item := range feed.Items {
		// if item.GUID does not exist in existingEpisodes
		if !slices.Contains(existingGuids, item.GUID) {
			// then create a new episode
			newEpisodes = append(newEpisodes, item)
		}
	}

	podcastEpisodes := make([]types.PodcastEpisodeRow, 0, len(newEpisodes))
	for _, item := range newEpisodes {

		authors := ""
		for _, author := range item.Authors {
			authors += author.Name + ", "
		}
		authors = authors[:len(authors)-2] // remove trailing comma and space

		coverArt, err := SavePodcastImage(ctx, item.Image.URL)
		if err != nil {
			return fmt.Errorf("saving podcast episode cover art: %v", err)
		}

		episodeLink := item.Enclosures[0]
		episodeDuration, err := strconv.Atoi(episodeLink.Length)
		if err != nil {
			logger.Printf("error parsing episode duration: %v", err)
			episodeDuration = 0
		}

		episode := types.PodcastEpisodeRow{
			ChannelId:   strconv.Itoa(podcastId),
			Guid:        item.GUID,
			Title:       item.Title,
			Artist:      authors,
			Year:        item.PublishedParsed.Format("2006"),
			CoverArt:    coverArt,
			ContentType: episodeLink.Type,
			Suffix:      filepath.Ext(episodeLink.URL),
			Duration:    episodeDuration,
			// BitRate:     item.Enclosure.Bitrate,
			Description: item.Description,
			PublishDate: logic.FormatTimeAsString(*item.PublishedParsed),
			Status:      string(types.PodcastStatusNew),
			// FilePath:    item.Enclosure.Url,
			// StreamId:    item.StreamId,
			CreatedAt: logic.GetCurrentTimeFormatted(),
		}
		podcastEpisodes = append(podcastEpisodes, episode)
	}

	if err := database.InsertPodcastEpisodes(podcastEpisodes); err != nil {
		logger.Printf("Error inserting podcast episodes: %v", err)
		return fmt.Errorf("inserting podcast episodes: %v", err)
	}

	database.UpdatePodcastChannelLastRefresh(podcastId)

	return nil
}

func SavePodcastImage(ctx context.Context, imageUrl string) (string, error) {
	coverArtUuid, err := logic.GenerateNewApiKey()
	if err != nil {
		return "", fmt.Errorf("generating coverArt UUID: %v", err)
	}

	image, err := art.GetImageFromInternet(imageUrl)
	if err != nil {
		return "", fmt.Errorf("getting image from internet: %v", err)
	}

	outputPath := filepath.Join(config.PodcastArtFolder, fmt.Sprintf("%s.jpg", coverArtUuid))

	err = art.ResizeImageAndSaveAsJPG(image, outputPath, 600)
	if err != nil {
		return "", fmt.Errorf("resizing and saving image: %v", err)
	}

	return coverArtUuid, nil
}

func RefreshAllPodcasts(ctx context.Context) error {
	logger.Printf("Refreshing all podcasts...")
	existingPodcasts, err := database.GetPodcasts(ctx, 0, false)
	if err != nil {
		logger.Printf("Error fetching podcasts: %v", err)
		return err
	}

	for _, podcast := range existingPodcasts {
		go RefreshPodcast(podcast)
	}
	return nil
}

func RefreshPodcast(podcast types.PodcastChannel) error {
	ctx := context.Background()

	existingPodcastId, err := strconv.Atoi(podcast.Id)
	if err != nil {
		logger.Printf("Error converting existing podcast ID to int: %v", err)
		return fmt.Errorf("converting existing podcast ID to int: %v", err)
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(podcast.Url)
	if err != nil {
		logger.Printf("Error parsing feed URL: %v", err)
		return fmt.Errorf("parsing feed URL: %v", err)
	}

	var coverArt string
	if feed.Image.URL != podcast.OriginalImageUrl {
		coverArt, err = SavePodcastImage(ctx, feed.Image.URL)
		if err != nil {
			logger.Printf("Error saving podcast channel image: %v", err)
			return fmt.Errorf("saving podcast channel image: %v", err)
		}
	} else {
		coverArt = podcast.CoverArt
	}

	if err := database.UpdatePodcastChannel(
		ctx,
		existingPodcastId,
		podcast.Url,
		feed.Title,
		feed.Description,
		feed.Image.URL,
		coverArt,
		podcast.LastRefresh,
		feed.Categories,
	); err != nil {
		logger.Printf("Error updating podcast channel: %v", err)
		return fmt.Errorf("updating podcast channel: %v", err)
	}

	logger.Printf("Updated podcast channel '%s', fetching any new episodes...", feed.Title)

	go createPodcastEpisodesForFeed(ctx, feed, existingPodcastId)

	return nil
}
