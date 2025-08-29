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

func HandleGetPlaylist(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	playlistId := form["id"]

	ctx := r.Context()

	if playlistId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter for playlist is required", "")
		return
	}

	playlistIdInt, err := strconv.Atoi(playlistId)
	if err != nil {
		logger.Printf("Error converting playlist id to int in GetPlaylist: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid id parameter for playlist", "")
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

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)
	response.SubsonicResponse.Playlist = &playlist

	playlistEntries, err := database.GetPlaylistEntries(ctx, playlistIdInt)
	if err != nil {
		logger.Printf("Error querying database for playlist entries in GetPlaylist: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to query database for playlist entries", "")
		return
	}
	response.SubsonicResponse.Playlist.Entries = playlistEntries

	net.WriteSubsonicResponse(w, r, response, format)
}
