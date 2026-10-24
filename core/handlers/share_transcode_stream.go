package handlers

import (
	"net/http"
	"strconv"

	"zene/core/database"
	"zene/core/ffmpeg"
	"zene/core/logger"
	"zene/core/net"
)

func HandleShareTranscodeStream(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	mediaId := form["mediaid"]
	offset := form["offset"]
	transcodeParams := form["transcodeparams"]
	token := r.PathValue("token")
	ctx := r.Context()

	if mediaId == "" {
		http.Error(w, "mediaId parameter is required", http.StatusBadRequest)
		return
	}

	_, tokenIsValid := database.ValidateTokenAndMediaId(ctx, token, mediaId)
	if !tokenIsValid {
		logger.Printf("Error validating share token %s", token)
		http.Error(w, "invalid token and media_id combo", http.StatusBadRequest)
		return
	}

	offsetInt := 0
	if offset != "" {
		var err error
		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			http.Error(w, "invalid offset parameter - must be an integer", http.StatusBadRequest)
			return
		}
	}

	transcodeParamsRow, err := database.GetTranscodeParams(r.Context(), transcodeParams)
	if err != nil {
		http.Error(w, "transcodeParams not found or expired", http.StatusNotFound)
		return
	}

	if transcodeParamsRow.MediaId != mediaId {
		http.Error(w, "transcodeParams do not match mediaId", http.StatusBadRequest)
		return
	}

	mediaFilepath, err := database.GetMediaFilePath(ctx, mediaId)
	if err != nil {
		logger.Printf("Error getting media file path for %s: %v", mediaId, err)
		http.Error(w, "Error getting media file path", http.StatusInternalServerError)
		return
	}

	err = ffmpeg.TranscodeAndStream(ctx, w, r, mediaFilepath, transcodeParamsRow.MediaId, transcodeParamsRow.Bitrate, offsetInt, transcodeParamsRow.TargetFormat)
	if err != nil {
		if !net.IsClientDisconnectError(err) {
			http.Error(w, "Error streaming audio", http.StatusInternalServerError)
		}
		return
	}
}
