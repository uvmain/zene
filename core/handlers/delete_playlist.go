package handlers

import (
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleDeletePlaylist(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	playlistId := form["id"]

	ctx := r.Context()

	playlistIdInt, err := strconv.Atoi(playlistId)
	if err != nil {
		logger.Printf("Error converting playlistId to int in DeletePlaylist: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid playlist ID", "")
		return
	}

	playlist, err := database.GetPlaylist(ctx, playlistIdInt)
	if err != nil {
		logger.Printf("Error checking if playlist exists in GetPlaylist: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to check if playlist exists", "")
		return
	}
	if playlist.Id < 1 {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Playlist not found", "")
		return
	}

	user, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context in DeletePlaylist: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to get user by context", "")
		return
	}

	if playlist.Owner != user.Username && !user.AdminRole {
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to delete this playlist", "")
		return
	}

	if err := database.DeletePlaylist(ctx, playlistIdInt); err != nil {
		logger.Printf("Error deleting playlist in DeletePlaylist: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to delete playlist", "")
		return
	}

	logger.Printf("Playlist %d (%s) deleted by user %s", playlistIdInt, playlist.Name, user.Username)

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	net.WriteSubsonicResponse(w, r, response, format)
}
