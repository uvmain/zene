package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleDeleteBookmark(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	id := form["id"]

	ctx := r.Context()

	if id == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	err := database.DeleteBookmark(ctx, id)
	if err != nil {
		logger.Printf("Error deleting bookmark: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to delete bookmark", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
