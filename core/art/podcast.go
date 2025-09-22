package art

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"zene/core/config"
	"zene/core/logic"
)

func GetArtForPodcast(ctx context.Context, coverArtId string, size int) ([]byte, time.Time, error) {
	// prevent path traversal
	if strings.Contains(coverArtId, "/") || strings.Contains(coverArtId, "\\") || strings.Contains(coverArtId, "..") {
		return nil, time.Now(), fmt.Errorf("invalid podcast coverArtId")
	}

	file_name := fmt.Sprintf("%s.jpg", coverArtId)
	filePath, _ := filepath.Abs(filepath.Join(config.PodcastArtFolder, file_name))

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, time.Now(), fmt.Errorf("podcast coverArt file does not exist: %s:  %s", filePath, err)
	}

	modTime := info.ModTime()

	blob, err := logic.ResizeJpegImage(ctx, filePath, size, 90)
	if err != nil {
		return nil, time.Now(), fmt.Errorf("error reading image for filepath %s: %s", filePath, err)
	}
	return blob, modTime, nil
}
