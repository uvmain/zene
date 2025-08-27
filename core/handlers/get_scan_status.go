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

func HandleGetScanStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]

	scanStatus, err := database.GetLatestScan(r.Context())
	if err != nil {
		logger.Printf("Error getting scan status: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to get scan status", "")
		return
	}

	scanStatusResponse := types.ScanStatus{
		Scanning:      scanStatus.CompletedDate == "",
		Count:         scanStatus.Count,
		FolderCount:   scanStatus.FolderCount,
		StartedDate:   scanStatus.StartedDate,
		Type:          scanStatus.Type,
		CompletedDate: scanStatus.CompletedDate,
	}

	response := subsonic.GetPopulatedSubsonicResponse(r.Context(), false)

	response.SubsonicResponse.ScanStatus = &scanStatusResponse

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
