package handlers

import (
	"encoding/json"

	"net/http"
	"zene/core/art"
	"zene/core/database"
	"zene/core/deezer"
	"zene/core/logger"
	"zene/core/musicbrainz"
	"zene/core/net"
	"zene/core/types"
)

type ArtistArtsResponse struct {
	Deezer           string `json:"deezer"`
	CoverArtArchive  string `json:"cover_art_archive"`
	LocalFolderArt   string `json:"local_folder_art"`
	LocalEmbeddedArt string `json:"local_embedded_art"`
}

func HandleGetArtistArts(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	artistName := form["artist"]
	artistId := form["id"]
	ctx := r.Context()

	if artistName == "" && artistId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "artist or id parameters are required", "")
		return
	}

	var err error
	var artist types.Artist

	if artistId == "" {
		artistId, err = database.GetArtistIdByName(ctx, artistName)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "artist not found", "")
			return
		}
	}

	artist, err = database.SelectArtistByMusicBrainzArtistId(ctx, artistId)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "artist not found", "")
		return
	}

	logger.Printf("Getting arts for artist: %s (ID: %s)", artist.Name, artist.MusicBrainzId)

	deezerImageUrl, _ := deezer.GetArtistArtUrlWithArtistName(ctx, artist.Name)
	coverArtArchiveUrl, _ := musicbrainz.GetArtistArtUrl(ctx, artist.MusicBrainzId)
	localArts, _ := art.GetLocalArtistArtAsBase64(ctx, artist.MusicBrainzId)

	response := ArtistArtsResponse{
		Deezer:          deezerImageUrl,
		CoverArtArchive: coverArtArchiveUrl,
		LocalFolderArt:  localArts.FolderArt,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
