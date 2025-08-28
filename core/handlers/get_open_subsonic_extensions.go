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

	response := subsonic.GetPopulatedSubsonicResponse(r.Context(), false)

	extension1 := types.OpenSubsonicExtensions{
		Name:     "formPost",
		Versions: []int{1},
	}
	extension2 := types.OpenSubsonicExtensions{
		Name:     "apiKeyAuthentication",
		Versions: []int{1},
	}
	extension3 := types.OpenSubsonicExtensions{
		Name:     "transcodeOffset",
		Versions: []int{1},
	}
	extension4 := types.OpenSubsonicExtensions{
		Name:     "songLyrics",
		Versions: []int{1},
	}

	response.SubsonicResponse.OpenSubsonicExtensions = []*types.OpenSubsonicExtensions{
		&extension1,
		&extension2,
		&extension3,
		&extension4,
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
