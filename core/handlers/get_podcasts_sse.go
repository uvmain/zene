package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetPodcastsServerSentEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	podcastId := form["id"]

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

	var podcastIdInt int
	if podcastId != "" {
		podcastIdInt, err = strconv.Atoi(podcastId)
		if err != nil {
			logger.Printf("Error converting podcast id to int: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter should be an integer", "")
			return
		}
	}

	// every 2 seconds, getPodcasts and write to w
	clientGone := r.Context().Done()

	rc := http.NewResponseController(w)
	t := time.NewTicker(time.Second * 2)
	defer t.Stop()
	for {
		select {
		case <-clientGone:
			return
		case <-t.C:
			podcasts, err := getPodcasts(ctx, podcastIdInt)
			if err != nil {
				logger.Printf("Error getting podcasts in HandleGetPodcastsServerSentEvents: %v", err)
				return
			}

			jsonData, err := json.Marshal(podcasts)
			if err != nil {
				logger.Printf("Error marshalling podcasts to JSON: %v", err)
				return
			}

			_, err = fmt.Fprintf(w, "data: %s\n\n", jsonData)
			if err != nil {
				logger.Printf("Error writing to response in HandleGetPodcastsServerSentEvents: %v", err)
				return
			}

			err = rc.Flush()
			if err != nil {
				logger.Printf("Error flushing response in HandleGetPodcastsServerSentEvents: %v", err)
				return
			}
		}
	}

}

func getPodcasts(ctx context.Context, podcastId int) ([]types.PodcastChannel, error) {
	podcasts, err := database.GetPodcasts(ctx, podcastId, true)
	if err != nil {
		return nil, err
	}
	return podcasts, nil
}
