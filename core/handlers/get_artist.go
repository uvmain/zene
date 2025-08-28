package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetArtist(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	musicBrainzArtistId := form["id"]

	ctx := r.Context()

	if musicBrainzArtistId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	row, err := database.SelectArtistByMusicBrainzArtistId(ctx, musicBrainzArtistId)
	if err != nil {
		logger.Printf("Error querying database in SelectArtistByMusicBrainzArtistId: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to query database", "")
		return
	}

	response.SubsonicResponse.Artist = &types.SubsonicArtistWrapper{}
	response.SubsonicResponse.Artist.Artist = row

	net.WriteSubsonicResponse(w, r, response, format)
}
