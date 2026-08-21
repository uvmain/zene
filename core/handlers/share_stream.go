package handlers

import (
	"net/http"
	"strconv"
	"zene/core/config"
	"zene/core/database"
	"zene/core/ffmpeg"
	"zene/core/logger"
	"zene/core/net"
)

func HandleShareStream(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	token := r.PathValue("share_token")
	ctx := r.Context()
	form := net.NormalisedForm(r, w)
	mediaId := form["mediaid"]
	maxBitRateString := form["maxbitrate"]
	streamFormat := form["format"]
	timeOffsetString := form["timeoffset"]

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

	var maxBitRate int
	if maxBitRateString == "" {
		maxBitRate = config.DefaultBitRate
	} else {
		var err error
		maxBitRate, err = strconv.Atoi(maxBitRateString)
		if err != nil {
			http.Error(w, "maxBitRate parameter must be an integer", http.StatusBadRequest)
			return
		}
	}

	if streamFormat == "" {
		streamFormat = "aac" // default format
	}

	timeOffset := 0
	if timeOffsetString != "" {
		if timeOffsetInt, err := strconv.Atoi(timeOffsetString); err == nil && timeOffsetInt >= 0 {
			timeOffset = timeOffsetInt
		}
	}

	mediaFilepath, err := database.GetMediaFilePath(ctx, mediaId)

	if err != nil {
		logger.Printf("Error querying database for media filepath %s: %v", mediaId, err)
		http.Error(w, "File not found in database.", http.StatusNotFound)
		return
	}

	if mediaFilepath == "" {
		http.Error(w, "File not available to stream.", http.StatusNotFound)
		return
	}

	if streamFormat == "raw" {
		fileInfo, _, file, err := getFile(mediaFilepath)
		if err != nil {
			http.Error(w, "Error opening file.", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		logger.Printf("serving %s without transcoding", mediaFilepath)

		net.ServeFileWithRangeSupport(w, r, file, fileInfo.ModTime(), streamFormat)

		return
	}

	err = ffmpeg.TranscodeAndStream(ctx, w, r, mediaFilepath, mediaId, maxBitRate, timeOffset, streamFormat)
	if err != nil {
		if !net.IsClientDisconnectError(err) {
			http.Error(w, "Error streaming audio", http.StatusInternalServerError)
		}
		return
	}
}
