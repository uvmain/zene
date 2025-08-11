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
	"zene/core/net"
	"zene/core/types"
)

func HandleStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()
	musicBrainzTrackId := r.FormValue("id")
	if musicBrainzTrackId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	var maxBitRate int64
	maxBitRateString := r.FormValue("maxBitRate")
	if maxBitRateString == "" {
		maxBitRate = config.DefaultBitRate
	} else {
		var err error
		maxBitRate, err = strconv.ParseInt(maxBitRateString, 10, 64)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "maxBitRate parameter must be an integer", "")
			return
		}
	}

	format := r.FormValue("format")
	if format == "" {
		format = "aac" // default format
	}

	timeOffsetString := r.FormValue("timeOffset")
	timeOffset := 0
	if timeOffsetString != "" {
		if timeOffsetInt, err := strconv.Atoi(timeOffsetString); err == nil && timeOffsetInt >= 0 {
			timeOffset = timeOffsetInt
		}
	}

	size := r.FormValue("size")
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
		http.Error(w, "File not found in database.", http.StatusNotFound)
		return
	}

	if format == "raw" {
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

	err = ffmpeg.TranscodeAndStream(ctx, w, r, track.FilePath, musicBrainzTrackId, maxBitRate, timeOffset, format)
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
