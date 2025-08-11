package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"zene/core/art"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleDeleteAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to get avatars", "")
		return
	}

	username := r.FormValue("username")

	if requestUser.AdminRole == false && username == requestUser.Username {
		logger.Printf("User %s attempted to fetch avatars for another user without admin role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to get avatars for another user", "")
		return
	}

	avatarUser, err := database.GetUserByUsername(ctx, username)
	if err != nil {
		logger.Printf("Error getting user ID for username %s: %v", username, err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "User not found", "")
		return
	}

	err = art.DeleteUserAvatarImage(avatarUser.Id)
	if err != nil {
		logger.Printf("Error deleting avatar image: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error deleting avatar image", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	format := r.FormValue("f")
	if format == "json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))
		xml.NewEncoder(w).Encode(response.SubsonicResponse)
	}
}
