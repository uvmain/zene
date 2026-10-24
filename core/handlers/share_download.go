package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/types"
)

func HandleShareDownload(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	token := r.PathValue("share_token")
	ctx := r.Context()
	mediaId := form["mediaid"]

	if mediaId == "" {
		http.Error(w, "mediaId parameter is required", http.StatusBadRequest)
		return
	}

	_, tokenIsValid := database.ValidateTokenAndMediaId(ctx, token, mediaId)
	if !tokenIsValid {
		logger.Printf("Error validating share token %s", token)
		http.Error(w, "invalid token and media_id combo", http.StatusBadRequest)
		return
	}

	mediaFilepath, err := database.GetMediaFilePath(ctx, mediaId)

	if err != nil {
		logger.Printf("Error querying database for media filepath %s: %v", mediaId, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "File not found in database.", "")
		return
	}

	if mediaFilepath == "" {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "File not available to stream.", "")
		return
	}

	fileBlob, err := io.GetFileBlob(ctx, mediaFilepath)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "file not found", "")
		return
	}

	mimeType := http.DetectContentType(fileBlob)
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", "attachment; filename="+mediaFilepath)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(fileBlob)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "failed to write file", "")
		return
	}
}
