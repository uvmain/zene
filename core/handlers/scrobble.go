package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
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
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()

	user, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to scrobble", "")
		return
	}

	var metadataIds []string
	metadataId := r.FormValue("id")
	if metadataId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id is required", "")
		return
	} else {
		_, metadataIds, err = net.ParseDuplicateFormKeys(r, "id", true)
		if err != nil {
			logger.Printf("Error parsing id: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid id", "")
			return
		}
	}

	var timeInt int
	timeInMilliseconds := r.FormValue("time")
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
	submission := r.FormValue("submission")
	if submission != "" {
		submissionBool = net.ParseBooleanFromString(w, r, submission)
	}

	playerName := r.FormValue("c")

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

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	format := r.FormValue("f")
	if format == "json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))
		xml.NewEncoder(w).Encode(response.SubsonicResponse)
	}
}
