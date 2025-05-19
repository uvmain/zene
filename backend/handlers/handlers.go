package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"zene/art"
	"zene/database"
	"zene/net"
	"zene/scanner"
	"zene/types"
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

func HandleGetFileByName(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Missing filename parameter", http.StatusBadRequest)
		return
	}

	row, err := database.SelectFileByFilename(filename)

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
	rows, err := database.SelectAllAlbums()
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
	rows, err := database.SelectAllMetadata()
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

func HandleGetRandomMetadataWithLimit(w http.ResponseWriter, r *http.Request) {
	var data []types.TrackMetadata
	var err error

	limitParam := r.URL.Query().Get("limit")
	log.Printf("Limit parameter: %s", limitParam)

	if limitParam == "" {
		data, err = database.SelectAllMetadataRandomized()
		if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}
	} else {
		limitParamInt, err := strconv.Atoi(limitParam)
		data, err = database.SelectRandomMetadataWithLimit(limitParamInt)
		if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
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
