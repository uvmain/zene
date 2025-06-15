package net

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"zene/core/art"
	"zene/core/database"
	"zene/core/io"
	"zene/core/scanner"

	// "zene/core/scanner"
	"zene/core/types"
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

func HandleGetArtists(w http.ResponseWriter, r *http.Request) {
	searchParam := r.URL.Query().Get("search")
	randomParam := r.URL.Query().Get("random")
	recentParam := r.URL.Query().Get("recent")
	chronoParam := r.URL.Query().Get("chronological")
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")

	if randomParam != "" && randomParam != "true" && randomParam != "false" {
		http.Error(w, "Invalid value for 'random' parameter", http.StatusBadRequest)
		return
	}
	if recentParam != "" && recentParam != "true" && recentParam != "false" {
		http.Error(w, "Invalid value for 'recent' parameter", http.StatusBadRequest)
		return
	}

	if limitParam != "" {
		if _, err := strconv.Atoi(limitParam); err != nil {
			http.Error(w, "Invalid value for 'limit' parameter", http.StatusBadRequest)
			return
		}
	}

	rows, err := database.SelectAlbumArtists(r.Context(), searchParam, randomParam, recentParam, chronoParam, limitParam, offsetParam)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
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
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(row); err != nil {
		log.Println("Error encoding database response:", err)
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
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandleGetArtistArt(w http.ResponseWriter, r *http.Request) {
	musicBrainzArtistId := r.PathValue("musicBrainzArtistId")
	imageBlob, err := art.GetArtForArtist(r.Context(), musicBrainzArtistId)

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
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandleGetAlbums(w http.ResponseWriter, r *http.Request) {
	randomParam := r.URL.Query().Get("random")
	limitParam := r.URL.Query().Get("limit")
	recentParam := r.URL.Query().Get("recent")

	rows, err := database.SelectAllAlbums(r.Context(), randomParam, limitParam, recentParam)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandleGetAlbum(w http.ResponseWriter, r *http.Request) {
	musicBrainzAlbumId := r.PathValue("musicBrainzAlbumId")

	rows, err := database.SelectAlbum(r.Context(), musicBrainzAlbumId)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandleGetAlbumTracks(w http.ResponseWriter, r *http.Request) {
	musicBrainzAlbumId := r.PathValue("musicBrainzAlbumId")

	rows, err := database.SelectTracksByAlbumID(r.Context(), musicBrainzAlbumId)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandleGetGenres(w http.ResponseWriter, r *http.Request) {
	searchParam := r.URL.Query().Get("search")
	rows, err := database.SelectDistinctGenres(r.Context(), searchParam)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandleGetTracks(w http.ResponseWriter, r *http.Request) {
	randomParam := r.URL.Query().Get("random")
	limitParam := r.URL.Query().Get("limit")
	recentParam := r.URL.Query().Get("recent")

	rows, err := database.SelectAllTracks(r.Context(), randomParam, limitParam, recentParam)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandleGetTrack(w http.ResponseWriter, r *http.Request) {
	musicBrainzTrackId := r.PathValue("musicBrainzTrackId")

	row, err := database.SelectTrack(r.Context(), musicBrainzTrackId)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(row); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandlePostScan(w http.ResponseWriter, r *http.Request) {
	scanResult := scanner.RunScan(r.Context())
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(scanResult); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandleSearchMetadata(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")

	rows, err := database.SearchMetadata(r.Context(), searchQuery)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandleGetAlbumArt(w http.ResponseWriter, r *http.Request) {
	musicBrainzAlbumId := r.PathValue("musicBrainzAlbumId")
	sizeParam := r.URL.Query().Get("size")
	if sizeParam == "" {
		sizeParam = "xl"
	}
	imageBlob, err := art.GetArtForAlbum(r.Context(), musicBrainzAlbumId, sizeParam)

	if err != nil {
		http.Redirect(w, r, "/default-square.png", http.StatusTemporaryRedirect)
		return
	}
	mimeType := http.DetectContentType(imageBlob)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(imageBlob)
}
