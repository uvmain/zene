package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
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
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
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

	genre := r.FormValue("genre")
	if genre == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "genre parameter is required", "")
		return
	}

	var countInt int
	count := r.FormValue("count")
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
	offset := r.FormValue("offset")
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

	// TODO limiting to musicfolderId not yet implemented
	// var musicFolderIdInt int
	// musicFolderId := r.FormValue("musicFolderId")
	// if musicFolderId != "" {
	// 	var err error
	// 	musicFolderIdInt, err = strconv.Atoi(musicFolderId)
	// 	if err != nil {
	// 		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "musicFolderId parameter must be an integer", "")
	// 		return
	// 	}
	// } else {
	// 	musicFolderIdInt = 0
	// }

	songs, err := database.GetSongsByGenre(ctx, genre, countInt, offsetInt)
	if err != nil {
		logger.Printf("Error getting songs by genre: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "No songs found for genre", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)
	response.SubsonicResponse.SongsByGenre = &types.SongsByGenre{
		Songs: songs,
	}

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
