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

func HandleUpdateShare(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	id := form["id"]
	description := form["description"]
	expiresUnixTime := form["expires"]

	ctx := r.Context()

	idInt, err := strconv.Atoi(id)
	if err != nil {
		logger.Printf("Error converting share id to int: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid share id parameter", "")
		return
	}

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to update shares", "")
		return
	}

	if !requestUser.ShareRole {
		logger.Printf("User %s attempted to update a share without share role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to update shares", "")
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

	updateShareOptions := database.UpdateShareOptions{
		ShareId:           idInt,
		UpdateDescription: description != "",
		Description:       description,
		UpdateExpiresAt:   !expiresTime.IsZero(),
		ExpiresAt:         expiresTime,
	}

	err = database.UpdateShare(ctx, updateShareOptions)
	if err != nil {
		logger.Printf("Error updating share: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to update share", "")
		return
	}

	logger.Printf("User %s updated share with ID %d", requestUser.Username, idInt)

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
