package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"zene/core/database"
	"zene/core/lyrics"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetLyrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	artist := form["artist"]
	title := form["title"]

	ctx := r.Context()

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	if artist == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "artist parameter is required", "")
		return
	}

	if title == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "title parameter is required", "")
		return
	}

	musicBrainzTrackId, err := database.GetTrackIdByArtistAndTitle(artist, title)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, fmt.Sprintf("Error fetching track ID: %v", err), "")
		return
	}

	lyricsData, err := lyrics.GetLyricsForMusicBrainzTrackId(ctx, musicBrainzTrackId)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, fmt.Sprintf("Error fetching lyrics: %v", err), "")
		return
	}

	response.SubsonicResponse.Lyrics = &types.SubsonicLyrics{
		Artist: artist,
		Title:  title,
		Value:  lyricsData.PlainLyrics,
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
