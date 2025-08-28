package handlers

import (
	"net/http"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetCaptions(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}
	net.WriteSubsonicError(w, r, types.ErrorIncompatibleVersion, "getCaptions endpoint is unsupported", "")
}
