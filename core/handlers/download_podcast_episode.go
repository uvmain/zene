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

func HandleDownloadPodcastEpisode(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	episodeId := form["id"]

	ctx := r.Context()

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "Error fetching user from context", "")
		return
	}

	if !requestUser.PodcastRole {
		logger.Printf("User %s attempted to trigger a podcast episode download without podcast role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to download podcast episodes", "")
		return
	}

	if episodeId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	err = podcasts.DownloadPodcastEpisode(ctx, episodeId)
	if err != nil {
		if err.Error() == "podcast episode is already downloading" {
			net.WriteSubsonicError(w, r, types.ErrorGeneric, "Podcast episode is already downloading", "")
			return
		}
		logger.Printf("Error downloading podcast episode: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error downloading podcast episode", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
