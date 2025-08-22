package handlers

import (
	"net/http"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	unknownEndpoint := r.PathValue("unknownEndpoint")

	logger.Printf("Request for path not found: %s", unknownEndpoint)
	net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Resource not found", "")
}
