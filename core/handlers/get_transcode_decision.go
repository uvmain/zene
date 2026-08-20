package handlers

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"zene/core/config"
	"zene/core/database"
	"zene/core/ffprobe"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetTranscodeDecision(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	mediaId := form["mediaid"]
	mediaType := form["mediatype"]

	if mediaId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "mediaId parameter is required", "")
		return
	}

	if mediaType != "song" && mediaType != "podcast" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "invalid mediaType parameter - must be 'song' or 'podcast'", "")
		return
	}

	var clientInfo types.ClientInfo

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&clientInfo); err != nil {
		if err == io.EOF {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "request body is required", "")
			return
		}
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "invalid request body", "")
		return
	}

	ctx := r.Context()

	metadata, err := GetMediaTranscodeMetadata(ctx, mediaId, mediaType)
	if err != nil {
		logger.Printf("Error fetching media metadata for %s: %v", mediaId, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Unable to inspect media file.", "")
		return
	}

	sourceStream := types.StreamDetails{
		Protocol:        "http",
		Container:       metadata.Container,
		Codec:           strings.ToLower(strings.TrimSpace(metadata.Codec)),
		AudioChannels:   metadata.AudioChannels,
		AudioBitrate:    metadata.AudioBitrate,
		AudioSamplerate: metadata.AudioSamplerate,
		AudioBitdepth:   metadata.AudioBitdepth,
	}

	decision, transcodeParams := logic.BuildTranscodeDecision(mediaId, mediaType, clientInfo, sourceStream)

	err = database.UpsertTranscodeParams(ctx, types.TranscodeParamsRow{
		ParamString:   decision.TranscodeParams,
		MediaId:       mediaId,
		MediaType:     mediaType,
		Container:     transcodeParams.Container,
		AudioCodec:    transcodeParams.AudioCodec,
		Protocol:      transcodeParams.Protocol,
		TargetFormat:  transcodeParams.TargetFormat,
		Bitrate:       transcodeParams.Bitrate,
		AudioChannels: transcodeParams.AudioChannels,
	})
	if err != nil {
		logger.Printf("Error upserting transcode params for %s: %v", mediaId, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "server failed to upsert transcode parameters", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)
	response.SubsonicResponse.TranscodeDecision = &decision

	net.WriteSubsonicResponse(w, r, response, format)
}

func GetMediaTranscodeMetadata(ctx context.Context, mediaId, mediaType string) (types.TranscodeMetadata, error) {
	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		return types.TranscodeMetadata{}, err
	}

	if !requestUser.PodcastRole && mediaType == "podcast" {
		return types.TranscodeMetadata{}, fmt.Errorf("user %s does not have permission to access podcasts", requestUser.Username)
	}

	metadata := types.TranscodeMetadata{}

	if mediaType == "podcast" {
		podcastMetadata, err := ffprobe.GetMetadataFromFile(ctx, mediaId)
		if err != nil {
			return types.TranscodeMetadata{}, fmt.Errorf("Error fetching podcast metadata for %s: %v", mediaId, err)
		}
		bitrateInt, err := strconv.Atoi(podcastMetadata.Bitrate)
		if err != nil {
			bitrateInt = config.DefaultBitRate
		}
		metadata.FilePath = mediaId
		metadata.Container = podcastMetadata.FormatName
		metadata.Codec = podcastMetadata.Codec
		metadata.AudioChannels = podcastMetadata.Channels
		metadata.AudioBitrate = cmp.Or(bitrateInt, config.DefaultBitRate)
		metadata.AudioSamplerate = podcastMetadata.SampleRate
		metadata.AudioBitdepth = podcastMetadata.BitDepth
	} else if mediaType == "song" {
		songMetadata, err := database.GetTranscodeMetadataForSong(ctx, mediaId)
		if err != nil {
			return types.TranscodeMetadata{}, fmt.Errorf("Error fetching song metadata for %s: %v", mediaId, err)
		}
		metadata.FilePath = songMetadata.FilePath
		metadata.Container = songMetadata.Container
		metadata.Codec = songMetadata.Codec
		metadata.AudioChannels = songMetadata.AudioChannels
		metadata.AudioBitrate = cmp.Or(songMetadata.AudioBitrate, config.DefaultBitRate)
		metadata.AudioSamplerate = songMetadata.AudioSamplerate
		metadata.AudioBitdepth = songMetadata.AudioBitdepth
	}

	return metadata, nil
}
