package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to get users", "")
		return
	}

	if requestUser.AdminRole == false {
		logger.Printf("User %s attempted to create a user without admin role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to get users", "")
		return
	}

	response := types.SubsonicUsersResponseWrapper{}
	stdRes := types.GetPopulatedSubsonicResponse(false)

	response.SubsonicResponse.XMLName = stdRes.SubsonicResponse.XMLName
	response.SubsonicResponse.Xmlns = stdRes.SubsonicResponse.Xmlns
	response.SubsonicResponse.Status = stdRes.SubsonicResponse.Status
	response.SubsonicResponse.Version = stdRes.SubsonicResponse.Version
	response.SubsonicResponse.Type = stdRes.SubsonicResponse.Type
	response.SubsonicResponse.ServerVersion = stdRes.SubsonicResponse.ServerVersion
	response.SubsonicResponse.OpenSubsonic = stdRes.SubsonicResponse.OpenSubsonic
	response.SubsonicResponse.Users = &types.SubsonicUsers{}
	response.SubsonicResponse.Users.User = make([]types.SubsonicUser, 0)

	allUsers, err := database.GetAllUsers(ctx)
	if err != nil {
		logger.Printf("Error getting all users: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Users not found", "")
		return
	}

	for _, user := range allUsers {
		response.SubsonicResponse.Users.User = append(response.SubsonicResponse.Users.User, types.SubsonicUser{
			Username:            user.Username,
			Email:               user.Email,
			ScrobblingEnabled:   user.ScrobblingEnabled,
			AdminRole:           user.AdminRole,
			SettingsRole:        user.SettingsRole,
			StreamRole:          user.StreamRole,
			JukeboxRole:         user.JukeboxRole,
			DownloadRole:        user.DownloadRole,
			UploadRole:          user.UploadRole,
			PlaylistRole:        user.PlaylistRole,
			CoverArtRole:        user.CoverArtRole,
			CommentRole:         user.CommentRole,
			PodcastRole:         user.PodcastRole,
			ShareRole:           user.ShareRole,
			VideoConversionRole: user.VideoConversionRole,
			MaxBitRate:          user.MaxBitRate,
			Folders:             user.Folders,
		})
	}

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
