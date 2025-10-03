package ffprobe

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"zene/core/config"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/musicbrainz"
	"zene/core/types"
)

func InitializeFfprobe(ctx context.Context) {
	logger.Printf("FFPROBE_PATH: %s", config.FfprobePath)

	if io.FileExists(config.FfprobePath) {
		logger.Printf("ffprobe binary found at %s", config.FfprobePath)
	} else {
		err := DownloadFfprobeBinary()
		if err != nil {
			log.Fatalf("failed to download ffprobe binary: %v", err)
		}
	}

	version, err := exec.CommandContext(ctx, config.FfprobePath, "-version").Output()
	if err != nil {
		log.Fatalf("ffprobe not found at %s: %v", config.FfprobePath, err)
	} else {
		logger.Printf("ffprobe version is %v", strings.Split(string(version), "\n")[0])
	}
}

func GetMetadata(ctx context.Context, audiofilePath string) (types.FileMetadata, error) {
	fileTags, err := GetMetadataFromFile(ctx, audiofilePath)
	if err != nil {
		return types.FileMetadata{}, err
	}
	parsedTags, err := ParseMetadata(ctx, fileTags)
	if err != nil {
		return types.FileMetadata{}, err
	}
	return parsedTags, nil
}

func GetDurationAndBitrate(ctx context.Context, audiofilePath string) (string, string, error) {
	cmd := exec.CommandContext(ctx, config.FfprobePath,
		"-v", "quiet",
		"-show_entries", "format=duration,bit_rate",
		"-of", "json",
		audiofilePath,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("Error running ffprobe: %s", output)
		return "", "", err
	}

	var ffprobeOutput struct {
		Format struct {
			Duration string `json:"duration"`
			BitRate  string `json:"bit_rate"`
		} `json:"format"`
	}

	if err := json.Unmarshal(output, &ffprobeOutput); err != nil {
		return "", "", err
	}

	if ffprobeOutput.Format.Duration == "" || ffprobeOutput.Format.BitRate == "" {
		return "0", "0", fmt.Errorf("ffprobe did not return duration/bitrate")
	}

	return ffprobeOutput.Format.Duration, ffprobeOutput.Format.BitRate, nil
}

func GetMetadataFromFile(ctx context.Context, audiofilePath string) (types.FfprobeStandard, error) {
	cmd := exec.CommandContext(ctx, config.FfprobePath, "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", audiofilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("Error running ffprobe: %s", output)
		return types.FfprobeStandard{}, err
	}

	if filepath.Ext(audiofilePath) == ".opus" {
		var ffprobeOpusOutput types.FfprobeOpusOutput

		err = json.Unmarshal(output, &ffprobeOpusOutput)
		if err != nil {
			logger.Printf("Error parsing ffprobe output: %v", err)
			return types.FfprobeStandard{}, err
		}

		sampleRateInt, err := strconv.Atoi(ffprobeOpusOutput.Streams[0].SampleRate)
		if err != nil {
			sampleRateInt = 0
		}

		return types.FfprobeStandard{
			Filename:   ffprobeOpusOutput.Format.Filename,
			FormatName: ffprobeOpusOutput.Format.FormatName,
			Tags:       ffprobeOpusOutput.Streams[0].Tags,
			Duration:   ffprobeOpusOutput.Format.Duration,
			Size:       ffprobeOpusOutput.Format.Size,
			Bitrate:    ffprobeOpusOutput.Format.Bitrate,
			BitDepth:   ffprobeOpusOutput.Streams[0].BitDepth,
			SampleRate: sampleRateInt,
			Channels:   ffprobeOpusOutput.Streams[0].Channels,
			Codec:      ffprobeOpusOutput.Streams[0].Codec,
		}, nil
	} else {
		var ffprobeStandardOutput types.FfprobeStandardOutput

		err = json.Unmarshal(output, &ffprobeStandardOutput)
		if err != nil {
			logger.Printf("Error parsing ffprobe output: %v", err)
			return types.FfprobeStandard{}, err
		}

		sampleRateInt, err := strconv.Atoi(ffprobeStandardOutput.Streams[0].SampleRate)
		if err != nil {
			sampleRateInt = 0
		}

		return types.FfprobeStandard{
			Filename:   ffprobeStandardOutput.Format.Filename,
			FormatName: ffprobeStandardOutput.Format.FormatName,
			Tags:       ffprobeStandardOutput.Format.Tags,
			Duration:   ffprobeStandardOutput.Format.Duration,
			Size:       ffprobeStandardOutput.Format.Size,
			Bitrate:    ffprobeStandardOutput.Format.Bitrate,
			BitDepth:   ffprobeStandardOutput.Streams[0].BitDepth,
			SampleRate: sampleRateInt,
			Channels:   ffprobeStandardOutput.Streams[0].Channels,
			Codec:      ffprobeStandardOutput.Streams[0].Codec,
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

func ParseMetadata(ctx context.Context, ffprobeOutput types.FfprobeStandard) (types.FileMetadata, error) {
	parsedArtist := getTagStringValue(ffprobeOutput.Tags, []string{"artist", "album_artist"})
	parsedAlbumArtist := getTagStringValue(ffprobeOutput.Tags, []string{"album_artist", "album-artist", "albumartist"})
	parsedTitle := getTagStringValue(ffprobeOutput.Tags, []string{"title"})
	parsedAlbum := getTagStringValue(ffprobeOutput.Tags, []string{"album"})
	parsedGenre := getTagStringValue(ffprobeOutput.Tags, []string{"genre"})
	parsedReleaseDate := getTagStringValue(ffprobeOutput.Tags, []string{"date", "release_date", "ORIGINAL_DATE", "ORIGINALDATE"})
	musicBrainzAlbumId := getTagStringValue(ffprobeOutput.Tags, []string{"MUSICBRAINZ_ALBUMID", "MusicBrainz Album Id", "musicbrainz Album Id"})
	musicBrainzArtistId := getTagStringValue(ffprobeOutput.Tags, []string{"MUSICBRAINZ_ARTISTID", "MusicBrainz Artist Id", "musicbrainz Artist Id"})
	musicBrainzTrackId := getTagStringValue(ffprobeOutput.Tags, []string{"MUSICBRAINZ_TRACKID", "MusicBrainz Release Track Id", "musicbrainz Release Track Id"})
	totalTracks := getTagStringValue(ffprobeOutput.Tags, []string{"TOTALTRACKS"})
	trackNumber := getTagStringValue(ffprobeOutput.Tags, []string{"track"})
	originalYear := getTagStringValue(ffprobeOutput.Tags, []string{"TORY", "ORY", "ORIGINAL_YEAR", "ORIGINAL YEAR"})
	totalDiscs := getTagStringValue(ffprobeOutput.Tags, []string{"TOTALDISCS"})
	discNumber := getTagStringValue(ffprobeOutput.Tags, []string{"disc"})
	label := getTagStringValue(ffprobeOutput.Tags, []string{"label", "publisher"})

	if musicBrainzArtistId == "" || musicBrainzAlbumId == "" || musicBrainzTrackId == "" {
		return types.FileMetadata{}, fmt.Errorf("Unable to parse musicbrainz metadata from file")
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

	if parsedReleaseDate == "" && originalYear != "" {
		parsedReleaseDate = fmt.Sprintf("%s-01-01", originalYear)
	}

	var musicBrainzData types.MbRelease

	if parsedReleaseDate == "" || trackNumber == "" {
		musicBrainzData, err := musicbrainz.GetMetadataForMusicBrainzAlbumId(musicBrainzAlbumId)
		if err != nil {
			logger.Printf("Error fetching parsedReleaseDate from musicbrainz: %v", err)
			return types.FileMetadata{}, err
		}
		parsedReleaseDate = musicBrainzData.Date
	}

	if trackNumber == "" {
		// if musicBrainzData.Media[0].tracks contains trackNumber, use that
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

	trackNumberInt, err := strconv.Atoi(trackNumber)
	if err != nil {
		logger.Printf("Error converting track number to int: %v", err)
		trackNumberInt = 0
	}

	discNumberInt, err := strconv.Atoi(discNumber)
	if err != nil {
		logger.Printf("Error converting disc number to int: %v", err)
		discNumberInt = 0
	}

	totalTracksInt, err := strconv.Atoi(totalTracks)
	if err != nil {
		logger.Printf("Error converting total tracks to int: %v", err)
		totalTracksInt = 0
	}

	totalDiscsInt, err := strconv.Atoi(totalDiscs)
	if err != nil {
		logger.Printf("Error converting total discs to int: %v", err)
		totalDiscsInt = 0
	}

	parsedMetadata := types.FileMetadata{
		Format:              ffprobeOutput.FormatName,
		Duration:            ffprobeOutput.Duration,
		Size:                ffprobeOutput.Size,
		Bitrate:             ffprobeOutput.Bitrate,
		Title:               parsedTitle,
		Artist:              parsedArtist,
		Album:               parsedAlbum,
		AlbumArtist:         parsedAlbumArtist,
		Genre:               parsedGenre,
		TrackNumber:         trackNumberInt,
		DiscNumber:          discNumberInt,
		ReleaseDate:         parsedReleaseDate,
		MusicBrainzArtistID: musicBrainzArtistId,
		MusicBrainzAlbumID:  musicBrainzAlbumId,
		MusicBrainzTrackID:  musicBrainzTrackId,
		Label:               label,
		TotalTracks:         totalTracksInt,
		TotalDiscs:          totalDiscsInt,
		Codec:               ffprobeOutput.Codec,
		BitDepth:            ffprobeOutput.BitDepth,
		SampleRate:          ffprobeOutput.SampleRate,
		Channels:            ffprobeOutput.Channels,
	}

	return parsedMetadata, nil
}
