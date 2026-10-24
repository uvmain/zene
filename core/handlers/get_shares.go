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

	if !requestUser.ShareRole && !requestUser.AdminRole {
		logger.Printf("User %s attempted to get shares without share or admin role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to get shares", "")
		return
	}

	var shares []types.ShareRow

	if requestUser.AdminRole {
		shares, err = database.GetAllShares(ctx)
		if err != nil {
			logger.Printf("Error getting all shares for admin user %s: %v", requestUser.Username, err)
			net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to get shares", "")
			return
		}
	} else {
		shares, err = database.GetSharesByUser(ctx)
		if err != nil {
			logger.Printf("Error getting shares for user %s: %v", requestUser.Username, err)
			net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to get shares", "")
			return
		}
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	response.SubsonicResponse.Shares = &types.Shares{}
	response.SubsonicResponse.Shares.Share = shares

	net.WriteSubsonicResponse(w, r, response, format)
}
