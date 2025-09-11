package handlers

import (
	"net/http"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/scanner"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleStartScan(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]

	scanStatus, err := scanner.RunScan(r.Context())
	if err != nil {
		if scanStatus.Scanning {
			logger.Printf("Error starting scan: %v", scanStatus)
			net.WriteSubsonicError(w, r, types.ErrorGeneric, "A scan is already in progress. Please wait for it to complete before starting a new one.", "")
			return
		}
		logger.Printf("Error starting scan: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to start scan", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(r.Context())

	response.SubsonicResponse.ScanStatus = &scanStatus

	net.WriteSubsonicResponse(w, r, response, format)
}
