package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"zene/core/auth"
	"zene/core/database"
	"zene/core/types"
)

func HandleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, _, err := auth.GetUserFromRequest(r)
	response := &UserResponse{}

	if err != nil {
		handleErrorResponse(w, response, "Failed to get user from request", err, http.StatusInternalServerError)
		return
	}

	response.User = toUserPointer(user)
	handleSuccessResponse(w, response)
}

func HandleGetUserById(w http.ResponseWriter, r *http.Request) {
	response := &UserResponse{}
	userIdString := r.PathValue("userId")

	userId, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		handleErrorResponse(w, response, "Failed to get parse user ID", err, http.StatusBadRequest)
		return
	}
	user, err := database.GetUserById(r.Context(), userId)
	if err != nil {
		handleErrorResponse(w, response, "Failed to get user from database", err, http.StatusInternalServerError)
		return
	}

	response.User = toUserPointer(user)
	handleSuccessResponse(w, response)
	return
}

func HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	response := &UsersResponse{}

	users, err := database.GetAllUsers(r.Context())
	if err != nil {
		handleErrorResponse(w, response, "Failed to get all users from database", err, http.StatusInternalServerError)
		return
	}

	response.Users = toUsersPointers(users)
	handleSuccessResponse(w, response)
	return
}

func HandlePostNewUser(w http.ResponseWriter, r *http.Request) {
	var newUser types.User
	response := &IdResponse{}

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		handleErrorResponse(w, response, "Failed to decode request body", err, http.StatusInternalServerError)
		return
	}

	userId, err := database.UpsertUser(r.Context(), newUser)
	if err != nil {
		handleErrorResponse(w, response, "Failed to insert new user into database", err, http.StatusInternalServerError)
		return
	}

	response.Id = userId
	handleSuccessResponse(w, response)
	return
}

func HandlePatchUserById(w http.ResponseWriter, r *http.Request) {
	userIdString := r.PathValue("userId")
	response := &IdResponse{}

	userIdInt, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		handleErrorResponse(w, response, "Failed to parse user ID", err, http.StatusInternalServerError)
		return
	}

	var requestUser types.User
	err = json.NewDecoder(r.Body).Decode(&requestUser)
	if err != nil {
		handleErrorResponse(w, response, "Failed to decode request body", err, http.StatusInternalServerError)
		return
	}

	username, err := database.GetUserById(r.Context(), userIdInt)
	if err != nil {
		handleErrorResponse(w, response, "Failed to validate user ID", err, http.StatusBadRequest)
		return
	}

	if username.Username != requestUser.Username {
		handleErrorResponse(w, response, "Invalid user ID for user", err, http.StatusBadRequest)
		return
	}

	userId, err := database.UpsertUser(r.Context(), requestUser)
	if err != nil {
		handleErrorResponse(w, response, "Failed to patch user", err, http.StatusInternalServerError)
		return
	}

	response.Id = userId
	handleSuccessResponse(w, response)
	return
}

func HandleDeleteUserById(w http.ResponseWriter, r *http.Request) {
	userIdString := r.PathValue("userId")
	response := &StandardResponse{}
	userIdInt, err := strconv.ParseInt(userIdString, 10, 64)

	if err != nil {
		handleErrorResponse(w, response, "Failed to parse user ID", err, http.StatusInternalServerError)
		return
	}

	user, err := database.GetUserById(r.Context(), userIdInt)
	if err != nil {
		handleErrorResponse(w, response, "Failed to validate user ID", err, http.StatusNotFound)
		return
	}

	allUsers, err := database.GetAllUsers(r.Context())
	if err != nil {
		handleErrorResponse(w, response, "Failed to get users from database", err, http.StatusInternalServerError)
		return
	}

	if len(allUsers) <= 1 {
		handleErrorResponse(w, response, "Cannot delete last user", fmt.Errorf(""), http.StatusBadRequest)
		return
	}

	err = database.DeleteUserById(r.Context(), user.Id)
	if err != nil {
		handleErrorResponse(w, response, "Failed to delete user from database", err, http.StatusInternalServerError)
		return
	}

	handleSuccessResponse(w, response)
	return
}
