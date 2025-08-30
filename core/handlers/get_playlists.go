package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetPlaylists(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	username := form["username"]

	ctx := r.Context()

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to get users", "")
		return
	}
	if username != "" && requestUser.Username != username && !requestUser.AdminRole {
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to get this user", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	playlistUsername := username
	if playlistUsername == "" {
		playlistUsername = requestUser.Username
	}

	playlists, err := database.GetPlaylists(ctx, playlistUsername)
	if err != nil {
		logger.Printf("Error querying database in GetPlaylistList: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to query database", "")
		return
	}
	if playlists == nil {
		playlists = []types.PlaylistRow{}
	}

	response.SubsonicResponse.Playlists = &types.Playlists{}
	response.SubsonicResponse.Playlists.Playlist = playlists

	net.WriteSubsonicResponse(w, r, response, format)
}
