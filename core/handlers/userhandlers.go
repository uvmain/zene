package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"zene/core/auth"
	"zene/core/database"
	"zene/core/types"
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

func HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetAllUsers(r.Context())
	if err != nil {
		log.Printf("Failed to get all users from database: %v", err)
		http.Error(w, "Failed to get all users from database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandlePostNewUser(w http.ResponseWriter, r *http.Request) {
	var newUser types.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Failed to decode request body", http.StatusInternalServerError)
		return
	}

	userId, err := database.UpsertUser(r.Context(), newUser)
	if err != nil {
		log.Printf("Failed to insert new user into database: %v", err)
		http.Error(w, "Failed to insert new user into database", http.StatusInternalServerError)
		return
	}

	var response struct {
		UserId int64 `json:"userId"`
	}
	response.UserId = userId

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandlePatchUserById(w http.ResponseWriter, r *http.Request) {
	userIdString := r.PathValue("userId")
	userIdInt, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		log.Printf("Failed to get parse user ID: %v", err)
		http.Error(w, "Failed to get parse user ID", http.StatusInternalServerError)
		return
	}

	var requestUser types.User
	err = json.NewDecoder(r.Body).Decode(&requestUser)
	if err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Failed to decode request body", http.StatusInternalServerError)
		return
	}

	username, err := database.GetUserById(r.Context(), userIdInt)
	if err != nil {
		log.Printf("Failed to validate user ID: %v", err)
		http.Error(w, "Failed to validate user ID", http.StatusInternalServerError)
		return
	}
	if username.Username != requestUser.Username {
		log.Println("Invalid user ID for user")
		http.Error(w, "Invalid user ID for user", http.StatusInternalServerError)
		return
	}

	userId, err := database.UpsertUser(r.Context(), requestUser)
	if err != nil {
		log.Printf("Failed to patch user: %v", err)
		http.Error(w, "Failed to patch user", http.StatusInternalServerError)
		return
	}

	var response struct {
		UserId int64 `json:"userId"`
	}
	response.UserId = userId

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func HandleDeleteUserById(w http.ResponseWriter, r *http.Request) {
	userIdString := r.PathValue("userId")
	userIdInt, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		log.Printf("Failed to get parse user ID: %v", err)
		http.Error(w, "Failed to get parse user ID", http.StatusInternalServerError)
		return
	}

	user, err := database.GetUserById(r.Context(), userIdInt)
	if err != nil {
		log.Printf("Failed to validate user ID: %v", err)
		http.Error(w, "Failed to validate user ID", http.StatusInternalServerError)
		return
	}

	allUsers, err := database.GetAllUsers(r.Context())
	if err != nil {
		log.Printf("Failed to getUsers from database: %v", err)
		http.Error(w, "Failed to getUsers from database", http.StatusInternalServerError)
		return
	}

	if len(allUsers) <= 1 {
		log.Printf("Cannot delete last user: %v", err)
		http.Error(w, "Cannot delete last user", http.StatusInternalServerError)
		return
	}

	err = database.DeleteUserById(r.Context(), user.Id)
	if err != nil {
		log.Printf("Failed to delete user from database: %v", err)
		http.Error(w, "Failed to delete user from database", http.StatusInternalServerError)
		return
	}

	var response struct {
		Status string `json:"status"`
	}

	response.Status = "success"

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}
