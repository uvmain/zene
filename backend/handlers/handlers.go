package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"zene/art"
	"zene/database"
	"zene/net"
	"zene/scanner"
	"zene/types"
)

func HandleGetFiles(w http.ResponseWriter, r *http.Request) {
	rows, err := database.SelectAllFiles()
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		return
	}
}

func HandleGetFile(w http.ResponseWriter, r *http.Request) {
	fileId := r.PathValue("fileId")
	if fileId == "" {
		http.Error(w, "Missing fileId parameter", http.StatusBadRequest)
		return
	}

	row, err := database.SelectFileByFileId(fileId)

	if err != nil {
		log.Printf("Error querying database: %v", err)
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

func HandleDownloadFile(w http.ResponseWriter, r *http.Request) {
	fileId := r.PathValue("fileId")
	fileBlob, err := database.GetFileBlob(fileId)

	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	mimeType := http.DetectContentType(fileBlob)
	net.EnableCdnCaching(w)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(fileBlob)
}

func HandleGetArtists(w http.ResponseWriter, r *http.Request) {
	searchParam := r.URL.Query().Get("search")
	var rows []types.ArtistResponse
	var err error

	if searchParam == "" {
		rows, err = database.SelectAllArtists()
	} else {
		rows, err = database.SearchForArtists(searchParam)
	}
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		return
	}
}

func HandleGetArtist(w http.ResponseWriter, r *http.Request) {
	musicBrainzArtistId := r.PathValue("musicBrainzArtistId")
	var row types.ArtistResponse
	var err error

	row, err = database.SelectArtistByMusicBrainzArtistId(musicBrainzArtistId)

	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(row); err != nil {
		log.Println("Error encoding database response:", err)
		return
	}
}

func GetArtistArt(w http.ResponseWriter, r *http.Request) {
	musicBrainzArtistId := r.PathValue("musicBrainzArtistId")
	imageBlob, err := art.GetArtForArtist(musicBrainzArtistId)

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

func HandleGetAlbums(w http.ResponseWriter, r *http.Request) {
	randomParam := r.URL.Query().Get("random")
	limitParam := r.URL.Query().Get("limit")
	recentParam := r.URL.Query().Get("recent")

	rows, err := database.SelectAllAlbums(randomParam, limitParam, recentParam)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		return
	}
}

func HandleGetAlbum(w http.ResponseWriter, r *http.Request) {
	musicBrainzAlbumId := r.PathValue("musicBrainzAlbumId")

	rows, err := database.SelectAlbum(musicBrainzAlbumId)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		return
	}
}

func HandleGetGenres(w http.ResponseWriter, r *http.Request) {
	searchParam := r.URL.Query().Get("search")
	rows, err := database.SelectDistinctGenres(searchParam)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		return
	}
}

func HandleGetTracks(w http.ResponseWriter, r *http.Request) {
	randomParam := r.URL.Query().Get("random")
	limitParam := r.URL.Query().Get("limit")
	recentParam := r.URL.Query().Get("recent")

	rows, err := database.SelectAllTracks(randomParam, limitParam, recentParam)
	if err != nil {
		log.Printf("Error querying database: %v", err)
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

func HandleSearchMetadata(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")

	rows, err := database.SearchMetadata(searchQuery)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		log.Println("Error encoding database response:", err)
		return
	}
}

func HandleGetAlbumArt(w http.ResponseWriter, r *http.Request) {
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
