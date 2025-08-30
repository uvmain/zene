package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleLicense(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]

	response := subsonic.GetPopulatedSubsonicResponse(r.Context())

	user, _ := database.GetUserByContext(r.Context())

	response.SubsonicResponse.License = &types.LicenseInfo{
		Valid:          true,
		Email:          user.Username,
		LicenseExpires: "",
		TrialExpires:   "",
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
