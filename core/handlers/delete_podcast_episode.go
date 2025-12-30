package handlers

import (
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/podcasts"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleDeletePodcastEpisode(w http.ResponseWriter, r *http.Request) {
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
		logger.Printf("User %s attempted to delete a podcast episode without podcast role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to delete podcast episodes", "")
		return
	}

	if episodeId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	episodeIdInt, err := strconv.Atoi(episodeId)
	if err != nil {
		logger.Printf("Error converting podcast episode id to int: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter should be an integer", "")
		return
	}

	err = podcasts.DeletePodcastEpisodeById(ctx, episodeIdInt, false)
	if err != nil {
		logger.Printf("Error deleting podcast episode %d: %v", episodeIdInt, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to delete podcast episode", "")
		return
	}

	logger.Printf("Podcast episode %s deleted by %s", episodeId, requestUser.Username)

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
