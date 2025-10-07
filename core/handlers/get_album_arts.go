package handlers

import (
	"encoding/json"

	"net/http"
	"zene/core/art"
	"zene/core/database"
	"zene/core/deezer"
	"zene/core/musicbrainz"
	"zene/core/net"
	"zene/core/types"
)

type AlbumArtsResponse struct {
	Deezer           string `json:"deezer"`
	CoverArtArchive  string `json:"cover_art_archive"`
	LocalFolderArt   string `json:"local_folder_art"`
	LocalEmbeddedArt string `json:"local_embedded_art"`
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

	deezerImageUrl, _ := deezer.GetAlbumArtUrlWithArtistNameAndAlbumName(ctx, artistName, albumName)
	coverArtArchiveUrl, _ := musicbrainz.GetAlbumArtUrl(ctx, album.MusicBrainzId)

	localArts, _ := art.GetLocalArtAsBase64(ctx, album.MusicBrainzId)

	response := AlbumArtsResponse{
		Deezer:           deezerImageUrl,
		CoverArtArchive:  coverArtArchiveUrl,
		LocalFolderArt:   localArts.FolderArt,
		LocalEmbeddedArt: localArts.EmbeddedArt,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
