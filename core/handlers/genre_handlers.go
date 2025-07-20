package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/types"
)

func HandleGetGenres(w http.ResponseWriter, r *http.Request) {
	searchParam := r.URL.Query().Get("search")
	limitParam := r.URL.Query().Get("limit")
	rows, err := database.SelectDistinctGenres(r.Context(), limitParam, searchParam)
	if err != nil {
		logger.Printf("Error querying database in SelectDistinctGenres: %v", err)
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

func HandleGetTracksByGenre(w http.ResponseWriter, r *http.Request) {
	condition := r.URL.Query().Get("condition")
	genres := r.URL.Query().Get("genres")
	limit := r.URL.Query().Get("limit")
	random := r.URL.Query().Get("random")

	limitInt := 0
	var err error

	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil || limitInt < 0 {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
	}

	if random != "true" {
		random = "false"
	}

	if genres == "" {
		http.Error(w, "No genres provided", http.StatusBadRequest)
		return
	}

	if condition != "and" && condition != "or" {
		condition = "and"
	}

	genresList := []string{}
	for _, genre := range strings.Split(genres, ",") {
		trimmedGenre := strings.TrimSpace(genre)
		if trimmedGenre != "" {
			genresList = append(genresList, trimmedGenre)
		}
	}

	var rows []types.MetadataWithPlaycounts

	rows, err = database.SelectTracksByGenres(r.Context(), genresList, condition, int64(limitInt), random)
	if err != nil {
		logger.Printf("Error querying database in SelectTracksByGenres: %v", err)
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
