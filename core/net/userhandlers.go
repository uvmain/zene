package net

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	// chi is no longer used for path parameters as router.go uses http.ServeMux
	// "github.com/go-chi/chi/v5"
	"github.com/ollama/ollama/api"
	"github.com/ollama/ollama/core/auth"
	"github.com/ollama/ollama/core/database"
	"github.com/ollama/ollama/core/types"
	"github.com/ollama/ollama/llm"
	log "github.com/sirupsen/logrus"
)

// isAdmin is a helper function to check if the authenticated user is an admin.
// It retrieves the user from the context, which should have been set by the AuthMiddleware.
func isAdmin(r *http.Request) bool {
	user, ok := r.Context().Value(auth.UserContextKey).(*types.User)
	if !ok || user == nil {
		return false // Should not happen if AuthMiddleware is applied
	}
	return user.IsAdmin
}

// GetCurrentUserHandler handles requests to get the current authenticated user's details.
func GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserContextKey).(*types.User)
	if !ok || user == nil {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	// Return user data (excluding password hash)
	responseUser := types.User{
		Id:       user.Id,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseUser); err != nil {
		log.Errorf("Error encoding current user response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetUsersHandler handles requests to list all users. Only accessible by admins.
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "Forbidden: You do not have permission to view users.", http.StatusForbidden)
		return
	}

	db := r.Context().Value(database.DBContextKey).(*database.DB)
	if db == nil {
		http.Error(w, "Database connection not found in context", http.StatusInternalServerError)
		return
	}

	users, err := database.GetAllUsers(r.Context(), db)
	if err != nil {
		log.Errorf("Error fetching users: %v", err)
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	// Ensure password hashes are not included in the response
	for i := range users {
		users[i].PasswordHash = ""
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Errorf("Error encoding users list response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// CreateUserRequest defines the expected JSON structure for creating a user.
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"` // Changed from is_admin to match Vue component
}

// CreateUserHandler handles requests to create a new user. Only accessible by admins.
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "Forbidden: You do not have permission to create users.", http.StatusForbidden)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Errorf("Error hashing password for user %s: %v", req.Username, err)
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	user := types.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		IsAdmin:      req.IsAdmin,
	}

	db := r.Context().Value(database.DBContextKey).(*database.DB)
	if db == nil {
		http.Error(w, "Database connection not found in context", http.StatusInternalServerError)
		return
	}

	// UpsertUser will handle ID being 0 for creation
	if err := database.UpsertUser(r.Context(), db, user); err != nil {
		log.Errorf("Error creating user %s: %v", req.Username, err)
		// Check for unique constraint error (example, specific error handling might be needed)
		if err.Error().Contains("UNIQUE constraint failed: users.username") { // This check is DB specific
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Fetch the created user to return its ID (without password hash)
	createdUser, err := database.GetUserByUsername(r.Context(), db, req.Username)
	if err != nil {
		log.Errorf("Error fetching created user %s: %v", req.Username, err)
		// Even if user was created, if we can't fetch it, it's an issue.
		// For simplicity, return 201 Created without the full user object if this happens,
		// or handle more gracefully.
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "User created successfully, but failed to retrieve details."}`))
		return
	}
	createdUser.PasswordHash = "" // Clear password hash

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		log.Errorf("Error encoding created user response: %v", err)
		// Already sent 201, so can't send another error header. Log is important.
	}
}

// UpdateUserRequest defines the expected JSON structure for updating a user.
// Password is optional; if provided, it will be updated.
type UpdateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"` // Optional password
	IsAdmin  bool   `json:"is_admin"`
}

// UpdateUserHandler handles requests to update an existing user. Only accessible by admins.
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "Forbidden: You do not have permission to update users.", http.StatusForbidden)
		return
	}

	// Use r.PathValue for Go 1.22+ http.ServeMux path parameters
	userIDStr := r.PathValue("userId")
	if userIDStr == "" {
		log.Warn("UpdateUserHandler: userId path parameter is empty.")
		http.Error(w, "User ID is required in path", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Errorf("UpdateUserHandler: Invalid user ID format '%s': %v", userIDStr, err)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	db := r.Context().Value(database.DBContextKey).(*database.DB)
	if db == nil {
		http.Error(w, "Database connection not found in context", http.StatusInternalServerError)
		return
	}

	// Fetch existing user to update
	user, err := database.GetUserById(r.Context(), db, userID)
	if err != nil {
		log.Errorf("Error fetching user %d for update: %v", userID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.Username = req.Username
	user.IsAdmin = req.IsAdmin

	if req.Password != "" {
		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			log.Errorf("Error hashing password for user update %d: %v", userID, err)
			http.Error(w, "Failed to process password", http.StatusInternalServerError)
			return
		}
		user.PasswordHash = hashedPassword
	}
	// If password is not provided in request, user.PasswordHash remains the existing one.

	if err := database.UpsertUser(r.Context(), db, user); err != nil {
		log.Errorf("Error updating user %d: %v", userID, err)
		if err.Error().Contains("UNIQUE constraint failed: users.username") {
			http.Error(w, "Username already exists for another user", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	updatedUser, err := database.GetUserById(r.Context(), db, userID)
	if err != nil {
		log.Errorf("Error fetching updated user %d: %v", userID, err)
		http.Error(w, "User updated, but failed to retrieve details", http.StatusInternalServerError)
		return
	}
	updatedUser.PasswordHash = ""


	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedUser); err != nil {
		log.Errorf("Error encoding updated user response: %v", err)
	}
}

// DeleteUserHandler handles requests to delete a user. Only accessible by admins.
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "Forbidden: You do not have permission to delete users.", http.StatusForbidden)
		return
	}

	// Use r.PathValue for Go 1.22+ http.ServeMux path parameters
	userIDStr := r.PathValue("userId")
	if userIDStr == "" {
		log.Warn("DeleteUserHandler: userId path parameter is empty.")
		http.Error(w, "User ID is required in path", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Errorf("DeleteUserHandler: Invalid user ID format '%s': %v", userIDStr, err)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Prevent admin from deleting themselves (optional safeguard)
	currentUser, _ := r.Context().Value(auth.UserContextKey).(*types.User)
	if currentUser != nil && currentUser.Id == userID {
		http.Error(w, "Cannot delete your own user account.", http.StatusBadRequest)
		return
	}

	db := r.Context().Value(database.DBContextKey).(*database.DB)
	if db == nil {
		http.Error(w, "Database connection not found in context", http.StatusInternalServerError)
		return
	}

	if err := database.DeleteUser(r.Context(), db, userID); err != nil {
		log.Errorf("Error deleting user %d: %v", userID, err)
		if err.Error() == "no rows affected" { // Assuming DeleteUser returns specific error for not found
			http.Error(w, "User not found or already deleted", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User deleted successfully"}`))
}

// This function is not directly used by user handlers but is kept for consistency with the original file structure.
// It might be used by other parts of the application or intended for future use.
func GetRunner(ctx context.Context, model string) (llm.LLM, error) {
	return GetLocalRunner(ctx, model, "", api.DefaultOptions(), false, nil)
}
