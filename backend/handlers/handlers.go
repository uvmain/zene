package handlers

import (
	"encoding/json"
	"net/http"
	"zene/database"
)

func HandleGetAllFiles(w http.ResponseWriter, r *http.Request) {
	rows, err := database.SelectAllFiles()
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
		}
	}
}

func HandleGetArtists(w http.ResponseWriter, r *http.Request) {
	rows, err := database.SelectAllArtists()
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
