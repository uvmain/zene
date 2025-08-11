package handlers

import (
	"fmt"
	"net/http"
	"zene/core/art"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetArtistArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()

	musicBrainzArtistId := r.FormValue("id")
	if musicBrainzArtistId == "" {
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

	imageBlob, lastModified, err := art.GetArtForArtist(ctx, musicBrainzArtistId)
	if err != nil {
		logger.Printf("Error getting cover art for artist %s: %v", musicBrainzArtistId, err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Artist cover art not found", "")
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
