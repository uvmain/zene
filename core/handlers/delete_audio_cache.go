package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleDeleteAudioCache(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	ctx := r.Context()

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "context user not found", "")
		return
	}

	if !requestUser.AdminRole {
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "user not authorized to delete audio cache", "")
		return
	}

	err = database.DeleteAllAudioCache(ctx)
	if err != nil {
		logger.Printf("Error deleting audio cache database entries: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorInvalidApiKey, "Server Error", "")
		return
	}

	err = io.DeleteAllAudioCacheFiles(ctx)
	if err != nil {
		logger.Printf("Error deleting audio cache files: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorInvalidApiKey, "Server Error", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
