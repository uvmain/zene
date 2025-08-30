package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleTokenInfo(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	token := form["apikey"]

	if token == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "apiKey parameter is required", "")
		return
	}

	ctx := r.Context()

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error getting user from context", "")
		return
	}

	tokenUser, err := database.ValidateApiKey(ctx, token)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error validating API key", "")
		return
	}

	if requestUser.Id != tokenUser.Id && !requestUser.AdminRole {
		net.WriteSubsonicError(w, r, types.ErrorInvalidApiKey, "Invalid API key", "")
		return
	}

	response.SubsonicResponse.TokenInfo = &types.TokenInfo{
		Username: tokenUser.Username,
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
