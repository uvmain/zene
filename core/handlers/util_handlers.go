package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"zene/core/scanner"
)

func HandlePostScan(w http.ResponseWriter, r *http.Request) {
	scanResult := scanner.RunScan(r.Context())
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(scanResult); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}
