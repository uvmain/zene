package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetInternetRadioStations(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]

	ctx := r.Context()

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	radioStations, err := database.GetInternetRadioStations(ctx)
	if err != nil {
		logger.Printf("Error fetching internet radio stations: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error fetching internet radio stations", "")
		return
	}

	response.SubsonicResponse.InternetRadioStations = &types.InternetRadioStations{
		InternetRadio: radioStations,
	}
	net.WriteSubsonicResponse(w, r, response, format)
}
