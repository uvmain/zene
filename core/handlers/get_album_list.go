package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetAlbumList(w http.ResponseWriter, r *http.Request) {
	var version int
	switch strings.ToLower(r.URL.Path) {
	case "/rest/getalbumlist.view":
		version = 1
	case "/rest/getalbumlist2.view":
		version = 2
	}

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	typeParam := strings.ToLower(form["type"])
	sizeParam := form["size"]
	offsetParam := form["offset"]
	fromYearParam := form["fromyear"]
	toYearParam := form["toyear"]
	genreParam := form["genre"]
	musicFolderIdParam := form["musicfolderid"]

	ctx := r.Context()
	var err error

	if typeParam == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "type parameter is required", "")
		return
	} else if !(typeParam == "random" || typeParam == "newest" || typeParam == "highest" || typeParam == "frequent" || typeParam == "recent" ||
		typeParam == "alphabeticalbyname" || typeParam == "alphabeticalbyartist" || typeParam == "starred" || typeParam == "byyear" || typeParam == "bygenre") {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid type parameter", "")
		return
	}

	var sizeInt int
	if sizeParam != "" {
		sizeInt, err = strconv.Atoi(sizeParam)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "size parameter must be an integer", "")
			return
		}
	} else {
		sizeInt = 20
	}

	var offsetInt int
	if offsetParam != "" {
		offsetInt, err = strconv.Atoi(offsetParam)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "offset parameter must be an integer", "")
			return
		}
	} else {
		offsetInt = 0
	}

	if typeParam == "byyear" {
		var fromYearInt int
		var toYearInt int

		if fromYearParam != "" {
			fromYearInt, err = strconv.Atoi(fromYearParam)
			if err != nil {
				net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "fromYear parameter must be an integer", "")
				return
			}
		} else {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "fromYear parameter is required if type=byyear", "")
			return
		}
		if toYearParam != "" {
			toYearInt, err = strconv.Atoi(toYearParam)
			if err != nil {
				net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "toYear parameter must be an integer", "")
				return
			}
		} else {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "toYear parameter is required if type=byyear", "")
			return
		}
		if toYearInt < fromYearInt {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "toYear parameter must be greater than or equal to fromYear", "")
			return
		}
	}

	if typeParam == "bygenre" && genreParam == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "genre parameter is required if type=bygenre", "")
		return
	}

	var musicFolderIdInt int
	if musicFolderIdParam != "" {
		musicFolderIdInt, err = strconv.Atoi(musicFolderIdParam)
		if err != nil || musicFolderIdInt < 0 {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "musicFolderId parameter must be a positive integer", "")
			return
		}
	} else {
		musicFolderIdInt = 0
	}

	requestUser, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to get users", "")
		return
	}

	// if musicFolderId not in requestUser.MusicFolderIds
	if musicFolderIdInt != 0 && !slices.Contains(requestUser.Folders, musicFolderIdInt) {
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to access this music folder", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	albums, err := database.GetAlbumList(ctx, typeParam, sizeInt, offsetInt, fromYearParam, toYearParam, genreParam, musicFolderIdInt)
	if err != nil {
		logger.Printf("Error querying database in GetAlbumList: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to query database", "")
		return
	}
	if albums == nil {
		albums = []types.AlbumId3{}
	}

	switch version {
	case 1:
		response.SubsonicResponse.AlbumList = &types.AlbumList{}
		response.SubsonicResponse.AlbumList.Albums = albums
	case 2:
		response.SubsonicResponse.AlbumList2 = &types.AlbumList2{}
		response.SubsonicResponse.AlbumList2.Albums = albums
	}

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
