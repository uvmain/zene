package handlers

import (
	"net/http"
	"strconv"
	"time"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleScrobble(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	metadataId := form["id"]
	timeInMilliseconds := form["time"]
	submission := form["submission"]
	playerName := form["c"]

	ctx := r.Context()

	user, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to scrobble", "")
		return
	}

	var metadataIds []string
	if metadataId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id is required", "")
		return
	} else {
		_, metadataIds, err = net.ParseDuplicateFormKeys(r, "id", false)
		if err != nil {
			logger.Printf("Error parsing id: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid id", "")
			return
		}
	}

	var timeInt int
	if timeInMilliseconds == "" {
		timeInt = int(time.Now().UnixMilli())
	} else {
		timeInt, err := strconv.Atoi(timeInMilliseconds)
		if err != nil || timeInt < 0 {
			logger.Printf("Error parsing time for user %d: %v", user.Id, err)
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid time, must be a positive integer", "")
			return
		}
	}

	var submissionBool = true
	if submission != "" {
		submissionBool = net.ParseBooleanFromString(w, r, submission)
	}

	for _, trackId := range metadataIds {
		if submissionBool {
			err = database.UpsertPlayCount(ctx, user.Id, trackId)
		}
		err = database.UpsertNowPlaying(ctx, user.Id, trackId, timeInt, 0, playerName)
		if err != nil {
			logger.Printf("Error upserting now playing for user %d: %v", user.Id, err)
			net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to upsert user now playing", "")
			return
		}
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
