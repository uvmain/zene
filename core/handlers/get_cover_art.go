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

	valid, metadataStruct, err := database.IsValidMetadataId(ctx, idParameter)
	if err != nil || !valid {
		errorString := "invalid id parameter"
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, errorString, "")
		return
	}

	var imageBlob []byte
	var lastModified time.Time
	if metadataStruct.MusicbrainzAlbumId {
		imageBlob, lastModified, err = art.GetArtForAlbum(ctx, idParameter, sizeInt)
		if err != nil {
			logger.Printf("Error getting album cover art for %s: %v", idParameter, err)
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Cover art not found", "")
			return
		}
	} else if metadataStruct.MusicbrainzArtistId {
		imageBlob, lastModified, err = art.GetArtForArtist(ctx, idParameter, sizeInt)
		if err != nil {
			logger.Printf("Error getting artist cover art for %s: %v", idParameter, err)
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Cover art not found", "")
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
