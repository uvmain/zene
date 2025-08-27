package auth

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"zene/core/database"
	"zene/core/encryption"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func generateExpectedToken(password, salt string) string {
	sum := md5.Sum([]byte(password + salt))
	return hex.EncodeToString(sum[:])
}

func validateToken(salt string, token string, decryptedPassword string) bool {
	expected := generateExpectedToken(decryptedPassword, salt)
	return token == expected
}

// ValidateAuth authenticates a user by either plaintext password or token and salt.
// It returns:
// string: the authenticated username
// int: the authenticated userId
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
func ValidateAuth(r *http.Request, w http.ResponseWriter) (string, int, bool) {
	ctx := r.Context()
	form := net.NormalisedForm(r, w)

	u := form["u"] // username
	p := form["p"] // plaintext password, or hex encrypted password prefixed with "enc:"
	t := form["t"] // token = sha256(password + salt)
	s := form["s"] // salt, used with token for authentication
	apiKey := form["apikey"]
	v := form["v"] // protocol version, required for all requests
	c := form["c"] // client name, required for all requests

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

	if t != "" && s != "" && validateWithTokenAndSalt(s, t, encryptedPassword) {
		return u, userId, true
	}

	net.WriteSubsonicError(w, r, types.ErrorWrongCredentials, "Wrong username or password", "")
	return "", 0, false
}

// validateWithApiKey checks if the provided API key is valid and returns the username and userId if successful.
func validateWithApiKey(ctx context.Context, apiKey string, w http.ResponseWriter, r *http.Request) (string, int, bool) {
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
func validateWithTokenAndSalt(salt, token, encryptedPassword string) bool {
	decryptedPassword, err := encryption.DecryptAES(encryptedPassword)
	if err != nil {
		logger.Printf("Error decrypting password: %v", err)
		return false
	}
	ok := validateToken(salt, token, decryptedPassword)
	return ok
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
		ctx = context.WithValue(ctx, "userId", userId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
