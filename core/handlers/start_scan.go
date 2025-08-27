package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/scanner"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleStartScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]

	scanStatus, err := scanner.RunScan(r.Context())
	if err != nil {
		if scanStatus.Scanning == true {
			logger.Printf("Error starting scan: %v", scanStatus)
			net.WriteSubsonicError(w, r, types.ErrorGeneric, "A scan is already in progress. Please wait for it to complete before starting a new one.", "")
			return
		}
		logger.Printf("Error starting scan: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to start scan", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(r.Context(), false)

	response.SubsonicResponse.ScanStatus = &scanStatus

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
