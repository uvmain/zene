package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"zene/core/database"
	"zene/core/encryption"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

// computeToken generates a SHA-256 token from a password and salt.
func computeToken(password, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return hex.EncodeToString(hash[:])
}

// validateToken checks if the provided token matches the computed token from the password and salt.
func validateToken(salt string, token string, encryptedPassword string) (bool, error) {
	decryptedPassword, err := encryption.DecryptAES(encryptedPassword)
	if err != nil {
		return false, err
	}
	expectedToken := computeToken(decryptedPassword, salt)
	return token == expectedToken, nil
}

// ValidateAuth authenticates a user by either plaintext password or token and salt.
// It returns:
// string: the authenticated username
// int64: the authenticated userId
// bool: true if successful; otherwise returns false and writes a 401 response.
//
// This supports the following request parameters in either form data or query parameters:
// - u: username
// - p: plaintext password
// - t: token = sha256(password + salt)
// - s: salt
// - apiKey: API key for authentication
// - v: protocol version
// - c: client name (required for all requests)
//
// If apiKey is specified, then none of p, t, s, nor u can be specified.
// Else either p or both t and s must be specified.
func ValidateAuth(r *http.Request, w http.ResponseWriter) (string, int64, bool) {
	ctx := r.Context()

	username := r.FormValue("u")
	password := r.FormValue("p")
	token := r.FormValue("t")
	salt := r.FormValue("s")
	apiKey := r.FormValue("apiKey")
	version := r.FormValue("v")
	clientName := r.FormValue("c")

	if clientName == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Required parameter 'c' is missing", "")
		return "", 0, false
	}

	if version == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Required parameter 'v' is missing", "")
		return "", 0, false
	}

	if apiKey != "" && (username != "" || password != "" || token != "" || salt != "") {
		net.WriteSubsonicError(w, r, types.ErrorTooManyAuthMechanisms, "Too many authentication mechanisms specified", "")
		return "", 0, false
	}

	if password != "" && (token != "" || salt != "") {
		net.WriteSubsonicError(w, r, types.ErrorTooManyAuthMechanisms, "Too many authentication mechanisms specified", "")
		return "", 0, false
	}

	if apiKey != "" {
		return validateWithApiKey(ctx, apiKey, w, r)
	}

	if username == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Required parameter 'u' is missing", "")
		return "", 0, false
	}

	if password == "" && (token == "" || salt == "") {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Either 'p' or both 't' and 's' parameters are required", "")
		return "", 0, false
	}

	encryptedPassword, userId, err := database.GetEncryptedPasswordFromDB(ctx, username)
	if err != nil {
		logger.Printf("Error getting encrypted password for user %s: %v", username, err)
		net.WriteSubsonicError(w, r, types.ErrorWrongCredentials, "Wrong username or password", "")
		return "", 0, false
	}

	if password != "" && validateWithPassword(username, password, encryptedPassword, w, r) {
		logger.Printf("User %s authenticated with plaintext password", username)
		return username, userId, true
	}

	if token != "" && salt != "" && validateWithTokenAndSalt(username, salt, token, encryptedPassword, w, r) {
		logger.Printf("User %s authenticated with salted password", username)
		return username, userId, true
	}

	net.WriteSubsonicError(w, r, types.ErrorWrongCredentials, "Wrong username or password", "")
	return "", 0, false
}

// validateWithApiKey checks if the provided API key is valid and returns the username and userId if successful.
func validateWithApiKey(ctx context.Context, apiKey string, w http.ResponseWriter, r *http.Request) (string, int64, bool) {
	user, err := database.ValidateApiKey(ctx, apiKey)
	if err != nil {
		logger.Printf("Error validating API key %s: %v", apiKey, err)
		net.WriteSubsonicError(w, r, types.ErrorInvalidApiKey, "Server Error", "")
		return "", 0, false
	}
	if user.Username == "" {
		logger.Printf("API key %s not found", apiKey)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "API Key not found", "")
		return "", 0, false
	}
	if user.AdminRole {
		logger.Printf("API key used for admin user %s", user.Username)
	} else {
		logger.Printf("API key used for user %s", user.Username)
	}
	return user.Username, user.Id, true
}

// validateWithPassword checks if the provided plaintext password matches the decrypted password from the database and returns true if valid.
func validateWithPassword(username, password, encryptedPassword string, w http.ResponseWriter, r *http.Request) bool {
	decryptedPassword, err := encryption.DecryptAES(encryptedPassword)
	if err != nil {
		logger.Printf("Error decrypting password for user %s: %v", username, err)
		net.WriteSubsonicError(w, r, types.ErrorWrongCredentials, "Wrong username or password", "")
		return false
	}
	return decryptedPassword == password
}

// validateWithTokenAndSalt checks if the provided token matches the computed token from the decrypted password and salt and returns true if valid.
func validateWithTokenAndSalt(username, salt, token, encryptedPassword string, w http.ResponseWriter, r *http.Request) bool {
	ok, err := validateToken(salt, token, encryptedPassword)
	if err != nil || !ok {
		logger.Printf("Error validating token for user %s: %v", username, err)
		net.WriteSubsonicError(w, r, types.ErrorWrongCredentials, "Wrong username or password", "")
		return false
	}
	return true
}

// AuthMiddleware is an HTTP middleware that authenticates requests using ValidateAuth.
// If authentication fails, it writes a 401 response and does not call the next handler.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userName, userId, ok := ValidateAuth(r, w)
		if !ok {
			// ValidateAuth already handled the error response.
			return
		}
		ctx := context.WithValue(r.Context(), "username", userName)
		ctx = context.WithValue(r.Context(), "userId", userId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminAuthMiddleware is an HTTP middleware that authenticates requests using ValidateAuth.
// If authentication fails, it writes a 401 response and does not call the next handler.
// It also ensures that the user is an admin by checking the "admin_role" field in the users table.
// If the user is not an admin, it writes a 403 Forbidden response.
func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userName, userId, ok := ValidateAuth(r, w)
		if !ok {
			// ValidateAuth already handled the error response.
			return
		}
		user, err := database.GetUserById(r.Context(), userId)
		if err != nil {
			logger.Printf("Error getting user by ID %d: %v", userId, err)
			net.WriteSubsonicError(w, r, types.ErrorGeneric, "Internal server error", "")
			return
		}
		if !user.AdminRole {
			logger.Printf("User %s (ID: %d) is not an admin", userName, userId)
			net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "User is not authorized for this operation", "")
			return
		}
		logger.Printf("Admin user %s (ID: %d) authenticated", userName, userId)
		ctx := context.WithValue(r.Context(), "username", userName)
		ctx = context.WithValue(r.Context(), "userId", userId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
