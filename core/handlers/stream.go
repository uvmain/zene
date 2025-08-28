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
	musicBrainzTrackId := form["id"]
	maxBitRateString := form["maxbitrate"]
	streamFormat := form["format"]
	timeOffsetString := form["timeoffset"]
	size := form["size"]

	ctx := r.Context()

	if musicBrainzTrackId == "" {
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

	// This is unused - we send the header anyway, even if they don't ask for it
	// var estimateContentLength bool
	// estimateContentLengthString := r.FormValue("estimateContentLength")
	// if estimateContentLengthString != "" {
	// 	estimateContentLength = net.ParseBooleanFromString(w, r, estimateContentLengthString)
	// }

	// var converted bool
	// convertedString := r.FormValue("converted")
	// if convertedString != "" {
	// 	converted = net.ParseBooleanFromString(w, r, convertedString)
	// }

	track, err := database.SelectTrack(ctx, musicBrainzTrackId)
	if err != nil {
		logger.Printf("Error querying database for track %s: %v", musicBrainzTrackId, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "File not found in database.", "")
		return
	}

	if streamFormat == "raw" {
		fileInfo, modTime, file, err := OpenFile(track.FilePath)
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

	err = ffmpeg.TranscodeAndStream(ctx, w, r, track.FilePath, musicBrainzTrackId, maxBitRate, timeOffset, streamFormat)
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
