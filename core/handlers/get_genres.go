package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetGenres(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]

	ctx := r.Context()

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	genres, err := database.SelectDistinctGenres(r.Context())
	if err != nil {
		logger.Printf("Error querying database in SelectDistinctGenres: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to query database", "")
		return
	}

	response.SubsonicResponse.Genres = &types.Genres{
		Genre: genres,
	}

	if format == "json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))
		xml.NewEncoder(w).Encode(response.SubsonicResponse)
	}
}
