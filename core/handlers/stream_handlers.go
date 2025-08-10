package handlers

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"zene/core/database"
	"zene/core/ffmpeg"
)

func HandleStreamTrack(w http.ResponseWriter, r *http.Request) {
	musicBrainzTrackId := r.PathValue("musicBrainzTrackId")
	qualityQuery := r.FormValue("quality")
	ctx := r.Context()

	track, err := database.SelectTrack(ctx, musicBrainzTrackId)

	if err != nil {
		http.Error(w, "File not found in database.", http.StatusNotFound)
		return
	}

	if qualityQuery == "native" {
		name, modTime, file, err := OpenFile(track.FilePath)
		if err != nil {
			http.Error(w, "Error opening file.", 500)
			return
		}
		defer file.Close()
		contentType := mime.TypeByExtension(filepath.Ext(name))
		w.Header().Set("Content-Type", contentType)
		http.ServeContent(w, r, name, modTime, file)
	} else if qualityQuery != "" {
		qualityInt, err := strconv.Atoi(qualityQuery)
		if err != nil {
			http.Error(w, "Error converting quality parameter to an integer.", 500)
			return
		}
		err = ffmpeg.TranscodeAndStream(ctx, w, r, track.FilePath, musicBrainzTrackId, qualityInt)
		if err != nil {
			http.Error(w, "Error streaming audio", http.StatusInternalServerError)
			return
		}
	} else {
		name, modTime, file, err := OpenFile(track.FilePath)
		if err != nil {
			http.Error(w, "Error opening file.", 500)
			return
		}
		defer file.Close()
		contentType := mime.TypeByExtension(filepath.Ext(name))
		w.Header().Set("Content-Type", contentType)
		http.ServeContent(w, r, name, modTime, file)
	}
}

func OpenFile(filePath string) (string, time.Time, *os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", time.Time{}, nil, fmt.Errorf("opening file: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return "", time.Time{}, nil, fmt.Errorf("stating file: %w", err)
	}

	return stat.Name(), stat.ModTime(), file, nil
}
