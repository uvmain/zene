package handlers

import (
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleDeleteApiKey(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	userId := form["user_id"]
	apiKeyId := form["id"]

	ctx := r.Context()

	var userIdInt int
	var err error

	if apiKeyId == "" {
		logger.Printf("Error deleting API key: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	apiKeyIdInt, err := strconv.Atoi(apiKeyId)
	if err != nil {
		logger.Printf("Error converting apiKeyId to int: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter should be an integer", "")
		return
	}

	if userId != "" {
		userIdInt, err = strconv.Atoi(userId)
		if err != nil {
			logger.Printf("Error converting user_id to int: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "user_id parameter should be an integer", "")
			return
		}
	}

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "context user not found", "")
		return
	}

	if requestUser.AdminRole == false && requestUser.Id != userIdInt {
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "user not authorized to delete API key for other users", "")
		return
	}

	if userIdInt == 0 {
		userIdInt = requestUser.Id
	}

	err = database.DeleteApiKey(ctx, apiKeyIdInt, userIdInt)
	if err != nil {
		logger.Printf("Error deleting API key %d: %v", apiKeyIdInt, err)
		net.WriteSubsonicError(w, r, types.ErrorInvalidApiKey, "Server Error", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
