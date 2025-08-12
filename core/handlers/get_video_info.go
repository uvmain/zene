package handlers

import (
	"fmt"
	"net/http"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetVideoInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}
	net.WriteSubsonicError(w, r, types.ErrorIncompatibleVersion, "getVideoInfo endpoint is unsupported", "")
}
