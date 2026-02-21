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

func HandleGetShareImg(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	imageId := r.PathValue("image_id")

	form := net.NormalisedForm(r, w)
	sizeQueryParameter := form["size"]

	var sizeInt = 400
	var err error
	if sizeQueryParameter != "" {
		sizeInt, err = strconv.Atoi(sizeQueryParameter)
		if err != nil {
			logger.Printf("Error parsing size parameter in HandleGetShareImg: %v", err)
			http.Error(w, "Failed to parse size parameter", http.StatusBadRequest)
			return
		}
	}

	ctx := r.Context()

	mediaArtType, err := database.GetMediaCoverType(ctx, imageId)
	if err != nil {
		errorString := "error getting media type from id parameter"
		logger.Printf("%s %s in HandleGetShareImg: %v", errorString, imageId, err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, errorString, "")
		return
	}

	var imageBlob []byte
	var lastModified time.Time

	switch mediaArtType {
	case "track":
		imageBlob, lastModified, err = art.GetArtForTrack(ctx, imageId, sizeInt)
	case "album":
		imageBlob, lastModified, err = art.GetArtForAlbum(ctx, imageId, sizeInt)
	case "artist":
		imageBlob, lastModified, err = art.GetArtForArtist(ctx, imageId, sizeInt)
	case "podcast":
		imageBlob, lastModified, err = art.GetArtForPodcast(ctx, imageId, sizeInt)
	default:
		logger.Printf("Error getting cover art for %s: %v", imageId, err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Cover art not found", "")
		return
	}

	if err != nil {
		logger.Printf("Error getting cover art for %s: %v", imageId, err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Cover art not found", "")
		return
	}

	if net.IfModifiedResponse(w, r, lastModified) {
		return
	}

	mimeType := http.DetectContentType(imageBlob)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(imageBlob)
}
