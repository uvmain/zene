package podcasts

import (
	"cmp"
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

	if len(newEpisodes) > 0 {
		logger.Printf("Found %d new episodes for podcast ID %d, %s", len(newEpisodes), podcastId, feed.Title)
	}

	podcastEpisodes := make([]types.PodcastEpisodeRow, 0, len(newEpisodes))
	for _, item := range newEpisodes {

		if len(item.Enclosures) == 0 {
			logger.Printf("Skipping episode '%s' (GUID: %s) for podcast ID %d: no enclosures found", item.Title, item.GUID, podcastId)
			continue
		}

		authors := ""
		for _, author := range item.Authors {
			authors += author.Name + ", "
		}
		if len(authors) > 2 {
			authors = authors[:len(authors)-2] // remove trailing comma and space
		}

		var imageUrl string
		var coverArt string
		if item.Image != nil {
			imageUrl = cmp.Or(item.Image.URL, item.ITunesExt.Image)
			coverArt, err = SavePodcastImage(ctx, imageUrl)
			if err != nil {
				logger.Printf("Error saving podcast episode cover art: %v", err)
				return fmt.Errorf("saving podcast episode cover art: %v", err)
			}
		} else {
			podcastChannel, err := database.GetPodcastsUserless(ctx, strconv.Itoa(podcastId))
			if err != nil {
				logger.Printf("Error getting podcast channels: %v", err)
				return fmt.Errorf("getting podcast channels: %v", err)
			}
			coverArt = podcastChannel[0].CoverArt
		}

		episodeLink := item.Enclosures[0]

		durationString := item.ITunesExt.Duration
		episodeDuration, err := strconv.Atoi(durationString)
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
			SourceUrl: episodeLink.URL,
		}
		podcastEpisodes = append(podcastEpisodes, episode)
	}

	if err := database.InsertPodcastEpisodes(podcastEpisodes); err != nil {
		logger.Printf("Error inserting podcast episodes: %v", err)
		return fmt.Errorf("inserting podcast episodes: %v", err)
	}

	database.UpdatePodcastChannelLastRefresh(podcastId)

	if len(podcastEpisodes) > 0 {
		logger.Printf("Inserted %d new episodes for podcast ID %d, %s", len(podcastEpisodes), podcastId, feed.Title)
	}

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
	existingPodcasts, err := database.GetPodcastsUserless(ctx, "")
	if err != nil {
		logger.Printf("Error fetching podcasts: %v", err)
		return err
	}

	for _, podcast := range existingPodcasts {
		logger.Printf("Refreshing podcast episodes for %s", podcast.Title)
		go RefreshPodcast(podcast)
	}
	return nil
}

func RefreshPodcastById(ctx context.Context, id string) error {
	existingPodcastId, err := strconv.Atoi(id)
	if err != nil {
		logger.Printf("Error converting existing podcast ID to int: %v", err)
		return fmt.Errorf("converting existing podcast ID to int: %v", err)
	}

	existingPodcast, err := database.GetPodcastChannelById(ctx, existingPodcastId)
	if err != nil {
		logger.Printf("Error fetching podcast by ID: %v", err)
		return fmt.Errorf("fetching podcast by ID: %v", err)
	}

	go RefreshPodcast(existingPodcast)
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
