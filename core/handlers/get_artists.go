package handlers

import (
	"slices"

	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetArtists(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	ifModifiedSince := form["ifmodifiedsince"]
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

	var ifModifiedSinceInt int
	var musicFolderIdInt int
	var err error

	if ifModifiedSince != "" {
		ifModifiedSinceInt, err = strconv.Atoi(ifModifiedSince)
		if err != nil || ifModifiedSinceInt < 0 {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "ifModifiedSince parameter must be a positive integer", "")
			return
		}
	} else {
		ifModifiedSinceInt = 0
	}

	if musicFolderId != "" {
		musicFolderIdInt, err := strconv.Atoi(musicFolderId)
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

	queryMusicFolderInts := []int{}
	// if the musicFolderId param is valid, use it - otherwise, use the user's folders
	if musicFolderIdInt != 0 {
		queryMusicFolderInts = append(queryMusicFolderInts, int(musicFolderIdInt))
	} else {
		queryMusicFolderInts = append(queryMusicFolderInts, requestUser.Folders...)
	}

	indexes, err := database.GetIndexes(ctx, requestUser.Id, queryMusicFolderInts, ifModifiedSinceInt)
	if err != nil {
		logger.Printf("Error querying database in GetIndexes: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to query database", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)
	response.SubsonicResponse.Artists = &types.SubsonicArtistsWrapper{
		Artists:         indexes.Indexes,
		IgnoredArticles: "",
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
