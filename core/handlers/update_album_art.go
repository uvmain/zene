package handlers

import (
	"net/http"
	"path/filepath"
	"zene/core/art"
	"zene/core/config"
	"zene/core/database"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleUpdateAlbumArt(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	albumId := form["id"]

	if albumId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	ctx := r.Context()

	user, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error getting user", "")
		return
	}
	if !user.AdminRole {
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "Admin role required to update album art", "")
		return
	}

	valid, err := database.GetMediaCoverType(ctx, albumId)
	if err != nil {
		logger.Printf("Error getting media cover type: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error getting media cover type", "")
		return
	}
	if valid != "album" {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Invalid album id", "")
		return
	}

	filename := filepath.Join(config.AlbumArtFolder, albumId+".jpg")

	fileExists := io.FileExists(filename)
	if fileExists {
		err = io.DeleteFile(filename)
		if err != nil {
			logger.Printf("Error deleting existing album art: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error deleting existing album art", "")
			return
		}
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		logger.Printf("Error parsing file: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid file", "")
		return
	}

	defer file.Close()

	if err := art.ResizeMultipartFileAndSaveAsJPG(file, filename, 512); err != nil {
		logger.Printf("Error uploading image: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Failed to upload new album cover art", "")
		return
	}

	err = database.UpsertAlbumArtRow(ctx, albumId)
	if err != nil {
		logger.Printf("Error inserting album art row: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error updating album_art row", "")
		return
	}

	album, err := database.GetAlbum(ctx, albumId)
	if err != nil {
		logger.Printf("Error getting album: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error getting album", "")
		return
	}

	logger.Printf("Updated album art for album ID %s: %s", albumId, album.Name)

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
