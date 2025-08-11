package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"zene/core/database"
	"zene/core/lyrics"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetLyricsBySongId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()

	response := types.SubsonicLyricsListResponseWrapper{}
	stdRes := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	response.SubsonicResponse.XMLName = stdRes.SubsonicResponse.XMLName
	response.SubsonicResponse.Xmlns = stdRes.SubsonicResponse.Xmlns
	response.SubsonicResponse.Status = stdRes.SubsonicResponse.Status
	response.SubsonicResponse.Version = stdRes.SubsonicResponse.Version
	response.SubsonicResponse.Type = stdRes.SubsonicResponse.Type
	response.SubsonicResponse.ServerVersion = stdRes.SubsonicResponse.ServerVersion
	response.SubsonicResponse.OpenSubsonic = stdRes.SubsonicResponse.OpenSubsonic

	musicBrainzTrackId := r.FormValue("id")
	if musicBrainzTrackId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	lyricsData, err := lyrics.GetLyricsForMusicBrainzTrackId(ctx, musicBrainzTrackId)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, fmt.Sprintf("Error fetching lyrics: %v", err), "")
		return
	}

	trackData, err := database.SelectTrack(ctx, musicBrainzTrackId)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, fmt.Sprintf("Error fetching track metadata: %v", err), "")
		return
	}

	lines := []types.StructuredLyricsLine{}
	if lyricsData.SyncedLyrics == "" {
		for _, line := range strings.Split(lyricsData.PlainLyrics, "\n") {
			lines = append(lines, types.StructuredLyricsLine{
				Value: line,
			})
		}
	} else {
		for _, line := range strings.Split(lyricsData.SyncedLyrics, "\n") {
			// each line is [MI:SS.ms] string, eg
			// [00:10.35] Thousands of hours, we wait to devour what's gone
			// split it into int (milliseconds) and text
			parts := strings.SplitN(line, "]", 2)
			if len(parts) == 2 {
				timePart := strings.TrimPrefix(parts[0], "[")
				textPart := strings.TrimSpace(parts[1])
				milliseconds, err := lyrics.ParseSyncLyricTimeToMilliseconds(timePart)
				if err == nil {
					lines = append(lines, types.StructuredLyricsLine{
						Start: milliseconds,
						Value: textPart,
					})
				}
			}
		}
	}

	response.SubsonicResponse.LyricsList = &types.SubsonicLyricsList{
		StructuredLyrics: []types.StructuredLyrics{
			{
				DisplayArtist: trackData.Artist,
				DisplayTitle:  trackData.Title,
				Lang:          "en",
				Offset:        0,
				Synced:        lyricsData.SyncedLyrics != "",
				Line:          lines,
			},
		},
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
