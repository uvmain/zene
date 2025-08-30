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

func HandleCreatePlaylist(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	playlistId := form["playlistid"]
	playlistName := form["name"]

	_, songIds, err := net.ParseDuplicateFormKeys(r, "songId", false)
	if err != nil {
		logger.Printf("Error parsing songIds: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid songIds", "")
		return
	}

	ctx := r.Context()

	if playlistId == "" && playlistName == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "either playlistId or name parameter is required", "")
		return
	}

	if playlistId != "" && len(songIds) == 0 {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "songId param should not be empty if playlistId is provided", "")
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

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	result, err := database.CreatePlaylist(ctx, playlistName, playlistIdInt, songIds)
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

	response.SubsonicResponse.Playlist = &result

	net.WriteSubsonicResponse(w, r, response, format)
}
