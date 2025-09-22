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

func HandleDeletePodcastChannel(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	channelId := form["id"]

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

	if channelId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	channelIdInt, err := strconv.Atoi(channelId)
	if err != nil {
		logger.Printf("Error converting channel ID to int: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid id parameter, it should be an integer", "")
		return
	}

	err = podcasts.DeletePodcastChannelAndEpisodes(ctx, channelIdInt)
	if err != nil {
		logger.Printf("Error deleting podcast channel and episodes %d: %v", channelIdInt, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to delete podcast channel and episodes", "")
		return
	}

	logger.Printf("Podcast channel %d and episodes deleted by %s", channelIdInt, requestUser.Username)

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
