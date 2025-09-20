package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/io"
	"zene/core/net"
	"zene/core/types"
)

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	musicBrainzTrackId := form["id"]

	ctx := r.Context()

	if musicBrainzTrackId == "" {
		errorString := "invalid id parameter"
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, errorString, "")
		return
	}

	track, err := database.SelectTrack(ctx, musicBrainzTrackId)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "track ID not found", "")
		return
	}

	fileBlob, err := io.GetFileBlob(ctx, track.FilePath)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "file not found", "")
		return
	}

	mimeType := http.DetectContentType(fileBlob)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(fileBlob)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "failed to write file", "")
		return
	}
}
