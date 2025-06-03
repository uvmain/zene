package net

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"zene/core/art"
	"zene/core/database"
	"zene/core/scanner"
	"zene/core/types"
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
	enableCdnCaching(w)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(fileBlob)
}

func HandleStreamFile(w http.ResponseWriter, r *http.Request) {
	fileId := r.PathValue("fileId")
	file, err := database.SelectFileByFileId(fileId)

	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	f, err := os.Open(file.FilePath)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		http.Error(w, "file stat error", http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
}

func HandleGetArtists(w http.ResponseWriter, r *http.Request) {
	searchParam := r.URL.Query().Get("search")
	randomParam := r.URL.Query().Get("random")
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")
	recentParam := r.URL.Query().Get("recent")

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

	rows, err := database.SelectAlbumArtists(searchParam, randomParam, limitParam, offsetParam, recentParam)
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
		http.Redirect(w, r, "/default-square.png", http.StatusTemporaryRedirect)
		return
	}
	mimeType := http.DetectContentType(imageBlob)
	enableCdnCaching(w)
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

func HandleGetAlbumTracks(w http.ResponseWriter, r *http.Request) {
	musicBrainzAlbumId := r.PathValue("musicBrainzAlbumId")

	rows, err := database.SelectTracksByAlbumID(musicBrainzAlbumId)
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

func HandleGetTrack(w http.ResponseWriter, r *http.Request) {
	musicBrainzTrackId := r.PathValue("musicBrainzTrackId")

	row, err := database.SelectTrack(musicBrainzTrackId)
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
		http.Redirect(w, r, "/default-square.png", http.StatusTemporaryRedirect)
		return
	}
	mimeType := http.DetectContentType(imageBlob)
	enableCdnCaching(w)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(imageBlob)
}
