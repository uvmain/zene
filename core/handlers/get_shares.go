package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetShares(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	ctx := r.Context()

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to get shares", "")
		return
	}

	if !requestUser.ShareRole {
		logger.Printf("User %s attempted to get shares without share role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to get shares", "")
		return
	}

	shares, err := database.GetSharesByUser(ctx)
	if err != nil {
		logger.Printf("Error getting shares: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to get shares", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	response.SubsonicResponse.Shares = &types.Shares{}
	response.SubsonicResponse.Shares.Share = shares

	net.WriteSubsonicResponse(w, r, response, format)
}
