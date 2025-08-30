package handlers

import (
	"html"
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleAddChatMessage(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	message := form["message"]
	format := form["f"]

	ctx := r.Context()

	user, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to add chat messages", "")
		return
	}

	if message == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Message is required", "")
		return
	}

	// sanitize message
	message = html.EscapeString(message)

	err = database.InsertChat(ctx, user.Id, message)
	if err != nil {
		logger.Printf("Error inserting chat message for user %d: %v", user.Id, err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to add chat message", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
