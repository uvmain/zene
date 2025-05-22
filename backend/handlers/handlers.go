package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"zene/art"
	"zene/database"
	"zene/net"
	"zene/scanner"
)

func HandleGetAllFiles(w http.ResponseWriter, r *http.Request) {
	rows, err := database.SelectAllFiles()
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		return
	}
}

func HandleGetFileById(w http.ResponseWriter, r *http.Request) {
	fileId := r.PathValue("fileId")
	if fileId == "" {
		http.Error(w, "Missing fileId parameter", http.StatusBadRequest)
		return
	}

	row, err := database.SelectFileByFileId(fileId)

	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(row); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func HandleGetArtists(w http.ResponseWriter, r *http.Request) {
	rows, err := database.SelectAllArtists()
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		return
	}
}

func HandleGetAlbums(w http.ResponseWriter, r *http.Request) {
	randomParam := r.URL.Query().Get("random")
	limitParam := r.URL.Query().Get("limit")
	recentParam := r.URL.Query().Get("recent")

	rows, err := database.SelectAllAlbums(randomParam, limitParam, recentParam)
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		return
	}
}

func HandleGetMetadata(w http.ResponseWriter, r *http.Request) {
	randomParam := r.URL.Query().Get("random")
	limitParam := r.URL.Query().Get("limit")
	recentParam := r.URL.Query().Get("recent")

	rows, err := database.SelectAllMetadata(randomParam, limitParam, recentParam)
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		return
	}
}

func HandlePostScan(w http.ResponseWriter, r *http.Request) {
	scanResult := scanner.RunScan()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(scanResult); err != nil {
		log.Println("Error encoding database response:", err)
		return
	}
}

func GetAlbumArtByMusicBrainzAlbumId(w http.ResponseWriter, r *http.Request) {
	musicBrainzAlbumId := r.PathValue("musicBrainzAlbumId")
	sizeParam := r.URL.Query().Get("size")
	if sizeParam == "" {
		sizeParam = "xl"
	}
	imageBlob, err := art.GetArtForAlbum(musicBrainzAlbumId, sizeParam)

	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	mimeType := http.DetectContentType(imageBlob)
	net.EnableCdnCaching(w)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(imageBlob)
}
