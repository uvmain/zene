package handlers

import (
	"net/http"
	"zene/core/art"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetAvatar(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	username := form["username"]

	ctx := r.Context()

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to get avatars", "")
		return
	}

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

	mimeType := http.DetectContentType(avatarBlob)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(avatarBlob)
}
