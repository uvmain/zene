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

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	var version int
	switch strings.ToLower(r.URL.Path) {
	case "/rest/search":
		version = 1
	case "/rest/search.view":
		version = 1
	case "/rest/search2":
		version = 2
	case "/rest/search2.view":
		version = 2
	case "/rest/search3":
		version = 3
	case "/rest/search3.view":
		version = 3
	}

	if version == 1 {
		net.WriteSubsonicError(w, r, types.ErrorIncompatibleClient, "Endpoint deprecated - use search2 instead", "/rest/search2.view")
		return
	}

	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	searchQueryParam := form["query"]
	artistCountParam := form["artistcount"]
	artistOffsetParam := form["artistoffset"]
	albumCountParam := form["albumcount"]
	albumOffsetParam := form["albumoffset"]
	songCountParam := form["songcount"]
	songOffsetParam := form["songoffset"]
	musicFolderIdParam := form["musicfolderid"]

	if searchQueryParam == "\"\"" || searchQueryParam == "''" {
		searchQueryParam = ""
	}

	ctx := r.Context()
	var err error

	var artistCount int
	if artistCountParam != "" {
		artistCount, err = strconv.Atoi(artistCountParam)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "artistCount parameter must be an integer", "")
			return
		}
	} else {
		artistCount = 20
	}

	var artistOffset int
	if artistOffsetParam != "" {
		artistOffset, err = strconv.Atoi(artistOffsetParam)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "artistOffset parameter must be an integer", "")
			return
		}
	} else {
		artistOffset = 0
	}

	var albumCount int
	if albumCountParam != "" {
		albumCount, err = strconv.Atoi(albumCountParam)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "albumCount parameter must be an integer", "")
			return
		}
	} else {
		albumCount = 20
	}

	var albumOffset int
	if albumOffsetParam != "" {
		albumOffset, err = strconv.Atoi(albumOffsetParam)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "albumOffset parameter must be an integer", "")
			return
		}
	} else {
		albumOffset = 0
	}

	var songCount int
	if songCountParam != "" {
		songCount, err = strconv.Atoi(songCountParam)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "songCount parameter must be an integer", "")
			return
		}
	} else {
		songCount = 20
	}

	var songOffset int
	if songOffsetParam != "" {
		songOffset, err = strconv.Atoi(songOffsetParam)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "songOffset parameter must be an integer", "")
			return
		}
	} else {
		songOffset = 0
	}

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

	artists := []types.Artist{}
	if artistCount > 0 {
		artists, err = database.SearchArtists(ctx, searchQueryParam, artistCount, artistOffset, musicFolderIdInt)
		if err != nil {
			logger.Printf("Error searching artists: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to search artists", "")
			return
		}
	}

	albums := []types.AlbumId3{}
	if albumCount > 0 {
		albums, err = database.SearchAlbums(ctx, searchQueryParam, albumCount, albumOffset, musicFolderIdInt)
		if err != nil {
			logger.Printf("Error searching albums: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to search albums", "")
			return
		}
	}

	songs := []types.SubsonicChild{}
	if songCount > 0 {
		songs, err = database.SearchSongs(ctx, searchQueryParam, songCount, songOffset, musicFolderIdInt)
		if err != nil {
			logger.Printf("Error searching songs: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to search songs", "")
			return
		}
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	switch version {
	case 2:
		searchResult2 := types.SearchResult2{}
		response.SubsonicResponse.SearchResult2 = &searchResult2
		response.SubsonicResponse.SearchResult2.Artists = artists
		response.SubsonicResponse.SearchResult2.Albums = albums
		response.SubsonicResponse.SearchResult2.Songs = songs
	case 3:
		searchResult3 := types.SearchResult3{}
		response.SubsonicResponse.SearchResult3 = &searchResult3
		response.SubsonicResponse.SearchResult3.Artists = artists
		response.SubsonicResponse.SearchResult3.Albums = albums
		response.SubsonicResponse.SearchResult3.Songs = songs
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
