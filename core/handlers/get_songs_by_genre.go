package handlers

import (
	"fmt"
	"slices"
	"strconv"

	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetSongsByGenre(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	genre := form["genre"]
	count := form["count"]
	offset := form["offset"]
	musicFolderId := form["musicfolderid"]

	ctx := r.Context()

	ifModifiedSinceHeader := r.Header.Get("If-Modified-Since")
	if ifModifiedSinceHeader != "" {
		latestScan, err := database.GetLatestCompletedScan(ctx)
		if err == nil {
			latestScanTime := logic.GetStringTimeFormatted(latestScan.CompletedDate)
			if net.IfModifiedResponse(w, r, latestScanTime) {
				return
			}
		}
	}

	if genre == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "genre parameter is required", "")
		return
	}

	var countInt int
	if count != "" {
		var err error
		countInt, err = strconv.Atoi(count)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "count parameter must be an integer", "")
			return
		}
	} else {
		countInt = 10 // default to ten if param is not provided
	}

	var offsetInt int
	if offset != "" {
		var err error
		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "offset parameter must be an integer", "")
			return
		}
	} else {
		offsetInt = 0 // default to zero if param is not provided
	}

	var musicFolderIdInt int
	if musicFolderId != "" {
		var err error
		musicFolderIdInt, err = strconv.Atoi(musicFolderId)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "musicFolderId parameter must be an integer", "")
			return
		}
		user, err := database.GetUserByContext(ctx)
		if err == nil {
			userHasAccessToMusicFolder := slices.Contains(user.Folders, musicFolderIdInt)
			if !userHasAccessToMusicFolder {
				net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, fmt.Sprintf("musicFolderId %d is not valid for user", musicFolderIdInt), "")
				return
			}
		}
	} else {
		musicFolderIdInt = 0
	}

	songs, err := database.GetSongsByGenre(ctx, genre, countInt, offsetInt, musicFolderIdInt)
	if err != nil {
		logger.Printf("Error getting songs by genre: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "No songs found for genre", "")
		return
	}
	if songs == nil {
		songs = []types.SubsonicChild{}
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)
	response.SubsonicResponse.SongsByGenre = &types.SongsByGenre{
		Songs: songs,
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
