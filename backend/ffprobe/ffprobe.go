package ffprobe

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os/exec"
	"strings"

	"zene/config"
	"zene/musicbrainz"
	"zene/types"
)

func GetTags(audfilePath string) (types.TrackMetadata, error) {
	cmd := exec.Command(config.FfprobePath, "-show_format", "-print_format", "json", audfilePath)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error running ffprobe: %v", err)
		return types.TrackMetadata{}, err
	}

	var ffprobeOutput struct {
		Format struct {
			Filename   string            `json:"filename"`
			FormatName string            `json:"format_name"`
			Tags       map[string]string `json:"tags"`
			Duration   string            `json:"duration"`
			Size       string            `json:"size"`
			Bitrate    string            `json:"bit_rate"`
		} `json:"format"`
	}

	err = json.Unmarshal(output, &ffprobeOutput)
	if err != nil {
		log.Printf("Error parsing ffprobe output: %v", err)
		return types.TrackMetadata{}, err
	}

	tags := ffprobeOutput.Format.Tags

	parsedArtist := getTagStringValue(tags, []string{"artist", "album_artist"})
	parsedTitle := getTagStringValue(tags, []string{"title"})
	parsedAlbum := getTagStringValue(tags, []string{"album"})
	parsedGenre := getTagStringValue(tags, []string{"genre"})
	parsedReleaseDate := getTagStringValue(tags, []string{"date", "release_date"})
	musicBrainzAlbumId := getTagStringValue(tags, []string{"MUSICBRAINZ_ALBUMID", "MusicBrainz Album Id"})
	musicBrainzArtistId := getTagStringValue(tags, []string{"MUSICBRAINZ_ARTISTID", "MusicBrainz Artist Id"})
	musicBrainzTrackId := getTagStringValue(tags, []string{"MUSICBRAINZ_TRACKID", "MusicBrainz Release Track Id"})
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
		metadata, _ := musicbrainz.GetMetadataForMusicBrainzAlbumId(musicBrainzAlbumId)
		parsedReleaseDate = metadata.Date
	}

	metadata := types.TrackMetadata{
		Filename:            ffprobeOutput.Format.Filename,
		Format:              ffprobeOutput.Format.FormatName,
		Duration:            ffprobeOutput.Format.Duration,
		Size:                ffprobeOutput.Format.Size,
		Bitrate:             ffprobeOutput.Format.Bitrate,
		Title:               parsedTitle,
		Artist:              parsedArtist,
		Album:               parsedAlbum,
		AlbumArtist:         ffprobeOutput.Format.Tags["album_artist"],
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

	return metadata, nil
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
