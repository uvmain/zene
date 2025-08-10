package handlers

import (
	"encoding/json"
	"net/http"
	"zene/core/art"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
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

func HandleGetAlbum(w http.ResponseWriter, r *http.Request) {
	musicBrainzAlbumId := r.PathValue("musicBrainzAlbumId")

	rows, err := database.SelectAlbum(r.Context(), musicBrainzAlbumId)
	if err != nil {
		logger.Printf("Error querying database in SelectAlbum: %v", err)
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

func HandleGetAlbumArt(w http.ResponseWriter, r *http.Request) {
	musicBrainzAlbumId := r.PathValue("musicBrainzAlbumId")
	sizeParam := r.FormValue("size")
	if sizeParam == "" {
		sizeParam = "xl"
	}
	imageBlob, lastModified, err := art.GetArtForAlbum(r.Context(), musicBrainzAlbumId, sizeParam)

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
