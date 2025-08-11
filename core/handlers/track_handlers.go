package handlers

import (
	"encoding/json"
	"net/http"
	"zene/core/database"
	"zene/core/io"
	"zene/core/logger"
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
	randomParam := r.FormValue("random")
	limitParam := r.FormValue("limit")
	offsetParam := r.FormValue("offset")
	recentParam := r.FormValue("recent")
	chronoParam := r.FormValue("chronological")

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

func HandleSearchMetadata(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.FormValue("search")

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
