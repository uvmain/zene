package handlers

import (
	"net/http"
	"strconv"
	"zene/core/art"
	"zene/core/logger"
	"zene/core/net"
)

func HandleGetShareImg(w http.ResponseWriter, r *http.Request) {
	imageId := r.PathValue("imageId")

	sizeQueryParameter := r.FormValue("size")
	var sizeInt int
	var err error
	if sizeQueryParameter != "" {
		sizeInt, err = strconv.Atoi(sizeQueryParameter)
		if err != nil {
			logger.Printf("Error parsing size parameter in HandleGetShareImg: %v", err)
			http.Error(w, "Failed to parse size parameter", http.StatusBadRequest)
			return
		}
		logger.Printf("sizeInt is %d", sizeInt)
	}

	ctx := r.Context()

	imageBlob, lastModified, err := art.GetArtForAlbum(ctx, imageId, "xl")
	if err != nil {
		imageBlob, lastModified, err = art.GetArtForArtist(ctx, imageId)
		if err != nil {
			logger.Printf("Error getting image for %s: %v", imageId, err)
			http.Error(w, "Failed to get image", http.StatusInternalServerError)
			return
		}
	}

	if net.IfModifiedResponse(w, r, lastModified) {
		return
	}

	mimeType := http.DetectContentType(imageBlob)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(imageBlob)
}
