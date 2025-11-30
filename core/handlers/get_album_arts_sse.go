package handlers

import (
	"encoding/json"
	"fmt"
	"sync"

	"net/http"
	"zene/core/art"
	"zene/core/database"
	"zene/core/deezer"
	"zene/core/musicbrainz"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetAlbumArtsServerSentEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	ctx := r.Context()

	form := net.NormalisedForm(r, w)
	artistName := form["artist"]
	albumName := form["album"]

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

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	deezerFunc := func() types.SseMessage {
		deezerImageUrl, _ := deezer.GetAlbumArtUrlWithArtistNameAndAlbumName(ctx, artistName, albumName)
		return types.SseMessage{Source: "Deezer", Data: deezerImageUrl}
	}
	coverArtArchiveFunc := func() types.SseMessage {
		coverArtArchiveUrl, _ := musicbrainz.GetAlbumArtUrl(ctx, album.MusicBrainzId)
		return types.SseMessage{Source: "CoverArtArchive", Data: coverArtArchiveUrl}
	}
	localArtFunc := func() types.SseMessage {
		localArts, _ := art.GetLocalArtAsBase64(ctx, album.MusicBrainzId)
		return types.SseMessage{Source: "LocalArt", Data: localArts}
	}

	results := make(chan types.SseMessage, 3)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() { defer wg.Done(); results <- deezerFunc() }()
	go func() { defer wg.Done(); results <- coverArtArchiveFunc() }()
	go func() { defer wg.Done(); results <- localArtFunc() }()

	go func() {
		wg.Wait()
		close(results)
	}()

	for msg := range results {
		jsonBytes, _ := json.Marshal(msg)
		fmt.Fprintf(w, "data: %s\n\n", jsonBytes)
		flusher.Flush()
	}

	fmt.Fprintf(w, "event: done\ndata: complete\n\n")
	flusher.Flush()
}
