package auth

import (
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
	"zene/core/logger"
)

var encryptionKey []byte

func init() {
	key := os.Getenv("AUTH_ENCRYPTION_KEY")
	if key == "" || len(key) != 32 {
		panic("AUTH_ENCRYPTION_KEY must be exactly 32 characters long")
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
// It returns the authenticated username and true if successful; otherwise returns false and writes a 401 response.
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
func ValidateAuth(r *http.Request, w http.ResponseWriter) (string, bool) {
	ctx := r.Context()

	apiKey := r.FormValue("apiKey")
	if apiKey != "" {
		user, err := database.ValidateApiKey(ctx, apiKey)
		if err != nil {
			logger.Printf("Error validating API key %s: %v", apiKey, err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return "", false
		}
		if user == nil {
			logger.Printf("API key %s not found", apiKey)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return "", false
		}
		if user.IsDisabled {
			logger.Printf("API key %s belongs to a disabled user", apiKey)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return "", false
		}
		if user.IsAdmin {
			logger.Printf("API key used for admin user %s", user.Username)
		} else {
			logger.Printf("API key used for user %s", user.Username)
		}
		return user.Username, true
	}

	username := r.FormValue("u")
	if username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return "", false
	}

	salt := r.FormValue("s")
	token := r.FormValue("t")
	password := r.FormValue("p")

	if password == "" && (token == "" || salt == "") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return "", false
	}

	encryptedPassword, err := database.GetEncryptedPasswordFromDB(ctx, username)
	if err != nil {
		logger.Printf("Error getting encrypted password for user %s: %v", username, err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return "", false
	}

	if password != "" && validateWithPassword(username, password, encryptedPassword, w) {
		return username, true
	}

	if token != "" && salt != "" && validateWithTokenAndSalt(username, salt, token, encryptedPassword, w) {
		return username, true
	}

	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	return "", false
}

// validateWithPassword checks if the provided plaintext password matches the decrypted password from the database and returns true if valid.
func validateWithPassword(username string, password string, encryptedPassword string, w http.ResponseWriter) bool {
	decryptedPassword, err := decryptAES(encryptedPassword)
	if err != nil {
		logger.Printf("Error decrypting password for user %s: %v", username, err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	return decryptedPassword == password
}

// validateWithTokenAndSalt checks if the provided token matches the computed token from the decrypted password and salt and returns true if valid.
func validateWithTokenAndSalt(username string, salt string, token string, encryptedPassword string, w http.ResponseWriter) bool {
	ok, err := validateToken(salt, token, encryptedPassword)
	if err != nil || !ok {
		logger.Printf("Error validating token for user %s: %v", username, err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}
