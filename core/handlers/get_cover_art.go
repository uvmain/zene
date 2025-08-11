package handlers

import (
	"fmt"
	"net/http"
	"zene/core/art"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetCoverArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()

	musicBrainzAlbumId := r.FormValue("id")
	if musicBrainzAlbumId == "" {
		errorString := "invalid id parameter"
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, errorString, "")
		return
	}

	// sizeParam := r.FormValue("size")
	// var sizeInt int64
	// var err error
	// if sizeParam == "" {
	// 	sizeInt, err = strconv.ParseInt(sizeParam, 10, 64)
	// 	if err != nil {
	// 		errorString := "invalid size parameter"
	// 		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, errorString, "")
	// 		return
	// 	}
	// }

	imageBlob, lastModified, err := art.GetArtForAlbum(ctx, musicBrainzAlbumId, "xl")
	if err != nil {
		logger.Printf("Error getting cover art for album %s: %v", musicBrainzAlbumId, err)
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
