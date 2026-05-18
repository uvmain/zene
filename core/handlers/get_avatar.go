package handlers

import (
	"net/http"
	"strconv"
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
	id := form["id"]

	if username == "" && id == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Missing required parameter: username or id", "")
		return
	}

	ctx := r.Context()

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to get avatars", "")
		return
	}

	var avatarUser types.User

	if id != "" {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			logger.Printf("Error converting id to int: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid id parameter", "")
			return
		}
		avatarUser, err = database.GetUserById(ctx, idInt)
		if err != nil {
			logger.Printf("Error getting user ID for user ID %d: %v", idInt, err)
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "User not found", "")
			return
		}
	} else {
		avatarUser, err = database.GetUserByUsername(ctx, username)
		if err != nil {
			logger.Printf("Error getting user ID for username %s: %v", username, err)
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "User not found", "")
			return
		}
	}

	if !requestUser.AdminRole && avatarUser.Id != requestUser.Id {
		logger.Printf("User %s attempted to fetch avatars for another user without admin role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to get avatars for another user", "")
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
	_, err = w.Write(avatarBlob)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "failed to write file", "")
		return
	}
}
