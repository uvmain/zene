package handlers

import (
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
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	var version int
	switch strings.ToLower(r.URL.Path) {
	case "/rest/getsimilarsongs":
		version = 1
	case "/rest/getsimilarsongs.view":
		version = 1
	case "/rest/getsimilarsongs2":
		version = 2
	case "/rest/getsimilarsongs2.view":
		version = 2
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	musicbrainzId := form["id"]
	count := form["count"]

	if musicbrainzId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

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

	songs, err := database.GetSimilarSongs(ctx, countInt, musicbrainzId)
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

	net.WriteSubsonicResponse(w, r, response, format)
}
