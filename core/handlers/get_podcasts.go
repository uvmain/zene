package handlers

import (
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetPodcasts(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	includeEpisodes := form["includeepisodes"]
	podcastId := form["id"]

	if includeEpisodes != "true" && includeEpisodes != "false" && includeEpisodes != "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Parameter 'includeEpisodes' must be 'true' or 'false'", "")
		return
	}

	if includeEpisodes == "" {
		includeEpisodes = "true"
	}

	var podcastIdInt int
	var err error
	if podcastId != "" {
		podcastIdInt, err = strconv.Atoi(podcastId)
		if err != nil {
			logger.Printf("Error converting podcast id to int: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter should be an integer", "")
			return
		}
	}

	includeEpisodesBool := includeEpisodes == "true"

	ctx := r.Context()

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	podcasts, err := database.GetPodcasts(ctx, podcastIdInt, includeEpisodesBool)
	if err != nil {
		logger.Printf("Error querying database in GetPodcasts: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to query database", "")
		return
	}
	if podcasts == nil {
		podcasts = []types.PodcastChannel{}
	}

	response.SubsonicResponse.PodcastChannels = &types.PodcastChannels{}
	response.SubsonicResponse.PodcastChannels.PodcastChannels = podcasts

	net.WriteSubsonicResponse(w, r, response, format)
}
