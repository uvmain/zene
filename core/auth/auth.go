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
func validateToken(username, salt string, token string, encryptedPassword string) (bool, error) {
	decryptedPassword, err := encryption.DecryptAES(encryptedPassword)
	if err != nil {
		logger.Printf("Error validating token for user %s: %v", username, err)
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
// - p: plaintext password, or hex encrypted password prefixed with "enc:"
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

	u := r.FormValue("u") // username
	p := r.FormValue("p") // plaintext password, or hex encrypted password prefixed with "enc:"
	t := r.FormValue("t") // token = sha256(password + salt)
	s := r.FormValue("s") // salt, used with token for authentication
	apiKey := r.FormValue("apiKey")
	v := r.FormValue("v") // protocol version, required for all requests
	c := r.FormValue("c") // client name, required for all requests

	if c == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Required parameter 'c' is missing", "")
		return "", 0, false
	}

	if v == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Required parameter 'v' is missing", "")
		return "", 0, false
	}

	if apiKey != "" && (u != "" || p != "" || t != "" || s != "") {
		net.WriteSubsonicError(w, r, types.ErrorTooManyAuthMechanisms, "Too many authentication mechanisms specified", "")
		return "", 0, false
	}

	if p != "" && (t != "" || s != "") {
		net.WriteSubsonicError(w, r, types.ErrorTooManyAuthMechanisms, "Too many authentication mechanisms specified", "")
		return "", 0, false
	}

	if apiKey != "" {
		return validateWithApiKey(ctx, apiKey, w, r)
	}

	if u == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Required parameter 'u' is missing", "")
		return "", 0, false
	}

	if p == "" && (t == "" || s == "") {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Either 'p' or both 't' and 's' parameters are required", "")
		return "", 0, false
	}

	encryptedPassword, userId, err := database.GetEncryptedPasswordFromDB(ctx, u)
	if err != nil {
		logger.Printf("Error getting encrypted password for user %s: %v", u, err)
		net.WriteSubsonicError(w, r, types.ErrorWrongCredentials, "Wrong username or password", "")
		return "", 0, false
	}

	if p != "" && validateWithPassword(u, p, encryptedPassword) {
		return u, userId, true
	}

	if t != "" && s != "" && validateWithTokenAndSalt(u, s, t, encryptedPassword) {
		return u, userId, true
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
	return user.Username, user.Id, true
}

// validateWithPassword checks if the provided password matches the decrypted password from the database and returns true if valid.
func validateWithPassword(username, password, encryptedPassword string) bool {
	decryptedPassword, err := encryption.DecryptAES(encryptedPassword)
	if err != nil {
		logger.Printf("Error decrypting password for user %s: %v", username, err)
		return false
	}
	// if password starts with "enc:" it is hex encrypted, we need to decrypt it first
	if len(password) > 4 && password[:4] == "enc:" {
		password, err = encryption.HexDecrypt(password[4:])
		if err != nil {
			logger.Printf("Error decrypting hex encoded password for user %s: %v", username, err)
			return false
		}
	}
	return decryptedPassword == password
}

// validateWithTokenAndSalt checks if the provided token matches the computed token from the decrypted password and salt and returns true if valid.
func validateWithTokenAndSalt(username, salt, token, encryptedPassword string) bool {
	ok, err := validateToken(username, salt, token, encryptedPassword)
	if err != nil || !ok {
		logger.Printf("Error validating token for user %s: %v", username, err)
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
