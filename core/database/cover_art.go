package database

import (
	"context"
	"fmt"
	"zene/core/logger"
)

func GetMediaCoverType(ctx context.Context, mediaId string) (string, error) {
	IsValidMetadataId, metadataStruct, err := IsValidMetadataId(ctx, mediaId)
	if err != nil && err.Error() != "checking metadata ID validity: sql: no rows in result set" {
		logger.Printf("error checking media id parameter: %v", err)
		return "", fmt.Errorf("error checking media id parameter: %v", err)
	}

	if IsValidMetadataId {
		if metadataStruct.MusicbrainzAlbumId {
			return "album", nil
		} else if metadataStruct.MusicbrainzArtistId {
			return "artist", nil
		} else if metadataStruct.MusicbrainzTrackId {
			return "track", nil
		}
	}

	isValidPodcastGuid, err := IsValidPodcastCover(ctx, mediaId)
	if err != nil && err.Error() != "checking podcast cover validity: sql: no rows in result set" {
		logger.Printf("error checking podcast cover validity: %v", err)
		return "", fmt.Errorf("error checking podcast cover validity: %v", err)
	}

	if isValidPodcastGuid {
		return "podcast", nil
	}

	return "", fmt.Errorf("unknown cover art type")
}
