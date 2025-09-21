package podcasts

import (
	"context"
	"fmt"
	"path/filepath"
	"zene/core/art"
	"zene/core/config"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"

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

	for _, item := range feed.Items {
		logger.Printf("guid: %s title: %v", item.GUID, item.Title)
	}

	coverArt, err := SavePodcastChannelImage(ctx, feed.Image.URL)
	if err != nil {
		return fmt.Errorf("saving podcast channel image: %v", err)
	}

	err = database.CreatePodcastChannel(ctx, feedUrl, feed.Title, feed.Description, feed.Image.URL, coverArt, "")
	if err != nil {
		return fmt.Errorf("creating podcast channel: %v", err)
	}

	return nil
}

func SavePodcastChannelImage(ctx context.Context, imageUrl string) (string, error) {
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
