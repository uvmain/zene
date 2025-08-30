package handlers

import (
	"net/http"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetVideos(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}
	net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "getVideos endpoint is unsupported", "")
}
