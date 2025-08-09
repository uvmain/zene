package handlers

import (
	"fmt"
	"net/http"
	"zene/core/art"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetAvatar(w http.ResponseWriter, r *http.Request) {
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

	avatarBlob, err := art.GetUserAvatarImage(avatarUser.Id)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error getting avatar image", "")
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(avatarBlob)
}
