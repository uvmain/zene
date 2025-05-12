package ffprobe

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os/exec"
	"strings"

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

	var parsedArtist = ffprobeOutput.Format.Tags["ARTIST"]
	if parsedArtist == "" {
		parsedArtist = ffprobeOutput.Format.Tags["artist"]
	}
	if parsedArtist == "" {
		parsedArtist = ffprobeOutput.Format.Tags["album_artist"]
	}

	var parsedTitle = ffprobeOutput.Format.Tags["TITLE"]
	if parsedTitle == "" {
		parsedTitle = ffprobeOutput.Format.Tags["title"]
	}

	var parsedAlbum = ffprobeOutput.Format.Tags["ALBUM"]
	if parsedAlbum == "" {
		parsedAlbum = ffprobeOutput.Format.Tags["album"]
	}

	var parsedGenre = ffprobeOutput.Format.Tags["GENRE"]
	if parsedGenre == "" {
		parsedGenre = ffprobeOutput.Format.Tags["genre"]
	}

	var parsedReleaseDate = ffprobeOutput.Format.Tags["DATE"]
	if parsedReleaseDate == "" {
		parsedReleaseDate = ffprobeOutput.Format.Tags["date"]
	}

	var musicBrainzAlbumId = ffprobeOutput.Format.Tags["MUSICBRAINZ_ALBUMID"]
	if musicBrainzAlbumId == "" {
		musicBrainzAlbumId = ffprobeOutput.Format.Tags["MusicBrainz Album Id"]
	}

	var musicBrainzArtistId = ffprobeOutput.Format.Tags["MUSICBRAINZ_ARTISTID"]
	if musicBrainzArtistId == "" {
		musicBrainzArtistId = ffprobeOutput.Format.Tags["MusicBrainz Artist Id"]
	}

	var musicBrainzTrackId = ffprobeOutput.Format.Tags["MUSICBRAINZ_TRACKID"]
	if musicBrainzTrackId == "" {
		musicBrainzTrackId = ffprobeOutput.Format.Tags["MusicBrainz Release Track Id"]
	}
	if musicBrainzTrackId == "" {
		var filenameBytes = []byte(ffprobeOutput.Format.Filename)
		musicBrainzTrackId = base64.StdEncoding.EncodeToString(filenameBytes)
	}

	var totalTracks = ffprobeOutput.Format.Tags["TOTALTRACKS"]
	if totalTracks == "" {
		totalTracks = ffprobeOutput.Format.Tags["totaltracks"]
	}

	var trackNumber = ffprobeOutput.Format.Tags["track"]
	if strings.Contains(trackNumber, "/") {
		splitValue := strings.Split(trackNumber, "/")
		trackNumber = splitValue[0]
		if totalTracks == "" && len(splitValue) > 1 {
			totalTracks = splitValue[1]
		}
	}

	var totalDiscs = ffprobeOutput.Format.Tags["TOTALDISCS"]
	if totalDiscs == "" {
		totalDiscs = ffprobeOutput.Format.Tags["totaldiscs"]
	}

	var discNumber = ffprobeOutput.Format.Tags["disc"]
	if strings.Contains(discNumber, "/") {
		splitValue := strings.Split(discNumber, "/")
		discNumber = splitValue[0]
		if totalDiscs == "" && len(splitValue) > 1 {
			totalDiscs = splitValue[1]
		}
	}

	var label = ffprobeOutput.Format.Tags["LABEL"]
	if label == "" {
		label = ffprobeOutput.Format.Tags["label"]
	}
	if label == "" {
		label = ffprobeOutput.Format.Tags["publisher"]
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
