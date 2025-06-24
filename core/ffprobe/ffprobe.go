package ffprobe

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"zene/core/config"
	"zene/core/logger"
	"zene/core/musicbrainz"
	"zene/core/types"
)

func GetTags(ctx context.Context, audiofilePath string) (types.Tags, error) {
	fileTags, err := GetTagsFromFile(ctx, audiofilePath)
	if err != nil {
		return types.Tags{}, err
	}
	parsedTags, err := ParseTags(ctx, fileTags)
	if err != nil {
		return types.Tags{}, err
	}
	return parsedTags, nil
}

func GetTagsFromFile(ctx context.Context, audiofilePath string) (map[string]string, error) {
	cmd := exec.Command(config.FfprobePath, "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", audiofilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("Error running ffprobe: %s", output)
		return map[string]string{}, err
	}

	if filepath.Ext(audiofilePath) == ".opus" {
		var ffprobeOutput types.FfprobeOpusOutput

		err = json.Unmarshal(output, &ffprobeOutput)
		if err != nil {
			logger.Printf("Error parsing ffprobe output: %v", err)
			return map[string]string{}, err
		}

		return ffprobeOutput.Streams[0].Tags, nil
	} else {
		var ffprobeOutput types.FfprobeOutput

		err = json.Unmarshal(output, &ffprobeOutput)
		if err != nil {
			logger.Printf("Error parsing ffprobe output: %v", err)
			return map[string]string{}, err
		}

		return ffprobeOutput.Format.Tags, nil
	}
}

func getTagStringValue(tags map[string]string, inputs []string) string {
	for _, input := range inputs {
		value := tags[input]
		if value != "" {
			return value
		}
		value = tags[strings.ToUpper(input)]
		if value != "" {
			return value
		}
		value = tags[strings.ToLower(input)]
		if value != "" {
			return value
		}
	}
	return ""
}

func ParseTags(ctx context.Context, tags map[string]string) (types.Tags, error) {
	parsedFormat := getTagStringValue(tags, []string{"format"})
	parsedDuration := getTagStringValue(tags, []string{"duration"})
	parsedSize := getTagStringValue(tags, []string{"size"})
	parsedBitrate := getTagStringValue(tags, []string{"bit_rate", "bitrate"})
	parsedArtist := getTagStringValue(tags, []string{"artist", "album_artist"})
	parsedAlbumArtist := getTagStringValue(tags, []string{"album_artist", "album-artist", "albumartist"})
	parsedTitle := getTagStringValue(tags, []string{"title"})
	parsedAlbum := getTagStringValue(tags, []string{"album"})
	parsedGenre := getTagStringValue(tags, []string{"genre"})
	parsedReleaseDate := getTagStringValue(tags, []string{"date", "release_date", "ORIGINAL_DATE"})
	musicBrainzAlbumId := getTagStringValue(tags, []string{"MUSICBRAINZ_ALBUMID", "MusicBrainz Album Id", "musicbrainz Album Id"})
	musicBrainzArtistId := getTagStringValue(tags, []string{"MUSICBRAINZ_ARTISTID", "MusicBrainz Artist Id", "musicbrainz Artist Id"})
	musicBrainzTrackId := getTagStringValue(tags, []string{"MUSICBRAINZ_TRACKID", "MusicBrainz Release Track Id", "musicbrainz Release Track Id"})
	totalTracks := getTagStringValue(tags, []string{"TOTALTRACKS"})
	trackNumber := getTagStringValue(tags, []string{"track"})
	totalDiscs := getTagStringValue(tags, []string{"TOTALDISCS"})
	discNumber := getTagStringValue(tags, []string{"disc"})
	label := getTagStringValue(tags, []string{"label", "publisher"})

	if strings.Contains(trackNumber, "/") {
		splitValue := strings.Split(trackNumber, "/")
		trackNumber = splitValue[0]
		if totalTracks == "" && len(splitValue) > 1 {
			totalTracks = splitValue[1]
		}
	}

	if strings.Contains(discNumber, "/") {
		splitValue := strings.Split(discNumber, "/")
		discNumber = splitValue[0]
		if totalDiscs == "" && len(splitValue) > 1 {
			totalDiscs = splitValue[1]
		}
	}

	if parsedReleaseDate == "" {
		musicBrainzData, err := musicbrainz.GetMetadataForMusicBrainzAlbumId(musicBrainzAlbumId)
		if err != nil {
			logger.Printf("Error fetching parsedReleaseDate from musicbrainz: %v", err)
			return types.Tags{}, err
		}
		parsedReleaseDate = musicBrainzData.Date
	}

	if musicBrainzArtistId == "" || musicBrainzAlbumId == "" || musicBrainzTrackId == "" {
		return types.Tags{}, fmt.Errorf("Unable to parse musicbrainz metadata from file")
	}

	parsedTags := types.Tags{
		Format:              parsedFormat,
		Duration:            parsedDuration,
		Size:                parsedSize,
		Bitrate:             parsedBitrate,
		Title:               parsedTitle,
		Artist:              parsedArtist,
		Album:               parsedAlbum,
		AlbumArtist:         parsedAlbumArtist,
		Genre:               parsedGenre,
		TrackNumber:         trackNumber,
		DiscNumber:          discNumber,
		ReleaseDate:         parsedReleaseDate,
		MusicBrainzArtistID: musicBrainzArtistId,
		MusicBrainzAlbumID:  musicBrainzAlbumId,
		MusicBrainzTrackID:  musicBrainzTrackId,
		Label:               label,
		TotalTracks:         totalTracks,
		TotalDiscs:          totalDiscs,
	}

	return parsedTags, nil
}
