package handlers

import (
	"encoding/json"
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
)

func HandleShareIncrementVisitCount(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	token := r.PathValue("share_token")
	ctx := r.Context()

	shareId, tokenIsValid := database.ValidateShareToken(ctx, token)
	if !tokenIsValid {
		logger.Printf("Error validating share token %s", token)
		http.Error(w, "invalid token and media_id combo", http.StatusBadRequest)
		return
	}

	err := database.IncrementShareVisitCount(ctx, shareId)
	if err != nil {
		logger.Printf("Error incrementing share visit count for share %s: %v", token, err)
		http.Error(w, "Failed to increment visit count.", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
