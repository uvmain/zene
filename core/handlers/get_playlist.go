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

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	playlistIdInt, err := strconv.Atoi(playlistId)
	if err != nil {
		logger.Printf("Error converting playlistId to int in GetPlaylist: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid playlist ID", "")
		return
	}

	playlistExists, err := database.PlaylistExists(ctx, playlistIdInt, "")
	if err != nil {
		logger.Printf("Error checking if playlist exists in GetPlaylist: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to check if playlist exists", "")
		return
	}
	if !playlistExists {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Playlist not found", "")
		return
	}

	playlist, err := database.GetPlaylist(ctx, playlistIdInt)
	if err != nil {
		logger.Printf("Error querying database in GetPlaylist: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to query database", "")
		return
	}
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
