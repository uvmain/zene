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

func HandleGetArtistArtsServerSentEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	ctx := r.Context()

	form := net.NormalisedForm(r, w)
	artistName := form["artist"]
	artistId := form["id"]

	if artistName == "" && artistId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "artist or id parameters are required", "")
		return
	}

	var err error

	if artistId == "" {
		artistId, err = database.GetArtistIdByName(ctx, artistName)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "artist not found", "")
			return
		}
	}

	artist, err := database.SelectArtistByMusicBrainzArtistId(ctx, artistId)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "artist not found", "")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	deezerFunc := func() types.SseMessage {
		deezerImageUrl, _ := deezer.GetArtistArtUrlWithArtistName(ctx, artist.Name)
		return types.SseMessage{Source: "Deezer", Data: deezerImageUrl}
	}
	coverArtArchiveFunc := func() types.SseMessage {
		coverArtArchiveUrl, _ := musicbrainz.GetArtistArtUrl(ctx, artist.MusicBrainzId)
		return types.SseMessage{Source: "CoverArtArchive", Data: coverArtArchiveUrl}
	}
	localArtFunc := func() types.SseMessage {
		localArts, _ := art.GetLocalArtistArtAsBase64(ctx, artist.MusicBrainzId)
		return types.SseMessage{Source: "LocalArt", Data: localArts}
	}

	results := make(chan types.SseMessage, 3)
	clientGone := ctx.Done()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() { defer wg.Done(); results <- deezerFunc() }()
	go func() { defer wg.Done(); results <- coverArtArchiveFunc() }()
	go func() { defer wg.Done(); results <- localArtFunc() }()

	go func() {
		wg.Wait()
		close(results)
	}()

	for {
		select {
		case <-clientGone:
			return // Client disconnected, stop processing
		case msg, ok := <-results:
			if !ok {
				// Channel closed, all results received
				// Only send completion event if client is still connected
				select {
				case <-clientGone:
					return
				default:
					fmt.Fprintf(w, "event: done\ndata: complete\n\n")
					flusher.Flush()
				}
				return
			}
			jsonBytes, _ := json.Marshal(msg)
			fmt.Fprintf(w, "data: %s\n\n", jsonBytes)
			flusher.Flush()
		}
	}
}
