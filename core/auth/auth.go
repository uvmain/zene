package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

var encryptionKey []byte

func getEncryptionKey() {
	key := os.Getenv("AUTH_ENCRYPTION_KEY")
	if key == "" || len(key) != 32 {
		logger.Println("*** AUTH_ENCRYPTION_KEY environment variable is not set or is not exactly 32 characters long")
		logger.Println("*** Using fallback key for development purposes only")
		key = "0123456789abcdef0123456789abcdef"

	}
	encryptionKey = []byte(key)
}

func encryptAES(plaintext string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptAES(cipherTextBase64 string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	return string(plaintext), err
}

// computeToken generates a SHA-256 token from a password and salt.
func computeToken(password, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return hex.EncodeToString(hash[:])
}

// validateToken checks if the provided token matches the computed token from the password and salt.
func validateToken(salt string, token string, encryptedPassword string) (bool, error) {
	decryptedPassword, err := decryptAES(encryptedPassword)
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
// - s: salt
// - t: token = sha256(password + salt)
// - apiKey: API key for authentication
//
// If apiKey is specified, then none of p, t, s, nor u can be specified.
// Else either p or both t and s must be specified.
func ValidateAuth(r *http.Request, w http.ResponseWriter) (string, int64, bool) {
	ctx := r.Context()

	apiKey := r.FormValue("apiKey")
	if apiKey != "" {
		user, err := database.ValidateApiKey(ctx, apiKey)
		if err != nil {
			logger.Printf("Error validating API key %s: %v", apiKey, err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return "", 0, false
		}
		if user.Username == "" {
			logger.Printf("API key %s not found", apiKey)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return "", 0, false
		}
		if user.IsDisabled {
			logger.Printf("API key %s belongs to a disabled user", apiKey)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return "", 0, false
		}
		if user.IsAdmin {
			logger.Printf("API key used for admin user %s", user.Username)
		} else {
			logger.Printf("API key used for user %s", user.Username)
		}
		return user.Username, user.Id, true
	}

	username := r.FormValue("u")
	if username == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Required parameter 'u' is missing")
		return "", 0, false
	}

	salt := r.FormValue("s")
	token := r.FormValue("t")
	password := r.FormValue("p")

	if password == "" && (token == "" || salt == "") {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Either 'p' or both 't' and 's' parameters are required")
		return "", 0, false
	}

	encryptedPassword, userId, err := database.GetEncryptedPasswordFromDB(ctx, username)
	if err != nil {
		logger.Printf("Error getting encrypted password for user %s: %v", username, err)
		net.WriteSubsonicError(w, r, types.ErrorWrongCredentials, "Wrong username or password")
		return "", 0, false
	}

	if password != "" && validateWithPassword(username, password, encryptedPassword, w, r) {
		// set context with username for further processing
		logger.Printf("User %s authenticated with plaintext password", username)
		return username, userId, true
	}

	if token != "" && salt != "" && validateWithTokenAndSalt(username, salt, token, encryptedPassword, w, r) {
		logger.Printf("User %s authenticated with salted password", username)
		return username, userId, true
	}

	net.WriteSubsonicError(w, r, types.ErrorWrongCredentials, "Wrong username or password")
	return "", 0, false
}

// validateWithPassword checks if the provided plaintext password matches the decrypted password from the database and returns true if valid.
func validateWithPassword(username string, password string, encryptedPassword string, w http.ResponseWriter, r *http.Request) bool {
	decryptedPassword, err := decryptAES(encryptedPassword)
	if err != nil {
		logger.Printf("Error decrypting password for user %s: %v", username, err)
		net.WriteSubsonicError(w, r, types.ErrorWrongCredentials, "Wrong username or password")
		return false
	}
	return decryptedPassword == password
}

// validateWithTokenAndSalt checks if the provided token matches the computed token from the decrypted password and salt and returns true if valid.
func validateWithTokenAndSalt(username string, salt string, token string, encryptedPassword string, w http.ResponseWriter, r *http.Request) bool {
	ok, err := validateToken(salt, token, encryptedPassword)
	if err != nil || !ok {
		logger.Printf("Error validating token for user %s: %v", username, err)
		net.WriteSubsonicError(w, r, types.ErrorWrongCredentials, "Wrong username or password")
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
// It also ensures that the user is an admin by checking the "isAdmin" field in the database.
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
			net.WriteSubsonicError(w, r, types.ErrorGeneric, "Internal server error")
			return
		}
		if !user.IsAdmin {
			logger.Printf("User %s (ID: %d) is not an admin", userName, userId)
			net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "User is not authorized for this operation")
			return
		}
		logger.Printf("Admin user %s (ID: %d) authenticated", userName, userId)
		ctx := context.WithValue(r.Context(), "username", userName)
		ctx = context.WithValue(r.Context(), "userId", userId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
