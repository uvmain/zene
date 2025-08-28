package handlers

import (
	"fmt"
	"net/http"
	"zene/core/database"
	"zene/core/lyrics"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetLyrics(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
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

	net.WriteSubsonicResponse(w, r, response, format)
}
