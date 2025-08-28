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

	response := subsonic.GetPopulatedSubsonicResponse(r.Context(), false)

	user, err := database.GetUserByContext(r.Context())
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error getting user from context", "")
		return
	}

	response.SubsonicResponse.TokenInfo = &types.TokenInfo{
		Username: user.Username,
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
