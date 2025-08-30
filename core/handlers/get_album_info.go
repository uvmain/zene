package handlers

import (
	"strings"

	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetAlbumInfo(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	var version int
	switch strings.ToLower(r.URL.Path) {
	case "/rest/getalbuminfo":
		version = 1
	case "/rest/getalbuminfo.view":
		version = 1
	case "/rest/getalbuminfo2":
		version = 2
	case "/rest/getalbuminfo2.view":
		version = 2
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

	valid, row, err := database.IsValidMetadataId(ctx, musicbrainzAlbumId)
	if err != nil || !valid || !row.MusicbrainzAlbumId {
		logger.Printf("album id is invalid: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "album id is invalid", "")
		return
	}

	shareUrl := logic.GetUnauthenticatedImageUrl(musicbrainzAlbumId, 600)

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

	net.WriteSubsonicResponse(w, r, response, format)
}
