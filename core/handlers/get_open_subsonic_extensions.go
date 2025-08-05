package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"zene/core/net"
	"zene/core/types"
)

func HandleOpenSubsonicExtensions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	response := types.SubsonicOpenSubsonicExtensionsResponse{}
	stdRes := types.GetPopulatedSubsonicResponse(false)

	response.SubsonicResponse.XMLName = stdRes.SubsonicResponse.XMLName
	response.SubsonicResponse.Xmlns = stdRes.SubsonicResponse.Xmlns
	response.SubsonicResponse.Status = stdRes.SubsonicResponse.Status
	response.SubsonicResponse.Version = stdRes.SubsonicResponse.Version
	response.SubsonicResponse.Type = stdRes.SubsonicResponse.Type
	response.SubsonicResponse.ServerVersion = stdRes.SubsonicResponse.ServerVersion
	response.SubsonicResponse.OpenSubsonic = stdRes.SubsonicResponse.OpenSubsonic

	extension1 := types.OpenSubsonicExtensions{
		Name:     "formPost",
		Versions: []int{1},
	}
	extension2 := types.OpenSubsonicExtensions{
		Name:     "apiKeyAuthentication",
		Versions: []int{1},
	}
	extension3 := types.OpenSubsonicExtensions{
		Name:     "transcodeOffset",
		Versions: []int{1},
	}
	extension4 := types.OpenSubsonicExtensions{
		Name:     "songLyrics",
		Versions: []int{1},
	}

	response.SubsonicResponse.OpenSubsonicExtensions = []*types.OpenSubsonicExtensions{
		&extension1,
		&extension2,
		&extension3,
		&extension4,
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
