package handlers

import (
	"encoding/json"
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetShare(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	token := r.PathValue("share_token")
	ctx := r.Context()

	share, err := database.GetShareByToken(ctx, token)
	if err != nil {
		logger.Printf("Error getting share by token %s: %v", token, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to get share", "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(share); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
