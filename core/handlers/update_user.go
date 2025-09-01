package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/encryption"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	username := form["username"]
	password := form["password"]
	email := form["email"]
	ldapAuthenticated := form["ldapauthenticated"]
	adminRole := form["adminrole"]
	settingsRole := form["settingsrole"]
	streamRole := form["streamrole"]
	jukeboxRole := form["jukeboxrole"]
	downloadRole := form["downloadrole"]
	uploadRole := form["uploadrole"]
	playlistRole := form["playlistrole"]
	coverArtRole := form["coverartrole"]
	commentRole := form["commentrole"]
	podcastRole := form["podcastrole"]
	shareRole := form["sharerole"]
	scrobblingEnabled := form["scrobblingenabled"]
	videoConversionRole := form["videoconversionrole"]
	maxBitRate := form["maxbitrate"]
	musicFolderId := form["musicfolderid"]

	ctx := r.Context()

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to update users", "")
		return
	}

	if requestUser.AdminRole == false {
		logger.Printf("User %s attempted to update a user without admin role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to update users", "")
		return
	}

	if username == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Username is required", "")
		return
	}

	userToUpdate, err := database.GetUserByUsername(ctx, username)
	if err != nil {
		logger.Printf("Error getting user by username %s: %v", username, err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "User not found", "")
		return
	}

	if password != "" {
		if len(password) > 4 && password[:4] == "enc:" {
			decryptedPassword, err := encryption.HexDecrypt(password[4:])
			if err != nil {
				logger.Printf("Error decrypting hex encoded password for user %s: %v", username, err)
				net.WriteSubsonicError(w, r, types.ErrorWrongCredentials, "Wrong username or password", "")
				return
			}
			password = decryptedPassword
		}
		encryptedPassword, err := encryption.EncryptAES(password)
		if err != nil {
			logger.Printf("Error encrypting password for user %s: %v", username, err)
			net.WriteSubsonicError(w, r, types.ErrorWrongCredentials, "Wrong username or password", "")
			return
		}
		userToUpdate.Password = encryptedPassword
	}

	if email != "" {
		userToUpdate.Email = email
	}

	if ldapAuthenticated != "" {
		userToUpdate.LdapAuthenticated = net.ParseBooleanFromString(w, r, ldapAuthenticated)
	}

	if adminRole != "" {
		userToUpdate.AdminRole = net.ParseBooleanFromString(w, r, adminRole)
	}

	if settingsRole != "" {
		userToUpdate.SettingsRole = net.ParseBooleanFromString(w, r, settingsRole)
	}

	if streamRole != "" {
		userToUpdate.StreamRole = net.ParseBooleanFromString(w, r, streamRole)
	}

	if jukeboxRole != "" {
		userToUpdate.JukeboxRole = net.ParseBooleanFromString(w, r, jukeboxRole)
	}

	if downloadRole != "" {
		userToUpdate.DownloadRole = net.ParseBooleanFromString(w, r, downloadRole)
	}

	if uploadRole != "" {
		userToUpdate.UploadRole = net.ParseBooleanFromString(w, r, uploadRole)
	}

	if playlistRole != "" {
		userToUpdate.PlaylistRole = net.ParseBooleanFromString(w, r, playlistRole)
	}

	if coverArtRole != "" {
		userToUpdate.CoverArtRole = net.ParseBooleanFromString(w, r, coverArtRole)
	}

	if commentRole != "" {
		userToUpdate.CommentRole = net.ParseBooleanFromString(w, r, commentRole)
	}

	if podcastRole != "" {
		userToUpdate.PodcastRole = net.ParseBooleanFromString(w, r, podcastRole)
	}

	if shareRole != "" {
		userToUpdate.ShareRole = net.ParseBooleanFromString(w, r, shareRole)
	}

	if scrobblingEnabled != "" {
		userToUpdate.ScrobblingEnabled = net.ParseBooleanFromString(w, r, scrobblingEnabled)
	}

	if videoConversionRole != "" {
		userToUpdate.VideoConversionRole = net.ParseBooleanFromString(w, r, videoConversionRole)
	}

	if maxBitRate != "" {
		maxBitRateInt, err := strconv.Atoi(maxBitRate)
		if err != nil {
			logger.Printf("Error parsing maxBitRate: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid maxBitRate", "")
			return
		}
		userToUpdate.MaxBitRate = maxBitRateInt
	} else {
		userToUpdate.MaxBitRate = 0
	}

	if musicFolderId != "" {
		folderIdInts, _, err := net.ParseDuplicateFormKeys(r, "musicfolderid", true)
		if err != nil {
			logger.Printf("Error parsing musicFolderId: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid musicFolderId", "")
			return
		}
		if len(folderIdInts) == 0 {
			logger.Printf("Error: musicFolderId must contain at least one folder ID")
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "At least one music folder ID is required", "")
			return
		} else {
			for _, folderId := range folderIdInts {
				_, err := database.GetMusicFolderById(ctx, folderId)
				if err != nil {
					logger.Printf("Error checking music folder ID %d: %v", folderId, err)
					net.WriteSubsonicError(w, r, types.ErrorDataNotFound, fmt.Sprintf("Music folder ID %d not found", folderId), "")
					return
				}
			}
		}
		userToUpdate.Folders = folderIdInts
	}

	userId, err := database.UpsertUser(ctx, userToUpdate)
	if err != nil {
		logger.Printf("Error updating user %s: %v", username, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to update user", "")
		return
	}

	logger.Printf("User %s updated successfully with ID %d", username, userId)
	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
