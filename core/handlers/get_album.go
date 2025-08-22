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

func HandleGetAlbum(w http.ResponseWriter, r *http.Request) {
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

	musicbrainzAlbumId := r.FormValue("id")

	if musicbrainzAlbumId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	album, err := database.GetAlbum(ctx, musicbrainzAlbumId)
	if err != nil {
		logger.Printf("Error getting album: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Album not found", "")
		return
	}
	response.SubsonicResponse.Album = &album

	songs, err := database.GetSongsForAlbum(ctx, musicbrainzAlbumId)
	if err != nil {
		logger.Printf("Error getting songs for album: %v", err)
	} else {
		response.SubsonicResponse.Album.Songs = songs
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
