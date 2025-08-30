package handlers

import (
	"net/http"
	"strconv"
	"strings"
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
	switch strings.ToLower(r.URL.Path) {
	case "/rest/getartistinfo":
		version = 1
	case "/rest/getartistinfo.view":
		version = 1
	case "/rest/getartistinfo2":
		version = 2
	case "/rest/getartistinfo2.view":
		version = 2
	}

	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	musicBrainzArtistId := form["id"]
	count := form["count"]
	includeNotPresent := form["includenotpresent"]

	ctx := r.Context()

	if musicBrainzArtistId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	var countLimit = 20
	var err error
	if count != "" {
		countLimit, err = strconv.Atoi(count)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "count parameter must be an integer", "")
			return
		}
	}

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

	shareUrl := logic.GetUnauthenticatedImageUrl(musicBrainzArtistId, 600)

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
						similarArtists = append(similarArtists, types.Artist{Name: artistName, Id: "-1"})
					}
				} else {
					similarArtists = append(similarArtists, types.Artist{Name: artistName, Id: "-1"})
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

	net.WriteSubsonicResponse(w, r, response, format)
}
