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

func HandleGetNewestPodcasts(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	count := form["count"]

	ctx := r.Context()

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "Error fetching user from context", "")
		return
	}

	if !requestUser.PodcastRole {
		logger.Printf("User %s attempted to delete a podcast channel without podcast role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to delete podcast channels", "")
		return
	}

	var countInt int
	if count != "" {
		countInt, err = strconv.Atoi(count)
		if err != nil {
			logger.Printf("Error converting count to int: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "count parameter should be an integer", "")
			return
		}
	} else {
		countInt = 20
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	episodes, err := database.GetNewestPodcastEpisodes(ctx, countInt)
	if err != nil {
		logger.Printf("Error querying database in GetNewestPodcastEpisodes: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to query database", "")
		return
	}
	if episodes == nil {
		episodes = []types.PodcastEpisode{}
	}

	response.SubsonicResponse.NewestPodcasts = &types.NewestPodcasts{
		Episodes: episodes,
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
