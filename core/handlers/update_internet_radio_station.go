package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleUpdateInternetRadioStation(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	radioStationId := form["id"]
	streamUrl := form["streamurl"]
	stationName := form["name"]
	homepageUrl := form["homepageurl"]

	if radioStationId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is mandatory", "")
		return
	}

	if streamUrl == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "streamUrl parameter is mandatory", "")
		return
	}

	if stationName == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "name parameter is mandatory", "")
		return
	}

	ctx := r.Context()

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	err := database.UpdateInternetRadioStation(ctx, radioStationId, stationName, streamUrl, homepageUrl)
	if err != nil {
		logger.Printf("Error updating internet radio station: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error updating internet radio station", "")
		return
	}

	net.WriteSubsonicResponse(w, r, response, format)
}
