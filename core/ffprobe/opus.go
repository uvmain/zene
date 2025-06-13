package ffprobe

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"zene/core/config"
	"zene/core/musicbrainz"
	"zene/core/types"
)

func GetOpusTags(audiofilePath string) (types.Tags, error) {
	cmd := exec.Command(config.FfprobePath, "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", audiofilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error running ffprobe: %s", output)
		return types.Tags{}, err
	}

	var ffprobeOutput types.FfprobeOpusOutput

	err = json.Unmarshal(output, &ffprobeOutput)
	if err != nil {
		log.Printf("Error parsing ffprobe output: %v", err)
		return types.Tags{}, err
	}

	tags := ffprobeOutput.Streams[0].Tags

	parsedArtist := getTagStringValue(tags, []string{"artist"})
	parsedAlbumArtist := getTagStringValue(tags, []string{"album_artist", "album-artist", "albumartist"})
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
		musicBrainzData, _ := musicbrainz.GetMetadataForMusicBrainzAlbumId(musicBrainzAlbumId)
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
