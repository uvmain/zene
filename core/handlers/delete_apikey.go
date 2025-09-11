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
	userId := form["userid"]
	apiKeyId := form["id"]

	ctx := r.Context()

	var userIdInt int
	var err error

	if apiKeyId == "" {
		logger.Printf("Error deleting API key: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	apiKeyIds, _, err := net.ParseDuplicateFormKeys(r, "id", true)
	if err != nil {
		logger.Printf("Error parsing id parameter(s): %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid id parameter(s) received", "")
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

	if !requestUser.AdminRole && requestUser.Id != userIdInt {
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "user not authorized to delete API key for other users", "")
		return
	}

	if userIdInt == 0 {
		userIdInt = requestUser.Id
	}

	err = database.DeleteApiKeys(ctx, apiKeyIds, userIdInt)
	if err != nil {
		logger.Printf("Error deleting API key(s) %v: %v", apiKeyIds, err)
		net.WriteSubsonicError(w, r, types.ErrorInvalidApiKey, "Server Error", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
