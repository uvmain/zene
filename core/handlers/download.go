package handlers

import (
	"fmt"
	"net/http"
	"zene/core/database"
	"zene/core/io"
	"zene/core/net"
	"zene/core/types"
)

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
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
	w.Write(fileBlob)
}
