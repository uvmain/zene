package handlers

import (
	"net/http"
	"zene/core/database"
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

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "Error fetching user from context", "")
		return
	}

	if !requestUser.PodcastRole {
		logger.Printf("User %s attempted to create a podcast channel without podcast role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to create podcast channels", "")
		return
	}

	if url == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "url parameter is required", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	err = podcasts.CreateNewPodcastFromFeedUrl(ctx, url)
	if err != nil {
		logger.Printf("Error creating podcast: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error creating podcast", "")
		return
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
