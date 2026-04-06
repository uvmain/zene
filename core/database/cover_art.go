package database

import (
	"context"
	"fmt"
	"zene/core/logger"
)

func GetMediaCoverType(ctx context.Context, mediaId string) (string, error) {
	IsValidMetadataId, metadataType, err := IsValidMetadataId(ctx, mediaId)
	if err != nil && err.Error() != "checking metadata ID validity: sql: no rows in result set" {
		logger.Printf("error checking media id parameter: %v", err)
		return "", fmt.Errorf("error checking media id parameter: %v", err)
	}

	if IsValidMetadataId {
		switch metadataType {
		case MetadataAlbum:
			return "album", nil
		case MetadataArtist:
			return "artist", nil
		case MetadataTrack:
			return "track", nil
		case MetadataPodcastChannel:
			return "podcast_channel", nil
		case MetadataPodcastEpisode:
			return "podcast_episode", nil
		}
	}

	return "", fmt.Errorf("unknown metadata type: %s", metadataType)
}
