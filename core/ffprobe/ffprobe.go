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

func GetTagsFromFile(ctx context.Context, audiofilePath string) (types.FfprobeStandard, error) {
	cmd := exec.Command(config.FfprobePath, "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", audiofilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("Error running ffprobe: %s", output)
		return types.FfprobeStandard{}, err
	}

	if filepath.Ext(audiofilePath) == ".opus" {
		var ffprobeOutput types.FfprobeOpusOutput

		err = json.Unmarshal(output, &ffprobeOutput)
		if err != nil {
			logger.Printf("Error parsing ffprobe output: %v", err)
			return types.FfprobeStandard{}, err
		}

		return types.FfprobeStandard{
			Filename:   ffprobeOutput.Format.Filename,
			FormatName: ffprobeOutput.Format.FormatName,
			Tags:       ffprobeOutput.Streams[0].Tags,
			Duration:   ffprobeOutput.Format.Duration,
			Size:       ffprobeOutput.Format.Size,
			Bitrate:    ffprobeOutput.Format.Bitrate,
		}, nil
	} else {
		var ffprobeOutput types.FfprobeOutput

		err = json.Unmarshal(output, &ffprobeOutput)
		if err != nil {
			logger.Printf("Error parsing ffprobe output: %v", err)
			return types.FfprobeStandard{}, err
		}

		return types.FfprobeStandard{
			Filename:   ffprobeOutput.Format.Filename,
			FormatName: ffprobeOutput.Format.FormatName,
			Tags:       ffprobeOutput.Format.Tags,
			Duration:   ffprobeOutput.Format.Duration,
			Size:       ffprobeOutput.Format.Size,
			Bitrate:    ffprobeOutput.Format.Bitrate,
		}, nil
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

func ParseTags(ctx context.Context, ffprobeOutput types.FfprobeStandard) (types.Tags, error) {
	parsedArtist := getTagStringValue(ffprobeOutput.Tags, []string{"artist", "album_artist"})
	parsedAlbumArtist := getTagStringValue(ffprobeOutput.Tags, []string{"album_artist", "album-artist", "albumartist"})
	parsedTitle := getTagStringValue(ffprobeOutput.Tags, []string{"title"})
	parsedAlbum := getTagStringValue(ffprobeOutput.Tags, []string{"album"})
	parsedGenre := getTagStringValue(ffprobeOutput.Tags, []string{"genre"})
	parsedReleaseDate := getTagStringValue(ffprobeOutput.Tags, []string{"date", "release_date", "ORIGINAL_DATE"})
	musicBrainzAlbumId := getTagStringValue(ffprobeOutput.Tags, []string{"MUSICBRAINZ_ALBUMID", "MusicBrainz Album Id", "musicbrainz Album Id"})
	musicBrainzArtistId := getTagStringValue(ffprobeOutput.Tags, []string{"MUSICBRAINZ_ARTISTID", "MusicBrainz Artist Id", "musicbrainz Artist Id"})
	musicBrainzTrackId := getTagStringValue(ffprobeOutput.Tags, []string{"MUSICBRAINZ_TRACKID", "MusicBrainz Release Track Id", "musicbrainz Release Track Id"})
	totalTracks := getTagStringValue(ffprobeOutput.Tags, []string{"TOTALTRACKS"})
	trackNumber := getTagStringValue(ffprobeOutput.Tags, []string{"track"})
	totalDiscs := getTagStringValue(ffprobeOutput.Tags, []string{"TOTALDISCS"})
	discNumber := getTagStringValue(ffprobeOutput.Tags, []string{"disc"})
	label := getTagStringValue(ffprobeOutput.Tags, []string{"label", "publisher"})

	if musicBrainzArtistId == "" || musicBrainzAlbumId == "" || musicBrainzTrackId == "" {
		return types.Tags{}, fmt.Errorf("Unable to parse musicbrainz metadata from file")
	}

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

	if discNumber == "" {
		discNumber = "1"
	}

	if parsedReleaseDate == "" {
		musicBrainzData, err := musicbrainz.GetMetadataForMusicBrainzAlbumId(musicBrainzAlbumId)
		if err != nil {
			logger.Printf("Error fetching parsedReleaseDate from musicbrainz: %v", err)
			return types.Tags{}, err
		}
		parsedReleaseDate = musicBrainzData.Date
	}

	if trackNumber == "" {
		musicBrainzData, err := musicbrainz.GetMetadataForMusicBrainzAlbumId(musicBrainzAlbumId)
		if err != nil {
			logger.Printf("Error fetching parsedReleaseDate from musicbrainz: %v", err)
			return types.Tags{}, err
		}
		//if musicBrainzData.Media[0].tracks contains trackNumber, use that
		for _, media := range musicBrainzData.Media {
			for _, track := range media.Tracks {
				if track.Recording.ID == musicBrainzTrackId {
					trackNumber = track.Number
					break
				}
			}
			if trackNumber != "" {
				break
			}
		}
		// else check if musicBrainzData.Media[0].Pregap contains trackNumber
		if trackNumber == "" && len(musicBrainzData.Media) > 0 && musicBrainzData.Media[0].Pregap.Recording.ID == musicBrainzTrackId {
			trackNumber = musicBrainzData.Media[0].Pregap.Number
		}
	}

	parsedTags := types.Tags{
		Format:              ffprobeOutput.FormatName,
		Duration:            ffprobeOutput.Duration,
		Size:                ffprobeOutput.Size,
		Bitrate:             ffprobeOutput.Bitrate,
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
