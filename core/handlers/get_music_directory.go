package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetMusicDirectory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()

	ifModifiedSinceHeader := r.Header.Get("If-Modified-Since")
	if ifModifiedSinceHeader != "" {
		latestScan, err := database.GetLatestCompletedScan(ctx)
		if err == nil {
			latestScanTime := logic.GetStringTimeFormatted(latestScan.CompletedDate)
			if net.IfModifiedResponse(w, r, latestScanTime) {
				return
			}
		}
	}

	musicbrainz_id := r.FormValue("id")

	if musicbrainz_id == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	isValidId, metadataResponse, err := database.IsValidMetadataId(ctx, musicbrainz_id)
	if err != nil || !isValidId {
		logger.Printf("Error getting metadata ID: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "provided ID is not a valid music directory ID", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	if metadataResponse.MusicbrainzAlbumId {
		directory, err := database.GetAlbumDirectory(ctx, musicbrainz_id)
		if err != nil {
			logger.Printf("Error getting album directory: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Album directory not found", "")
			return
		}
		response.SubsonicResponse.MusicDirectory = &directory
	} else if metadataResponse.MusicbrainzArtistId {
		directory, err := database.GetArtistDirectory(ctx, musicbrainz_id)
		if err != nil {
			logger.Printf("Error getting artist directory: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Artist directory not found", "")
			return
		}
		response.SubsonicResponse.MusicDirectory = &directory
	} else if metadataResponse.MusicbrainzTrackId {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Directory not found", "")
		return
	}

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
