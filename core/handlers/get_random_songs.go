package handlers

import (
	"encoding/json"
	"encoding/xml"
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

func HandleGetRandomSongs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	count := form["size"]
	genre := form["genre"]
	fromYear := form["fromyear"]
	toYear := form["toyear"]
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

	var fromYearInt int
	var toYearInt int
	var err error

	if fromYear != "" {
		fromYearInt, err = strconv.Atoi(fromYear)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "fromYear parameter must be an integer", "")
			return
		}
	}
	if toYear != "" {
		toYearInt, err = strconv.Atoi(toYear)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "toYear parameter must be an integer", "")
			return
		}
	}
	if toYearInt > 0 && toYearInt < fromYearInt {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "toYear parameter must be greater than or equal to fromYear", "")
		return
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

	songs, err := database.GetRandomSongs(ctx, countInt, genre, fromYear, toYear, musicFolderIdInt)
	if err != nil {
		logger.Printf("Error getting random songs: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error getting random songs", "")
		return
	}
	if songs == nil {
		songs = []types.SubsonicChild{}
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)
	response.SubsonicResponse.RandomSongs = &types.RandomSongs{
		Songs: songs,
	}

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
