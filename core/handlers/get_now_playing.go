package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"zene/core/database"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetNowPlaying(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
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
