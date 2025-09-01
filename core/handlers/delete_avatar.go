package handlers

import (
	"net/http"
	"zene/core/art"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleDeleteAvatar(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
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

	if username == "" {
		username = requestUser.Username
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

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
