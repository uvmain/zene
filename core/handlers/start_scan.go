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

	scanResult := scanner.RunScan(r.Context())
	if scanResult.Success != true {
		logger.Printf("Error starting scan: %v", scanResult.Status)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to start scan", "")
		return
	}

	response := types.SubsonicScanStatusResponse{}
	stdRes := subsonic.GetPopulatedSubsonicResponse(r.Context(), false)

	response.SubsonicResponse.XMLName = stdRes.SubsonicResponse.XMLName
	response.SubsonicResponse.Xmlns = stdRes.SubsonicResponse.Xmlns
	response.SubsonicResponse.Status = stdRes.SubsonicResponse.Status
	response.SubsonicResponse.Version = stdRes.SubsonicResponse.Version
	response.SubsonicResponse.Type = stdRes.SubsonicResponse.Type
	response.SubsonicResponse.ServerVersion = stdRes.SubsonicResponse.ServerVersion
	response.SubsonicResponse.OpenSubsonic = stdRes.SubsonicResponse.OpenSubsonic

	response.SubsonicResponse.ScanStatus = &types.ScanStatus{
		Scanning:    scanResult.Success,
		Count:       0,
		FolderCount: 0,
		LastScan:    "",
		ScanType:    "",
		ElapsedTime: 0,
	}

	format := r.FormValue("f")
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
