package handlers

import (
	"net/http"
	"strconv"
	"time"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleCreateShare(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	description := form["description"]
	expiresUnixTime := form["expires"]

	_, mediaIds, err := net.ParseDuplicateFormKeys(r, "id", false)
	if err != nil {
		logger.Printf("Error parsing media Ids: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid media Ids", "")
		return
	}

	if len(mediaIds) == 0 {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	ctx := r.Context()

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to create shares", "")
		return
	}

	if !requestUser.ShareRole {
		logger.Printf("User %s attempted to create a share without share role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to create shares", "")
		return
	}

	var expiresUnixTimeInt int
	var expiresTime time.Time
	if expiresUnixTime != "" {
		expiresUnixTimeInt, err = strconv.Atoi(expiresUnixTime)
		if err != nil {
			logger.Printf("Error converting expiresUnixTime to int: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid expires parameter", "")
			return
		}
		expiresTime = logic.GetTimeFromUnixTimestamp(expiresUnixTimeInt)
	}

	createShareOptions := database.CreateShareOptions{
		Description: description,
		ExpiresAt:   expiresTime,
		MediaIds:    mediaIds,
	}
	logger.Printf("%v", createShareOptions)

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	response.SubsonicResponse.Shares = &types.Shares{}

	net.WriteSubsonicResponse(w, r, response, format)
}
