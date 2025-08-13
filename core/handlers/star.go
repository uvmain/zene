package handlers

import (
	"cmp"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleStar(w http.ResponseWriter, r *http.Request) {
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

	metadata_id := cmp.Or(r.FormValue("id"), r.FormValue("albumId"), r.FormValue("artistId"))
	if metadata_id == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "one of id, albumId, or artistId is required", "")
		return
	}

	err = database.UpsertUserStar(ctx, user.Id, metadata_id)
	if err != nil {
		logger.Printf("Error inserting user star for user %d: %v", user.Id, err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to add user star", "")
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
