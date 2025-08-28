package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
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
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to delete users", "")
		return
	}

	if requestUser.AdminRole == false {
		logger.Printf("User %s attempted to delete a user without admin role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to delete users", "")
		return
	}

	if username == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Username is required", "")
		return
	}

	userToDelete, err := database.GetUserByUsername(ctx, username)
	if err != nil || userToDelete.Id <= 0 {
		logger.Printf("Error checking if user exists: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "User does not exist", "")
		return
	}

	err = database.DeleteUserById(ctx, userToDelete.Id)
	if err != nil {
		logger.Printf("Error deleting user %s: %v", username, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to delete user", "")
		return
	}

	logger.Printf("User %s deleted with ID %d by %s", username, userToDelete.Id, requestUser.Username)

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	net.WriteSubsonicResponse(w, r, response, format)
}
