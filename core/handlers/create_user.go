package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/encryption"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
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
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to create users", "")
		return
	}

	if !requestUser.AdminRole {
		logger.Printf("User %s attempted to create a user without admin role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to create users", "")
		return
	}

	userToCreate := types.User{}

	if username == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Username is required", "")
		return
	}

	usernameAlreadyExists, _ := database.GetUserByUsername(ctx, username)
	if usernameAlreadyExists.Id > 0 {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "User already exists", "")
		return
	}

	userToCreate.Username = username

	if password == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Password is required", "")
		return
	}

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
		logger.Printf("Error decrypting hex encoded password for user %s: %v", username, err)
		net.WriteSubsonicError(w, r, types.ErrorWrongCredentials, "Wrong username or password", "")
	}

	userToCreate.Password = encryptedPassword

	if email == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Email is required", "")
		return
	}
	userToCreate.Email = email

	if ldapAuthenticated != "" {
		userToCreate.LdapAuthenticated = net.ParseBooleanFromString(w, r, ldapAuthenticated)
	} else {
		userToCreate.LdapAuthenticated = logic.GetDefaultRoleValue("ldapAuthenticated")
	}

	if adminRole != "" {
		userToCreate.AdminRole = net.ParseBooleanFromString(w, r, adminRole)
	} else {
		userToCreate.AdminRole = logic.GetDefaultRoleValue("adminRole")
	}

	if settingsRole != "" {
		userToCreate.SettingsRole = net.ParseBooleanFromString(w, r, settingsRole)
	} else {
		userToCreate.SettingsRole = logic.GetDefaultRoleValue("settingsRole")
	}

	if streamRole != "" {
		userToCreate.StreamRole = net.ParseBooleanFromString(w, r, streamRole)
	} else {
		userToCreate.StreamRole = logic.GetDefaultRoleValue("streamRole")
	}

	if jukeboxRole != "" {
		userToCreate.JukeboxRole = net.ParseBooleanFromString(w, r, jukeboxRole)
	} else {
		userToCreate.JukeboxRole = logic.GetDefaultRoleValue("jukeboxRole")
	}

	if downloadRole != "" {
		userToCreate.DownloadRole = net.ParseBooleanFromString(w, r, downloadRole)
	} else {
		userToCreate.DownloadRole = logic.GetDefaultRoleValue("downloadRole")
	}

	if uploadRole != "" {
		userToCreate.UploadRole = net.ParseBooleanFromString(w, r, uploadRole)
	} else {
		userToCreate.UploadRole = logic.GetDefaultRoleValue("uploadRole")
	}

	if playlistRole != "" {
		userToCreate.PlaylistRole = net.ParseBooleanFromString(w, r, playlistRole)
	} else {
		userToCreate.PlaylistRole = logic.GetDefaultRoleValue("playlistRole")
	}

	if coverArtRole != "" {
		userToCreate.CoverArtRole = net.ParseBooleanFromString(w, r, coverArtRole)
	} else {
		userToCreate.CoverArtRole = logic.GetDefaultRoleValue("coverArtRole")
	}

	if commentRole != "" {
		userToCreate.CommentRole = net.ParseBooleanFromString(w, r, commentRole)
	} else {
		userToCreate.CommentRole = logic.GetDefaultRoleValue("commentRole")
	}

	if podcastRole != "" {
		userToCreate.PodcastRole = net.ParseBooleanFromString(w, r, podcastRole)
	} else {
		userToCreate.PodcastRole = logic.GetDefaultRoleValue("podcastRole")
	}

	if shareRole != "" {
		userToCreate.ShareRole = net.ParseBooleanFromString(w, r, shareRole)
	} else {
		userToCreate.ShareRole = logic.GetDefaultRoleValue("shareRole")
	}

	if scrobblingEnabled != "" {
		userToCreate.ScrobblingEnabled = net.ParseBooleanFromString(w, r, scrobblingEnabled)
	} else {
		userToCreate.ScrobblingEnabled = logic.GetDefaultRoleValue("scrobblingEnabled")
	}

	if videoConversionRole != "" {
		logger.Printf("Setting videoConversionRole for user %s to %s", username, videoConversionRole)
		userToCreate.VideoConversionRole = net.ParseBooleanFromString(w, r, videoConversionRole)
	} else {
		userToCreate.VideoConversionRole = logic.GetDefaultRoleValue("videoConversionRole")
	}

	if maxBitRate != "" {
		maxBitRateInt, err := strconv.Atoi(maxBitRate)
		if err != nil {
			logger.Printf("Error parsing maxBitRate: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid maxBitRate", "")
			return
		}
		userToCreate.MaxBitRate = maxBitRateInt
	} else {
		userToCreate.MaxBitRate = 0
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
		userToCreate.Folders = folderIdInts
	} else {
		allMusicFolders, err := database.GetMusicFolders(ctx)
		if err != nil {
			logger.Printf("Error getting music folders: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Music folders not found", "")
			return
		}
		for _, folder := range allMusicFolders {
			userToCreate.Folders = append(userToCreate.Folders, folder.Id)
		}
	}

	userId, err := database.UpsertUser(ctx, userToCreate)
	if err != nil {
		logger.Printf("Error creating user %s: %v", username, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to create user", "")
		return
	}

	logger.Printf("User %s created with ID %d by %s", username, userId, ctx.Value("username"))

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
