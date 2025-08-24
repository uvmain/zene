package handlers

import (
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

func HandleGetArtist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()

	musicBrainzArtistId := r.FormValue("id")
	if musicBrainzArtistId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	row, err := database.SelectArtistByMusicBrainzArtistId(ctx, musicBrainzArtistId)
	if err != nil {
		logger.Printf("Error querying database in SelectArtistByMusicBrainzArtistId: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to query database", "")
		return
	}

	response.SubsonicResponse.Artist = &types.SubsonicArtistWrapper{}
	response.SubsonicResponse.Artist.Artist = row

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
