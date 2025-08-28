package handlers

import (
	"fmt"
	"net/http"
	"zene/core/database"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetNowPlaying(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]

	ctx := r.Context()

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	response.SubsonicResponse.NowPlaying = &types.SubsonicNowPlaying{}

	nowPlayingEntries, err := database.GetNowPlaying(ctx)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, fmt.Sprintf("Failed to get now playing: %v", err), "")
		return
	}

	response.SubsonicResponse.NowPlaying.Entry = nowPlayingEntries

	net.WriteSubsonicResponse(w, r, response, format)
}
