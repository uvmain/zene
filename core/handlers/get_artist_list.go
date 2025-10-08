package handlers

import (
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetArtistList(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	ifModifiedSince := form["ifmodifiedsince"]

	musicFolderIds, _, err := net.ParseDuplicateFormKeys(r, "musicfolderid", true)
	if err != nil {
		musicFolderIds = []int{}
	}
	logger.Printf("Parsed music folder IDs: %v", musicFolderIds)

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

	if ifModifiedSince != "" {
		ifModifiedSinceInt, err = strconv.Atoi(ifModifiedSince)
		if err != nil || ifModifiedSinceInt < 0 {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "ifModifiedSince parameter must be a positive integer", "")
			return
		}
	} else {
		ifModifiedSinceInt = 0
	}

	artists, err := database.GetArtistList(ctx, musicFolderIds)
	if err != nil {
		logger.Printf("Error querying database in GetIndexes: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to query database", "")
		return
	}
	if artists == nil {
		artists = []types.Artist{}
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)
	response.SubsonicResponse.ArtistList = &types.SubsonicArtistListWrapper{}
	response.SubsonicResponse.ArtistList.Artists = artists

	net.WriteSubsonicResponse(w, r, response, format)
}
