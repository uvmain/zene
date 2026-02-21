package handlers

import (
	"net/http"
	"strconv"
	"strings"
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
	typeParam := strings.ToLower(form["type"])
	sizeParam := form["size"]
	offsetParam := form["offset"]
	seedParam := form["seed"]

	musicFolderIds, _, err := net.ParseDuplicateFormKeys(r, "musicfolderid", true)
	if err != nil {
		musicFolderIds = []int{}
	}

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

	/*
		type param is required and can be one of:
		starred
		random
		newest (recently modified)
		highest (highest rated)
		frequent (most frequently played)
		recent (recently played)
		alphabetical

		if type=random, accept an optional `seed` param (integer) to get deterministic results
	*/
	if typeParam == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "type parameter is required", "")
		return
	} else if typeParam != "starred" && typeParam != "random" && typeParam != "newest" && typeParam != "highest" &&
		typeParam != "frequent" && typeParam != "recent" && typeParam != "alphabetical" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid type parameter", "")
		return
	}

	var sizeInt int
	if sizeParam != "" {
		sizeInt, err = strconv.Atoi(sizeParam)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "size parameter must be an integer", "")
			return
		}
	} else {
		sizeInt = 0
	}

	var offsetInt int
	if offsetParam != "" {
		offsetInt, err = strconv.Atoi(offsetParam)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "offset parameter must be an integer", "")
			return
		}
	} else {
		offsetInt = 0
	}

	if sizeInt <= 0 && offsetInt > 0 {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "size must be greater than 0 if offset is greater than 0", "")
		return
	}

	var seedInt int
	if seedParam != "" {
		seedInt, err = strconv.Atoi(seedParam)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "seed parameter must be an integer", "")
			return
		}
	} else {
		seedInt = 0
	}

	artists, err := database.GetArtistList(ctx, musicFolderIds, sizeInt, offsetInt, typeParam, seedInt)
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
