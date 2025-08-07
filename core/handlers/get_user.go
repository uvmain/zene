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

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()

	response := types.SubsonicUserResponseWrapper{}
	stdRes := types.GetPopulatedSubsonicResponse(false)

	response.SubsonicResponse.XMLName = stdRes.SubsonicResponse.XMLName
	response.SubsonicResponse.Xmlns = stdRes.SubsonicResponse.Xmlns
	response.SubsonicResponse.Status = stdRes.SubsonicResponse.Status
	response.SubsonicResponse.Version = stdRes.SubsonicResponse.Version
	response.SubsonicResponse.Type = stdRes.SubsonicResponse.Type
	response.SubsonicResponse.ServerVersion = stdRes.SubsonicResponse.ServerVersion
	response.SubsonicResponse.OpenSubsonic = stdRes.SubsonicResponse.OpenSubsonic

	username := r.FormValue("username")
	var user types.User
	var err error

	requestUser, _ := database.GetUserByContext(ctx)
	if requestUser.AdminRole == true && username != "" {
		// Admin can request any user
		user, err = database.GetUserByUsername(ctx, username)
		if err != nil {
			logger.Printf("Error getting user by username %s: %v", username, err)
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "User not found", "")
			return
		}
	} else {
		user, err = database.GetUserByContext(ctx)
		if err != nil {
			logger.Printf("Error getting user by username %s: %v", username, err)
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "User not found", "")
			return
		}
	}

	response.SubsonicResponse.User = &types.SubsonicUser{}
	response.SubsonicResponse.User.Username = user.Username
	response.SubsonicResponse.User.Email = user.Email
	response.SubsonicResponse.User.ScrobblingEnabled = user.ScrobblingEnabled
	response.SubsonicResponse.User.AdminRole = user.AdminRole
	response.SubsonicResponse.User.SettingsRole = user.SettingsRole
	response.SubsonicResponse.User.DownloadRole = user.DownloadRole
	response.SubsonicResponse.User.UploadRole = user.UploadRole
	response.SubsonicResponse.User.PlaylistRole = user.PlaylistRole
	response.SubsonicResponse.User.CoverArtRole = user.CoverArtRole
	response.SubsonicResponse.User.CommentRole = user.CommentRole
	response.SubsonicResponse.User.PodcastRole = user.PodcastRole
	response.SubsonicResponse.User.StreamRole = user.StreamRole
	response.SubsonicResponse.User.JukeboxRole = user.JukeboxRole
	response.SubsonicResponse.User.ShareRole = user.ShareRole
	response.SubsonicResponse.User.VideoConversionRole = user.VideoConversionRole
	response.SubsonicResponse.User.Folders = user.Folders

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
