package handlers

import (
	"encoding/json"
	"net/http"
	"zene/core/art"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetArtists(w http.ResponseWriter, r *http.Request) {
	searchParam := r.URL.Query().Get("search")
	randomParam := r.URL.Query().Get("random")
	recentParam := r.URL.Query().Get("recent")
	chronoParam := r.URL.Query().Get("chronological")
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")

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
	randomParam := r.URL.Query().Get("random")
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")
	recentParam := r.URL.Query().Get("recent")

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

func HandleGetArtistArt(w http.ResponseWriter, r *http.Request) {
	musicBrainzArtistId := r.PathValue("musicBrainzArtistId")
	imageBlob, lastModified, err := art.GetArtForArtist(r.Context(), musicBrainzArtistId)

	if net.IfModifiedResponse(w, r, lastModified) {
		return
	}

	if err != nil {
		http.Redirect(w, r, "/default-square.png", http.StatusTemporaryRedirect)
		return
	}

	mimeType := http.DetectContentType(imageBlob)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(imageBlob)
}

func HandleGetArtistAlbums(w http.ResponseWriter, r *http.Request) {
	musicBrainzArtistId := r.PathValue("musicBrainzArtistId")
	randomParam := r.URL.Query().Get("random")
	chronoParam := r.URL.Query().Get("chronological")
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")
	recentParam := r.URL.Query().Get("recent")

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
