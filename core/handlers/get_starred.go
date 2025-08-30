package handlers

import (
	"net/http"
	"slices"
	"strconv"
	"strings"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetStarred(w http.ResponseWriter, r *http.Request) {
	var version int
	switch strings.ToLower(r.URL.Path) {
	case "/rest/getstarred":
		version = 1
	case "/rest/getstarred.view":
		version = 1
	case "/rest/getstarred2":
		version = 2
	case "/rest/getstarred2.view":
		version = 2
	}

	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	musicFolderIdParam := form["musicfolderid"]

	ctx := r.Context()
	var err error

	var musicFolderIdInt int
	if musicFolderIdParam != "" {
		musicFolderIdInt, err = strconv.Atoi(musicFolderIdParam)
		if err != nil || musicFolderIdInt < 0 {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "musicFolderId parameter must be a positive integer", "")
			return
		}
	} else {
		musicFolderIdInt = 0
	}

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to get users", "")
		return
	}

	// if musicFolderId not in requestUser.MusicFolderIds
	if musicFolderIdInt != 0 && !slices.Contains(requestUser.Folders, musicFolderIdInt) {
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to access this music folder", "")
		return
	}

	artists, err := database.GetStarredArtists(ctx, musicFolderIdInt)
	if err != nil {
		logger.Printf("Error getting starred artists: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to get starred artists", "")
		return
	}
	if artists == nil {
		artists = []types.Artist{}
	}

	albums, err := database.GetStarredAlbums(ctx, musicFolderIdInt)
	if err != nil {
		logger.Printf("Error getting starred albums: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to get starred albums", "")
		return
	}
	if albums == nil {
		albums = []types.AlbumId3{}
	}

	songs, err := database.GetStarredSongs(ctx, musicFolderIdInt)
	if err != nil {
		logger.Printf("Error getting starred songs: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to get starred songs", "")
		return
	}
	if songs == nil {
		songs = []types.SubsonicChild{}
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	switch version {
	case 1:
		response.SubsonicResponse.Starred = &types.Starred{}
		response.SubsonicResponse.Starred.Artists = artists
		response.SubsonicResponse.Starred.Albums = albums
		response.SubsonicResponse.Starred.Songs = songs
	case 2:
		response.SubsonicResponse.Starred2 = &types.Starred2{}
		response.SubsonicResponse.Starred2.Artists = artists
		response.SubsonicResponse.Starred2.Albums = albums
		response.SubsonicResponse.Starred2.Songs = songs
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
