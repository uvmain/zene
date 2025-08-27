package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetSimilarSongs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	var version int
	switch strings.ToLower(r.URL.Path) {
	case "/rest/getsimilarsongs.view":
		version = 1
	case "/rest/getsimilarsongs2.view":
		version = 2
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	artistId := form["id"]
	count := form["count"]

	ctx := r.Context()

	var countInt int
	if count != "" {
		var err error
		countInt, err = strconv.Atoi(count)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "count parameter must be an integer", "")
			return
		}
	} else {
		countInt = 50 // default to 50 if param is not provided
	}

	songs, err := database.GetSimilarSongs(ctx, countInt, artistId)
	if err != nil {
		logger.Printf("Error getting similar songs: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error getting similar songs", "")
		return
	}
	if songs == nil {
		songs = []types.SubsonicChild{}
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	switch version {
	case 1:
		response.SubsonicResponse.SimilarSongs = &types.SimilarSongs{}
		response.SubsonicResponse.SimilarSongs.Songs = songs
	case 2:
		response.SubsonicResponse.SimilarSongs2 = &types.SimilarSongs2{}
		response.SubsonicResponse.SimilarSongs2.Songs = songs
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
