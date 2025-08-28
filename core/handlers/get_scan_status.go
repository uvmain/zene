package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetScanStatus(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
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

	net.WriteSubsonicResponse(w, r, response, format)
}
