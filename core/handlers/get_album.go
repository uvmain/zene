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

func HandleGetAlbum(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	musicbrainzAlbumId := form["id"]

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

	if musicbrainzAlbumId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

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

	net.WriteSubsonicResponse(w, r, response, format)
}
