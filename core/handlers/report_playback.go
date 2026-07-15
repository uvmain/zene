package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleReportPlayback(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	mediaId := form["mediaid"]
	mediaType := form["mediatype"]
	positionMs := form["positionms"]
	state := form["state"]
	playbackRate := form["playbackrate"]
	ignoreScrobble := form["ignorescrobble"]

	ctx := r.Context()

	if mediaId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Media ID is required", "")
		return
	}
	if mediaType == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Media type is required", "")
		return
	}
	if positionMs == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Position (ms) is required", "")
		return
	}
	if state == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "State is required", "")
		return
	}
	if playbackRate == "" {
		playbackRate = "1.0"
	}
	if ignoreScrobble == "" {
		ignoreScrobble = "false"
	}

	user, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to report playback", "")
		return
	}

	err = database.InsertChat(ctx, user.Id, "message")
	if err != nil {
		logger.Printf("Error inserting chat message for user %d: %v", user.Id, err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to add chat message", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
