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

func HandleUpdateArtistArt(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	artistId := form["id"]
	artUrl := form["url"]

	if artistId == "" {
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
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "Admin role required to update artist art", "")
		return
	}

	valid, err := database.GetMediaCoverType(ctx, artistId)
	if err != nil {
		logger.Printf("Error getting media cover type: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error getting media cover type", "")
		return
	}
	if valid != "artist" {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Invalid artist id", "")
		return
	}

	filename := filepath.Join(config.ArtistArtFolder, artistId+".jpg")

	fileExists := io.FileExists(filename)
	if fileExists {
		err = io.DeleteFile(filename)
		if err != nil {
			logger.Printf("Error deleting existing artist art: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error deleting existing artist art", "")
			return
		}
	}

	if artUrl != "" {
		artImage, err := art.GetImageFromInternet(artUrl)
		if err != nil {
			logger.Printf("Error fetching artist art from URL: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error fetching artist art from URL", "")
			return
		}
		art.ResizeImageAndSaveAsJPG(artImage, filename, 512)
	} else {
		file, _, err := r.FormFile("file")
		if err != nil {
			logger.Printf("Error parsing file: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid file or url parameters", "")
			return
		}

		defer file.Close()

		if err := art.ResizeMultipartFileAndSaveAsJPG(file, filename, 512); err != nil {
			logger.Printf("Error uploading image: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Failed to upload new artist cover art", "")
			return
		}
	}

	err = database.UpsertArtistArtRow(ctx, artistId)
	if err != nil {
		logger.Printf("Error inserting artist art row: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error updating artist_art row", "")
		return
	}

	artist, err := database.SelectArtistByMusicBrainzArtistId(ctx, artistId)
	if err != nil {
		logger.Printf("Error getting artist: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Error getting artist", "")
		return
	}

	logger.Printf("Updated artist art for artist ID %s: %s", artistId, artist.Name)

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
