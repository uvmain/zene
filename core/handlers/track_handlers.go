package handlers

import (
	"encoding/json"
	"net/http"
	"zene/core/database"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/lyrics"
)

func HandleDownloadTrack(w http.ResponseWriter, r *http.Request) {
	musicBrainzTrackId := r.PathValue("musicBrainzTrackId")
	track, err := database.SelectTrack(r.Context(), musicBrainzTrackId)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	fileBlob, err := io.GetFileBlob(r.Context(), track.FilePath)

	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	mimeType := http.DetectContentType(fileBlob)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(fileBlob)
}

func HandleGetTracks(w http.ResponseWriter, r *http.Request) {
	randomParam := r.URL.Query().Get("random")
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")
	recentParam := r.URL.Query().Get("recent")
	chronoParam := r.URL.Query().Get("chronological")

	rows, err := database.SelectAllTracks(r.Context(), randomParam, limitParam, offsetParam, recentParam, chronoParam)
	if err != nil {
		logger.Printf("Error querying database in SelectAllTracks: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		logger.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandleGetTrack(w http.ResponseWriter, r *http.Request) {
	musicBrainzTrackId := r.PathValue("musicBrainzTrackId")

	row, err := database.SelectTrack(r.Context(), musicBrainzTrackId)
	if err != nil {
		logger.Printf("Error querying database in SelectTrack: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(row); err != nil {
		logger.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandleGetTrackLyrics(w http.ResponseWriter, r *http.Request) {
	musicBrainzTrackId := r.PathValue("musicBrainzTrackId")

	lyrics, err := lyrics.GetLyricsForMusicBrainzTrackId(r.Context(), musicBrainzTrackId)
	if err != nil {
		logger.Printf("Error querying lyrics in GetLyricsForMusicBrainzTrackId: %v", err)
		http.Error(w, "Failed to fetch lyrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lyrics); err != nil {
		logger.Println("Error encoding lyrics response:", err)
		http.Error(w, "Error encoding lyrics response", http.StatusInternalServerError)
		return
	}
}

func HandleSearchMetadata(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")

	rows, err := database.SearchMetadata(r.Context(), searchQuery)
	if err != nil {
		logger.Printf("Error querying database in SearchMetadata: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		logger.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}
