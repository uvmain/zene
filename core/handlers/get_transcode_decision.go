package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

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

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "User not found", "")
		return
	}

	mediaFilepath, err := database.GetMediaFilePath(ctx, mediaId)

	if mediaFilepath == "" || err != nil {
		// check if the file is a podcast episode
		if requestUser.PodcastRole {
			episode, _ := database.GetPodcastEpisodeByGuid(ctx, mediaId)
			if episode.SourceUrl != "" {
				http.Redirect(w, r, episode.SourceUrl, http.StatusFound)
				return
			}
		}
	}

	if err != nil {
		logger.Printf("Error querying database for media filepath %s: %v", mediaId, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "File not found in database.", "")
		return
	}

	if mediaFilepath == "" {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "File not available to stream.", "")
		return
	}

	metadata, err := ffprobe.GetMetadataFromFile(ctx, mediaFilepath)
	if err != nil {
		logger.Printf("Error probing media file %s: %v", mediaFilepath, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Unable to inspect media file.", "")
		return
	}

	bitrateInt, err := strconv.Atoi(strings.TrimSpace(metadata.Bitrate))
	if err != nil || bitrateInt < 0 {
		bitrateInt = 0
	}

	sourceStream := types.StreamDetails{
		Protocol:        "http",
		Container:       logic.ParseMediaContainer(mediaFilepath, metadata),
		Codec:           strings.ToLower(strings.TrimSpace(metadata.Codec)),
		AudioChannels:   metadata.Channels,
		AudioBitrate:    bitrateInt,
		AudioSamplerate: metadata.SampleRate,
		AudioBitdepth:   metadata.BitDepth,
	}

	decision := logic.BuildTranscodeDecision(mediaId, mediaType, clientInfo, sourceStream)

	response := subsonic.GetPopulatedSubsonicResponse(ctx)
	response.SubsonicResponse.TranscodeDecision = &decision

	net.WriteSubsonicResponse(w, r, response, format)
}
