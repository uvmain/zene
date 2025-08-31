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

func HandleGetApiKeys(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	userId := form["user_id"]

	ctx := r.Context()

	var userIdInt int
	var err error

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
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "user not authorized to create API key for this user", "")
		return
	}

	if userIdInt == 0 {
		userIdInt = requestUser.Id
	}

	apiKeys, err := database.GetApiKeys(ctx, userIdInt)
	if err != nil {
		logger.Printf("Error getting API keys for user ID %d: %v", userIdInt, err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "User not found", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	response.SubsonicResponse.ApiKeys = &types.ApiKeys{}
	response.SubsonicResponse.ApiKeys.ApiKeys = apiKeys

	net.WriteSubsonicResponse(w, r, response, format)
}
