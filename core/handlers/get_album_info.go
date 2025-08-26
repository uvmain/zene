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

func HandleGetAlbumInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	var version int
	switch r.URL.Path {
	case "/rest/getAlbumInfo.view":
		version = 1
	case "/rest/getAlbumInfo2.view":
		version = 2
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

	valid, row, err := database.IsValidMetadataId(ctx, musicbrainzAlbumId)
	if err != nil || !valid || !row.MusicbrainzAlbumId {
		logger.Printf("album id is invalid: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "album id is invalid", "")
		return
	}

	shareUrl := logic.GetUnauthenticatedImageUrl(musicbrainzAlbumId)

	albumInfo := types.AlbumInfo{
		MusicBrainzId:  musicbrainzAlbumId,
		SmallImageUrl:  shareUrl + "?size=300",
		MediumImageUrl: shareUrl + "?size=600",
		LargeImageUrl:  shareUrl + "?size=1200",
	}

	switch version {
	case 1:
		response.SubsonicResponse.AlbumInfo = &albumInfo
	case 2:
		response.SubsonicResponse.AlbumInfo2 = &albumInfo
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
