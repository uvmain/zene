package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"zene/core/auth"
	"zene/core/database"
)

func HandleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, _, err := auth.GetUserFromRequest(r)
	if err != nil {
		log.Printf("Failed to get user from request: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandleGetUserById(w http.ResponseWriter, r *http.Request) {
	userIdString := r.PathValue("userId")
	userId, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		log.Printf("Failed to get parse user ID: %v", err)
		http.Error(w, "Failed to get parse user ID", http.StatusInternalServerError)
		return
	}
	user, err := database.GetUserById(r.Context(), userId)
	if err != nil {
		log.Printf("Failed to get user from database: %v", err)
		http.Error(w, "Failed to get user from database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}
