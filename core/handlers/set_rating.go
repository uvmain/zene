package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleSetRating(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()

	user, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to add chat messages", "")
		return
	}

	metadataId := r.FormValue("id")
	if metadataId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id is required", "")
		return
	}

	rating := r.FormValue("rating")
	if rating == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "rating is required", "")
		return
	}

	ratingInt, err := strconv.Atoi(rating)
	if err != nil || ratingInt < 0 || ratingInt > 5 {
		logger.Printf("Error parsing rating for user %d: %v", user.Id, err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid rating, must be between 0 and 5", "")
		return
	}

	err = database.UpsertUserRating(ctx, user.Id, metadataId, ratingInt)
	if err != nil {
		logger.Printf("Error upserting rating for user %d: %v", user.Id, err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to upsert user rating", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

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
