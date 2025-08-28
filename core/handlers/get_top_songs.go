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

func HandleGetTopSongs(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	artistName := form["artist"]
	artistId := form["id"]
	count := form["count"]

	ctx := r.Context()

	if artistName == "" && artistId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "artist or id parameter is required", "")
		return
	}

	if artistId != "" && artistName == "" {
		var err error
		artistName, err = database.GetArtistNameById(ctx, artistId)
		if err != nil {
			logger.Printf("failed to get artist name by ID %s: %v", artistId, err)
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "artist not found", "")
			return
		}
		if artistName == "" {
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "artist not found", "")
			return
		}
	}

	var countLimit = 50
	var err error
	if count != "" {
		countLimit, err = strconv.Atoi(count)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "count parameter must be an integer", "")
			return
		}
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	artistTopSongs, err := database.SelectTopSongsForArtistName(ctx, artistName, countLimit)
	if err != nil {
		logger.Printf("failed to get top songs from database: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "failed to get top songs", "")
		return
	}

	response.SubsonicResponse.TopSongs = &types.TopSongs{Songs: artistTopSongs}

	net.WriteSubsonicResponse(w, r, response, format)
}
