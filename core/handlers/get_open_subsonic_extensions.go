package handlers

import (
	"net/http"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleOpenSubsonicExtensions(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]

	response := subsonic.GetPopulatedSubsonicResponse(r.Context())

	extensions := []*types.OpenSubsonicExtensions{
		{
			Name:     "apiKeyAuthentication",
			Versions: []int{1},
		},
		{
			Name:     "getPodcastEpisode",
			Versions: []int{1},
		},
		{
			Name:     "formPost",
			Versions: []int{1},
		},
		{
			Name:     "indexBasedQueue",
			Versions: []int{1},
		},
		{
			Name:     "playbackReport",
			Versions: []int{1},
		},
		{
			Name:     "songLyrics",
			Versions: []int{1},
		},
		{
			Name:     "topSongsByArtistId",
			Versions: []int{1},
		},
		{
			Name:     "transcodeOffset",
			Versions: []int{1},
		},
		{
			Name:     "transcoding",
			Versions: []int{1},
		},
	}

	response.SubsonicResponse.OpenSubsonicExtensions = extensions

	net.WriteSubsonicResponse(w, r, response, format)
}
