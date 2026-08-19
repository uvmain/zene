package handlers

import (
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleDeleteShare(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	id := form["id"]

	if id == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		logger.Printf("Error converting id to int: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid id parameter", "")
		return
	}

	ctx := r.Context()

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to delete shares", "")
		return
	}

	if !requestUser.ShareRole {
		logger.Printf("User %s attempted to delete a share without share role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to delete shares", "")
		return
	}

	err = database.DeleteShare(ctx, idInt)
	if err != nil {
		logger.Printf("Error deleting share: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to delete share", "")
		return
	}

	logger.Printf("User %s deleted share with ID %d", requestUser.Username, idInt)

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
