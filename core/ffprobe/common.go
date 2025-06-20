package ffprobe

import (
	"encoding/base64"
	"encoding/json"
	"os/exec"
	"strings"
	"zene/core/config"
	"zene/core/logger"
	"zene/core/musicbrainz"
	"zene/core/types"
)

func GetCommonTags(audiofilePath string) (types.Tags, error) {
	cmd := exec.Command(config.FfprobePath, "-v", "quiet", "-show_format", "-show_streams", "-print_format", "json", audiofilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("Error running ffprobe: %s", output)
		return types.Tags{}, err
	}

	var ffprobeOutput types.FfprobeOutput

	err = json.Unmarshal(output, &ffprobeOutput)
	if err != nil {
		logger.Printf("Error parsing ffprobe output: %v", err)
		return types.Tags{}, err
	}

	tags := ffprobeOutput.Format.Tags

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

	if musicBrainzTrackId == "" {
		var filenameBytes = []byte(ffprobeOutput.Format.Filename)
		musicBrainzTrackId = base64.StdEncoding.EncodeToString(filenameBytes)
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

	if parsedReleaseDate == "" {
		musicBrainzData, err := musicbrainz.GetMetadataForMusicBrainzAlbumId(musicBrainzAlbumId)
		if err != nil {
			logger.Printf("Error fetching parsedReleaseDate from musicbrainz: %v", err)
			return types.Tags{}, err
		}
		parsedReleaseDate = musicBrainzData.Date
	}

	parsedTags := types.Tags{
		Format:              ffprobeOutput.Format.FormatName,
		Duration:            ffprobeOutput.Format.Duration,
		Size:                ffprobeOutput.Format.Size,
		Bitrate:             ffprobeOutput.Format.Bitrate,
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
