package handlers

import (
	"net/http"
	"strconv"

	"zene/core/database"
	"zene/core/ffmpeg"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetTranscodeStream(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	mediaId := form["mediaid"]
	mediaType := form["mediatype"]
	offset := form["offset"]
	transcodeParams := form["transcodeparams"]

	if mediaId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "mediaId parameter is required", "")
		return
	}

	if mediaType != "song" && mediaType != "podcast" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "invalid mediaType parameter - must be 'song' or 'podcast'", "")
		return
	}

	offsetInt := 0
	if offset != "" {
		var err error
		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "invalid offset parameter - must be an integer", "")
			return
		}
	}

	transcodeParamsRow, err := database.GetTranscodeParams(r.Context(), transcodeParams)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "transcodeParams not found or expired", "")
		return
	}

	if transcodeParamsRow.MediaId != mediaId || transcodeParamsRow.MediaType != mediaType {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "transcodeParams do not match mediaId or mediaType", "")
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
		if requestUser.PodcastRole {
			episode, _ := database.GetPodcastEpisodeByGuid(ctx, mediaId)
			if episode.SourceUrl != "" {
				http.Redirect(w, r, episode.SourceUrl, http.StatusFound)
				return
			}
		}
	}

	err = ffmpeg.TranscodeAndStream(ctx, w, r, mediaFilepath, transcodeParamsRow.MediaId, transcodeParamsRow.Bitrate, offsetInt, transcodeParamsRow.TargetFormat)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error streaming audio", "")
		return
	}
}
