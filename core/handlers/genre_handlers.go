package handlers

import (
	"encoding/json"
	"net/http"
	"zene/core/database"
	"zene/core/logger"
)

func HandleGetGenres(w http.ResponseWriter, r *http.Request) {
	searchParam := r.URL.Query().Get("search")
	rows, err := database.SelectDistinctGenres(r.Context(), searchParam)
	if err != nil {
		logger.Printf("Error querying database: %v", err)
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
