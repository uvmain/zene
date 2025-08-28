package handlers

import (
	"net/http"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetVideoInfo(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}
	net.WriteSubsonicError(w, r, types.ErrorIncompatibleVersion, "getVideoInfo endpoint is unsupported", "")
}
