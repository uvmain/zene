package handlers

import (
	"strconv"

	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleUpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	playlistId := form["playlistid"]
	playlistName := form["name"]
	comment := form["comment"]
	public := form["public"]
	coverArt := form["coverart"]

	allowedUsers, _, err := net.ParseDuplicateFormKeys(r, "allowedUserId", true)
	if err != nil {
		logger.Printf("Error parsing allowedUserId: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid allowedUserId", "")
		return
	}

	_, songIdsToAdd, err := net.ParseDuplicateFormKeys(r, "songIdToAdd", false)
	if err != nil {
		logger.Printf("Error parsing songIdToAdd: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid songIdToAdd", "")
		return
	}

	songIndexesToRemove, _, err := net.ParseDuplicateFormKeys(r, "songIndexToRemove", true)
	if err != nil {
		logger.Printf("Error parsing songIndexToRemove: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid songIndexToRemove", "")
		return
	}

	ctx := r.Context()

	if playlistId == "" && playlistName == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "either playlistId or name parameter is required", "")
		return
	}

	var playlistIdInt int
	if playlistId != "" {
		var err error
		playlistIdInt, err = strconv.Atoi(playlistId)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "playlistId parameter must be an integer", "")
			return
		}
	}

	playlistExists, err := database.PlaylistExists(ctx, playlistIdInt, playlistName)
	if err != nil {
		logger.Printf("Error getting playlist: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error getting playlist", "")
		return
	}
	if !playlistExists {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "playlist does not exist", "")
		return
	}

	err = database.UpdatePlaylist(ctx, playlistIdInt, playlistName, comment, public, coverArt, allowedUsers, songIdsToAdd, songIndexesToRemove)
	if err != nil && err.Error() == "existing playlist provided with no new songIds" {
		logger.Printf("Error creating playlist: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error creating playlist, existing playlist provided with no new songIds", "")
		return
	} else if err != nil && err.Error() == "existing playlists should be referenced by playlistId, not name" {
		logger.Printf("Error creating playlist: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error creating playlist, existing playlists should be referenced by playlistId, not name", "")
		return
	} else if err != nil {
		logger.Printf("Error creating playlist: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error creating playlist", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	net.WriteSubsonicResponse(w, r, response, format)
}
