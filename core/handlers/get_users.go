package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]

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

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

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

	net.WriteSubsonicResponse(w, r, response, format)
}
