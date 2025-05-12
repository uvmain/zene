package ffprobe

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os/exec"

	"zene/config"
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

	var trackID = ffprobeOutput.Format.Tags["MUSICBRAINZ_TRACKID"]
	if trackID == "" {
		var filenameBytes = []byte(ffprobeOutput.Format.Filename)
		trackID = base64.StdEncoding.EncodeToString(filenameBytes)
	}

	metadata := types.TrackMetadata{
		MusicBrainzTrackID:  trackID,
		Filename:            ffprobeOutput.Format.Filename,
		Format:              ffprobeOutput.Format.FormatName,
		Duration:            ffprobeOutput.Format.Duration,
		Size:                ffprobeOutput.Format.Size,
		Bitrate:             ffprobeOutput.Format.Bitrate,
		Title:               ffprobeOutput.Format.Tags["TITLE"],
		Artist:              ffprobeOutput.Format.Tags["ARTIST"],
		Album:               ffprobeOutput.Format.Tags["ALBUM"],
		AlbumArtist:         ffprobeOutput.Format.Tags["album_artist"],
		Genre:               ffprobeOutput.Format.Tags["GENRE"],
		TrackNumber:         ffprobeOutput.Format.Tags["track"],
		DiscNumber:          ffprobeOutput.Format.Tags["disc"],
		ReleaseDate:         ffprobeOutput.Format.Tags["DATE"],
		MusicBrainzArtistID: ffprobeOutput.Format.Tags["MUSICBRAINZ_ARTISTID"],
		MusicBrainzAlbumID:  ffprobeOutput.Format.Tags["MUSICBRAINZ_ALBUMID"],
		Label:               ffprobeOutput.Format.Tags["LABEL"],
		TotalTracks:         ffprobeOutput.Format.Tags["TOTALTRACKS"],
		TotalDiscs:          ffprobeOutput.Format.Tags["TOTALDISCS"],
	}

	return metadata, nil
}
