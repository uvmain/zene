package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	username := form["username"]

	ctx := r.Context()

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	var user types.User
	var err error

	requestUser, _ := database.GetUserByContext(ctx)
	if requestUser.AdminRole && username != "" {
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
	response.SubsonicResponse.User.MaxBitRate = user.MaxBitRate
	response.SubsonicResponse.User.Folders = user.Folders

	net.WriteSubsonicResponse(w, r, response, format)
}
