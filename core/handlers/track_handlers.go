package handlers

import (
	"encoding/json"
	"net/http"
	"zene/core/database"
	"zene/core/logger"
)

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
