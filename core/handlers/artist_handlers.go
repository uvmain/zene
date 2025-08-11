package handlers

import (
	"encoding/json"
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/types"
)

func HandleGetArtists(w http.ResponseWriter, r *http.Request) {
	searchParam := r.FormValue("search")
	randomParam := r.FormValue("random")
	recentParam := r.FormValue("recent")
	chronoParam := r.FormValue("chronological")
	limitParam := r.FormValue("limit")
	offsetParam := r.FormValue("offset")

	rows, err := database.SelectAlbumArtists(r.Context(), searchParam, randomParam, recentParam, chronoParam, limitParam, offsetParam)
	if err != nil {
		logger.Printf("Error querying database in SelectAlbumArtists: %v", err)
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

func HandleGetArtist(w http.ResponseWriter, r *http.Request) {
	musicBrainzArtistId := r.PathValue("musicBrainzArtistId")
	var row types.ArtistResponse
	var err error

	row, err = database.SelectArtistByMusicBrainzArtistId(r.Context(), musicBrainzArtistId)

	if err != nil {
		logger.Printf("Error querying database in SelectArtistByMusicBrainzArtistId: %v", err)
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

func HandleGetArtistTracks(w http.ResponseWriter, r *http.Request) {
	musicBrainzArtistId := r.PathValue("musicBrainzArtistId")
	randomParam := r.FormValue("random")
	limitParam := r.FormValue("limit")
	offsetParam := r.FormValue("offset")
	recentParam := r.FormValue("recent")

	rows, err := database.SelectTracksByArtistId(r.Context(), musicBrainzArtistId, randomParam, limitParam, offsetParam, recentParam)
	if err != nil {
		logger.Printf("Error querying database in SelectTracksByArtistId: %v", err)
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

func HandleGetArtistAlbums(w http.ResponseWriter, r *http.Request) {
	musicBrainzArtistId := r.PathValue("musicBrainzArtistId")
	randomParam := r.FormValue("random")
	chronoParam := r.FormValue("chronological")
	limitParam := r.FormValue("limit")
	offsetParam := r.FormValue("offset")
	recentParam := r.FormValue("recent")

	rows, err := database.SelectAlbumsByArtistId(r.Context(), musicBrainzArtistId, randomParam, recentParam, chronoParam, limitParam, offsetParam)
	if err != nil {
		logger.Printf("Error querying database in SelectAlbumsByArtistId: %v", err)
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
