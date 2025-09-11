package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetPlayqueueByIndex(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]

	ctx := r.Context()

	playqueue, err := database.GetPlayqueue(ctx)
	if err != nil {
		logger.Printf("Error getting playqueue: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to get playqueue", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	playqueueResponse := &types.PlayqueueByIndex{}
	playqueueResponse.Changed = playqueue.Changed
	playqueueResponse.ChangedBy = playqueue.ChangedBy
	playqueueResponse.CurrentIndex = playqueue.CurrentIndex
	playqueueResponse.Position = playqueue.Position
	playqueueResponse.Username = playqueue.Username

	children, err := database.GetSongsByIDs(ctx, playqueue.TrackIds)
	if err != nil {
		logger.Printf("Error getting playqueue children: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to get playqueue children", "")
	}

	playqueueResponse.Entry = children

	response.SubsonicResponse.PlayQueueByIndex = playqueueResponse

	net.WriteSubsonicResponse(w, r, response, format)
}
