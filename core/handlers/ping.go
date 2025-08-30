package handlers

import (
	"net/http"
	"zene/core/net"
	"zene/core/subsonic"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]

	response := subsonic.GetPopulatedSubsonicResponse(r.Context())

	net.WriteSubsonicResponse(w, r, response, format)
}
