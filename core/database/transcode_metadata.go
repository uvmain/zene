package database

import (
	"context"
	"database/sql"
	"fmt"
	"zene/core/types"
)

func GetTranscodeMetadataForSong(ctx context.Context, mediaId string) (types.TranscodeMetadata, error) {
	metadata := types.TranscodeMetadata{}

	query := "select file_path, format, codec, channels, bitrate, sample_rate, bit_depth from metadata WHERE musicbrainz_track_id = ? LIMIT 1;"
	err := DB.QueryRow(query, mediaId).Scan(&metadata.FilePath, &metadata.Container, &metadata.Codec, &metadata.AudioChannels, &metadata.AudioBitrate, &metadata.AudioSamplerate, &metadata.AudioBitdepth)
	if err == sql.ErrNoRows {
		return metadata, fmt.Errorf("no track found for media ID '%s'", mediaId)
	} else if err != nil {
		return metadata, fmt.Errorf("error querying track ID: %v", err)
	}
	return metadata, nil
}
