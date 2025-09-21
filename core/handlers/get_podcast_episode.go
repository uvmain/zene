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

func HandleGetPodcastEpisode(w http.ResponseWriter, r *http.Request) {
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
		logger.Printf("User %s attempted to delete a podcast channel without podcast role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to delete podcast channels", "")
		return
	}

	if episodeId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Parameter 'id' is required", "")
		return
	}

	var episodeIdInt int
	if episodeId != "" {
		episodeIdInt, err = strconv.Atoi(episodeId)
		if err != nil {
			logger.Printf("Error converting podcast episode id to int: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter should be an integer", "")
			return
		}
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	episode, err := database.GetPodcastEpisodeById(ctx, episodeIdInt)
	if err != nil {
		logger.Printf("Error querying database in GetPodcasts: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to query database", "")
		return
	}

	response.SubsonicResponse.PodcastEpisode = &episode

	net.WriteSubsonicResponse(w, r, response, format)
}
