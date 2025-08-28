package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetChatMessages(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	timeSinceParam := form["since"]

	ctx := r.Context()

	var timeSince int
	if timeSinceParam != "" {
		var err error
		timeSince, err = strconv.Atoi(timeSinceParam)
		if err != nil {
			errorString := fmt.Sprintf("Invalid since parameter: %s", timeSinceParam)
			net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
			return
		}
	} else {
		timeSince = 0
	}

	chats, err := database.GetChats(ctx, timeSince)
	if err != nil {
		logger.Printf("Error fetching chats from database: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to fetch chat messages", "")
		return
	}

	logger.Printf("Fetched %d chats since %d", len(chats), timeSince)

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	response.SubsonicResponse.ChatMessages = &types.ChatMessages{
		ChatMessage: chats,
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
