package handlers

import (
	"encoding/json"
	"net/http"
	"zene/core/database"
	"zene/core/ffmpeg"
	"zene/core/ffprobe"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleDownloadFfBinaries(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	ctx := r.Context()

	var err error

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to get users", "")
		return
	}

	if !requestUser.AdminRole {
		logger.Printf("User %s attempted to update ffmpeg binaries without admin role", requestUser.Username)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to update ffmpeg binaries", "")
		return
	}

	err = ffmpeg.DownloadFfmpegBinary()
	if err != nil {
		logger.Printf("Error downloading ffmpeg binary: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error downloading ffmpeg binary", "")
		return
	}

	err = ffprobe.DownloadFfprobeBinary()
	if err != nil {
		logger.Printf("Error downloading ffprobe binary: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error downloading ffprobe binary", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}

func HandleGetFfVersions(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	ctx := r.Context()

	ffmpegVersion, err := ffmpeg.GetFfmpegVersion(ctx)
	if err != nil {
		logger.Printf("Error getting ffmpeg version: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error getting ffmpeg version", "")
		return
	}
	ffprobeVersion, err := ffprobe.GetFfprobeVersion(ctx)
	if err != nil {
		logger.Printf("Error getting ffprobe version: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error getting ffprobe version", "")
		return
	}

	versions := types.FfVersions{
		FfmpegVersion:  ffmpegVersion,
		FfprobeVersion: ffprobeVersion,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(versions); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
