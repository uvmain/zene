package handlers

import (
	"encoding/json"
	"net/http"
	"zene/core/database"
	"zene/core/logger"
)

func HandleGetAlbums(w http.ResponseWriter, r *http.Request) {
	randomParam := r.FormValue("random")
	limitParam := r.FormValue("limit")
	recentParam := r.FormValue("recent")

	rows, err := database.SelectAllAlbums(r.Context(), randomParam, limitParam, recentParam)
	if err != nil {
		logger.Printf("Error querying database in SelectAllAlbums: %v", err)
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

func HandleGetAlbumTracks(w http.ResponseWriter, r *http.Request) {
	musicBrainzAlbumId := r.PathValue("musicBrainzAlbumId")

	rows, err := database.SelectTracksByAlbumId(r.Context(), musicBrainzAlbumId)
	if err != nil {
		logger.Printf("Error querying database in SelectTracksByAlbumId: %v", err)
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
