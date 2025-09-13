package podcasts

import (
	"context"
	"fmt"
	"zene/core/database"

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

	err = database.CreatePodcastChannel(ctx, feedUrl, feed.Title, feed.Description, feed.Image.URL, "")
	if err != nil {
		return fmt.Errorf("creating podcast channel: %v", err)
	}

	return nil
}
