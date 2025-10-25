package handlers

import (
	"net/http"
	"strconv"
	"time"
	"zene/core/art"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetCoverArt(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	idParameter := form["id"]
	sizeParam := form["size"]

	ctx := r.Context()

	if idParameter == "" {
		errorString := "invalid id parameter"
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, errorString, "")
		return
	}

	var sizeInt = 400
	var err error
	if sizeParam != "" {
		sizeInt, err = strconv.Atoi(sizeParam)
		if err != nil {
			errorString := "invalid size parameter"
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, errorString, "")
			return
		}
	}

	mediaArtType, err := database.GetMediaCoverType(ctx, idParameter)
	if err != nil {
		errorString := "error getting media type from id parameter"
		logger.Printf("%s %s in HandleGetCoverArt: %v", errorString, idParameter, err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, errorString, "")
		return
	}

	var imageBlob []byte
	var lastModified time.Time

	switch mediaArtType {
	case "track":
		imageBlob, lastModified, err = art.GetArtForTrack(ctx, idParameter, sizeInt)
	case "album":
		imageBlob, lastModified, err = art.GetArtForAlbum(ctx, idParameter, sizeInt)
	case "artist":
		imageBlob, lastModified, err = art.GetArtForArtist(ctx, idParameter, sizeInt)
	case "podcast":
		imageBlob, lastModified, err = art.GetArtForPodcast(ctx, idParameter, sizeInt)
	default:
		logger.Printf("Error getting cover art for %s: %v", idParameter, err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Cover art not found", "")
		return
	}

	if err != nil {
		logger.Printf("Error getting cover art for %s: %v", idParameter, err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Cover art not found", "")
		return
	}

	if net.IfModifiedResponse(w, r, lastModified) {
		return
	}

	mimeType := http.DetectContentType(imageBlob)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(imageBlob)
}
