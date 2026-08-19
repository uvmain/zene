package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetShare(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	token := r.PathValue("share_token")
	ctx := r.Context()

	share, err := database.GetShareByToken(ctx, token)
	if err != nil {
		logger.Printf("Error getting share by token %s: %v", token, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to get share", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	response.SubsonicResponse.Shares = &types.Shares{}
	response.SubsonicResponse.Shares.Share = []types.ShareRow{}

	response.SubsonicResponse.Shares.Share = append(response.SubsonicResponse.Shares.Share, share)

	net.WriteSubsonicResponse(w, r, response, format)
}
