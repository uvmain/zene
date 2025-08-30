package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleDeleteInternetRadioStation(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	radioStationId := form["id"]
	if radioStationId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is mandatory", "")
		return
	}

	ctx := r.Context()

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	err := database.DeleteInternetRadioStation(ctx, radioStationId)
	if err != nil {
		logger.Printf("Error deleting internet radio station: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error deleting internet radio station", "")
		return
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
