package handlers

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
	"zene/core/config"
	"zene/core/database"
	"zene/core/ffmpeg"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func HandleStream(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	streamId := form["id"]
	maxBitRateString := form["maxbitrate"]
	streamFormat := form["format"]
	timeOffsetString := form["timeoffset"]
	size := form["size"]

	ctx := r.Context()

	if streamId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "User not found", "")
		return
	}

	var maxBitRate int
	if maxBitRateString == "" {
		maxBitRate = config.DefaultBitRate
	} else {
		var err error
		maxBitRate, err = strconv.Atoi(maxBitRateString)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "maxBitRate parameter must be an integer", "")
			return
		}
	}

	if requestUser.MaxBitRate > 0 && requestUser.MaxBitRate < maxBitRate {
		maxBitRate = requestUser.MaxBitRate
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

	if size != "" {
		// check size is in "WxH" syntax, eg "640x480"
		re := regexp.MustCompile(`^\d+x\d+$`)
		if !re.MatchString(size) {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "size parameter must be in 'WxH' format", "")
			return
		}
	}

	// This parameter is unused - we send the Content-Length header anyway, even if they don't ask for it
	// var estimateContentLength bool
	// estimateContentLengthString := r.FormValue("estimateContentLength")
	// if estimateContentLengthString != "" {
	// 	estimateContentLength = net.ParseBooleanFromString(w, r, estimateContentLengthString)
	// }

	// This parameter is unused - it is for Video, which Zene does not support
	// var converted bool
	// convertedString := r.FormValue("converted")
	// if convertedString != "" {
	// 	converted = net.ParseBooleanFromString(w, r, convertedString)
	// }

	mediaFilepath, err := database.GetMediaFilePath(ctx, streamId)

	if mediaFilepath == "" || err != nil {
		// check if the file is a podcast episode
		if requestUser.PodcastRole {
			logger.Printf("Checking if streamId %s is a podcast episode", streamId)
			episode, _ := database.GetPodcastEpisodeByGuid(ctx, streamId)
			if episode.SourceUrl != "" {
				http.Redirect(w, r, episode.SourceUrl, http.StatusFound)
				return
			}
		}
	}

	if err != nil {
		logger.Printf("Error querying database for media filepath %s: %v", streamId, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "File not found in database.", "")
		return
	}

	if mediaFilepath == "" {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "File not available to stream.", "")
		return
	}

	if streamFormat == "raw" {
		fileInfo, modTime, file, err := OpenFile(mediaFilepath)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error opening file.", "")
			return
		}
		defer file.Close()

		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		contentType := mime.TypeByExtension(filepath.Ext(fileInfo.Name()))
		w.Header().Set("Content-Type", contentType)
		http.ServeContent(w, r, fileInfo.Name(), modTime, file)
		return
	}

	err = ffmpeg.TranscodeAndStream(ctx, w, r, mediaFilepath, streamId, maxBitRate, timeOffset, streamFormat)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error streaming audio", "")
		return
	}
}

func OpenFile(filePath string) (os.FileInfo, time.Time, *os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, time.Time{}, nil, fmt.Errorf("opening file: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, time.Time{}, nil, fmt.Errorf("stating file: %w", err)
	}

	return stat, stat.ModTime(), file, nil
}
