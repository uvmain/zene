package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetTopSongs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()

	artistName := r.FormValue("artist")
	artistId := r.FormValue("id")
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
	count := r.FormValue("count")
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

	format := r.FormValue("f")
	if format == "json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))
		xml.NewEncoder(w).Encode(response.SubsonicResponse)
	}
}
