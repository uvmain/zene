package handlers

import (
	"net/http"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/podcasts"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleCreatePodcastChannel(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	url := form["url"]

	ctx := r.Context()

	if url == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "url parameter is required", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	err := podcasts.CreateNewPodcastFromFeedUrl(ctx, url)
	if err != nil {
		logger.Printf("Error creating podcast: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error creating podcast", "")
		return
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
