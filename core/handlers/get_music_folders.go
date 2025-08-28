package handlers

import (
	"net/http"
	"zene/core/config"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetMusicFolders(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]

	response := subsonic.GetPopulatedSubsonicResponse(r.Context(), false)

	var musicFolders []types.MusicFolder

	for i := range config.MusicDirs {
		musicFolder := types.MusicFolder{
			Id:   i + 1,
			Name: config.MusicDirs[i],
		}
		musicFolders = append(musicFolders, musicFolder)
	}

	response.SubsonicResponse.MusicFolders = &types.MusicFolders{
		MusicFolder: musicFolders,
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
