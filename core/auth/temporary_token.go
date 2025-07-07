package auth

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"
)

type TemporaryToken struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
}

func GetTemporaryTokenHandler(w http.ResponseWriter, r *http.Request) {
	duration := r.URL.Query().Get("duration")

	durationInt := 5 // default to 5 minutes
	if duration != "" {
		var err error
		durationInt, err = strconv.Atoi(duration)
		if err != nil || durationInt < 0 {
			http.Error(w, "Invalid duration", http.StatusBadRequest)
			return
		}
	}

	ctx := r.Context()

	userId, err := logic.GetUserIdFromContext(ctx)
	if err != nil {
		logger.Printf("Unauthorized access attempt, invalid session: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := generateToken()
	if err != nil {
		logger.Println("Error generating token:", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	expiresAt, err := database.SaveTemporaryToken(ctx, userId, token, time.Minute*time.Duration(durationInt))
	if err != nil {
		logger.Println("Error saving token:", err)
		http.Error(w, "Error saving token", http.StatusInternalServerError)
		return
	}

	response := TemporaryToken{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	logger.Printf("Token provisioned for user ID %d, duration %d minutes, expires at %s", userId, durationInt, expiresAt)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Println("Error encoding GetTemporaryTokenHandler response:", err)
		http.Error(w, "Error encoding GetTemporaryTokenHandler response", http.StatusInternalServerError)
		return
	}
}

func ValidateToken(tokenStr string) (bool, error) {
	return true, nil
}
