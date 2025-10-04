package handlers

import (
	"encoding/json"

	"net/http"
	"zene/core/database"
	"zene/core/deezer"
	"zene/core/musicbrainz"
	"zene/core/net"
	"zene/core/types"
)

type AlbumArtsResponse struct {
	Deezer          string `json:"deezer"`
	CoverArtArchive string `json:"cover_art_archive"`
}

func HandleGetAlbumArts(w http.ResponseWriter, r *http.Request) {
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

	album, err := database.GetAlbumByArtistNameAndAlbumName(ctx, artistName, albumName)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "album not found", "")
		return
	}

	deezerImageUrl, err := deezer.GetAlbumArtUrlWithArtistNameAndAlbumName(ctx, artistName, albumName)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "album art not found", "")
		return
	}

	coverArtArchiveUrl, err := musicbrainz.GetAlbumArtUrl(ctx, album.MusicBrainzId)

	response := AlbumArtsResponse{
		Deezer:          deezerImageUrl,
		CoverArtArchive: coverArtArchiveUrl,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
