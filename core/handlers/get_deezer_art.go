package handlers

import (
	"encoding/json"

	"net/http"
	"zene/core/deezer"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetDeezerArt(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	artistName := form["artist"]
	albumName := form["album"]

	ctx := r.Context()

	if artistName == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "artist parameter is required", "")
		return
	}

	if albumName == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "album parameter is required", "")
		return
	}

	deezerImageUrl, err := deezer.GetAlbumArtUrlWithArtistNameAndAlbumName(ctx, artistName, albumName)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "album art not found", "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"url": deezerImageUrl,
	}); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
