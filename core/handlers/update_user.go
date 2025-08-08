package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/encryption"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "User not found", "")
		return
	}

	if requestUser.AdminRole == false {
		logger.Printf("User %s attempted to create a user without admin role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to create users", "")
		return
	}

	username := r.FormValue("username")
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

	password := r.FormValue("password")
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

	email := r.FormValue("email")
	if email != "" {
		userToUpdate.Email = email
	}

	ldapAuthenticated := r.FormValue("ldapAuthenticated")
	if ldapAuthenticated != "" {
		userToUpdate.LdapAuthenticated = net.ParseBooleanFormValue(w, r, ldapAuthenticated)
	}

	adminRole := r.FormValue("adminRole")
	if adminRole != "" {
		userToUpdate.AdminRole = net.ParseBooleanFormValue(w, r, adminRole)
	}

	settingsRole := r.FormValue("settingsRole")
	if settingsRole != "" {
		userToUpdate.SettingsRole = net.ParseBooleanFormValue(w, r, settingsRole)
	}

	streamRole := r.FormValue("streamRole")
	if streamRole != "" {
		userToUpdate.StreamRole = net.ParseBooleanFormValue(w, r, streamRole)
	}

	jukeboxRole := r.FormValue("jukeboxRole")
	if jukeboxRole != "" {
		userToUpdate.JukeboxRole = net.ParseBooleanFormValue(w, r, jukeboxRole)
	}

	downloadRole := r.FormValue("downloadRole")
	if downloadRole != "" {
		userToUpdate.DownloadRole = net.ParseBooleanFormValue(w, r, downloadRole)
	}

	uploadRole := r.FormValue("uploadRole")
	if uploadRole != "" {
		userToUpdate.UploadRole = net.ParseBooleanFormValue(w, r, uploadRole)
	}

	playlistRole := r.FormValue("playlistRole")
	if playlistRole != "" {
		userToUpdate.PlaylistRole = net.ParseBooleanFormValue(w, r, playlistRole)
	}

	coverArtRole := r.FormValue("coverArtRole")
	if coverArtRole != "" {
		userToUpdate.CoverArtRole = net.ParseBooleanFormValue(w, r, coverArtRole)
	}

	commentRole := r.FormValue("commentRole")
	if commentRole != "" {
		userToUpdate.CommentRole = net.ParseBooleanFormValue(w, r, commentRole)
	}

	podcastRole := r.FormValue("podcastRole")
	if podcastRole != "" {
		userToUpdate.PodcastRole = net.ParseBooleanFormValue(w, r, podcastRole)
	}

	shareRole := r.FormValue("shareRole")
	if shareRole != "" {
		userToUpdate.ShareRole = net.ParseBooleanFormValue(w, r, shareRole)
	}

	scrobblingEnabled := r.FormValue("scrobblingEnabled")
	if scrobblingEnabled != "" {
		userToUpdate.ScrobblingEnabled = net.ParseBooleanFormValue(w, r, scrobblingEnabled)
	}

	videoConversionRole := r.FormValue("videoConversionRole")
	if videoConversionRole != "" {
		userToUpdate.VideoConversionRole = net.ParseBooleanFormValue(w, r, videoConversionRole)
	}

	maxBitRate := r.FormValue("maxBitRate")
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

	musicFolderId := r.FormValue("musicFolderId")
	if musicFolderId != "" {
		folderIdInts, _, err := net.ParseDuplicateFormKeys(r, "musicFolderId", true)
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
	response := types.GetPopulatedSubsonicResponse(false)

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
