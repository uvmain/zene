package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/encryption"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	username := form["username"]
	password := form["password"]

	ctx := r.Context()

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to change passwords", "")
		return
	}

	if requestUser.AdminRole == false {
		logger.Printf("User %s attempted to create a user without admin role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to change passwords", "")
		return
	}

	if username == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Username is required", "")
		return
	}

	if password == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Password is required", "")
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

	_, err = database.UpsertUser(ctx, userToUpdate)
	if err != nil {
		logger.Printf("Error updating user %s: %v", username, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to update password for user", "")
		return
	}

	logger.Printf("Password for user %s updated successfully by %s", username, requestUser.Username)
	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	net.WriteSubsonicResponse(w, r, response, format)
}
