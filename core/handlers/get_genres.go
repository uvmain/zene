package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetGenres(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]

	ctx := r.Context()

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	genres, err := database.SelectDistinctGenres(r.Context())
	if err != nil {
		logger.Printf("Error querying database in SelectDistinctGenres: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to query database", "")
		return
	}

	response.SubsonicResponse.Genres = &types.Genres{
		Genre: genres,
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
