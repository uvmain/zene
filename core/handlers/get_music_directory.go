package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetMusicDirectory(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	musicbrainz_id := form["id"]

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

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

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

	net.WriteSubsonicResponse(w, r, response, format)
}
