package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/deezer"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetArtistInfo(w http.ResponseWriter, r *http.Request) {
	var version int
	switch r.URL.Path {
	case "/rest/getArtistInfo.view":
		version = 1
	case "/rest/getArtistInfo2.view":
		version = 2
	}

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()

	musicBrainzArtistId := r.FormValue("id")
	if musicBrainzArtistId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	var countLimit = 20
	var err error
	count := r.FormValue("count")
	if count != "" {
		countLimit, err = strconv.Atoi(count)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "count parameter must be an integer", "")
			return
		}
	}

	includeNotPresent := r.FormValue("includeNotPresent")
	includeNotPresentBool := false
	if includeNotPresent != "" {
		includeNotPresentBool = net.ParseBooleanFromString(w, r, includeNotPresent)
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	valid, row, err := database.IsValidMetadataId(ctx, musicBrainzArtistId)
	if err != nil || !valid || !row.MusicbrainzArtistId {
		logger.Printf("artist id is invalid: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "artist id is invalid", "")
		return
	}

	shareUrl := logic.GetUnauthenticatedImageUrl(musicBrainzArtistId)

	artistInfo := types.ArtistInfo{
		MusicBrainzId:  musicBrainzArtistId,
		SmallImageUrl:  shareUrl + "?size=300",
		MediumImageUrl: shareUrl + "?size=600",
		LargeImageUrl:  shareUrl + "?size=1200",
	}

	switch version {
	case 1:
		response.SubsonicResponse.ArtistInfo = &artistInfo
	case 2:
		response.SubsonicResponse.ArtistInfo2 = &artistInfo
	}

	similarArtists := []types.Artist{}

	artistName := database.GetArtistNameByMusicBrainzArtistId(ctx, musicBrainzArtistId)

	if includeNotPresentBool {
		similarArtistNames, err := deezer.GetSimilarArtistNames(ctx, artistName)
		if err != nil {
			logger.Printf("failed to get similar artists: %v", err)
		} else {
			// if count is specified, limit the number of similar artists, otherwise default to a limit of 20
			for i, artistName := range similarArtistNames {
				if countLimit > 0 && i >= countLimit {
					break
				}
				artistId, err := database.GetArtistIdByName(ctx, artistName)
				if err == nil && artistId != "" {
					artist, err := database.SelectArtistByMusicBrainzArtistId(ctx, artistId)
					if err == nil {
						similarArtists = append(similarArtists, artist)
					} else {
						similarArtists = append(similarArtists, types.Artist{Name: artistName})
					}
				} else {
					similarArtists = append(similarArtists, types.Artist{Name: artistName})
				}
			}
		}
	} else {
		similarArtistsRows, err := database.SelectSimilarArtists(ctx, musicBrainzArtistId)
		if err != nil {
			logger.Printf("failed to get similar artists: %v", err)
		}
		similarArtists = append(similarArtists, similarArtistsRows...)
	}

	switch version {
	case 1:
		response.SubsonicResponse.ArtistInfo.SimilarArtists = similarArtists
	case 2:
		response.SubsonicResponse.ArtistInfo2.SimilarArtists = similarArtists
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
