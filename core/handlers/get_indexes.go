package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"slices"

	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetIndexes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()

	ifModifiedSince := r.FormValue("ifModifiedSince")
	musicFolderId := r.FormValue("musicFolderId")

	var ifModifiedSinceInt int64
	var musicFolderIdInt int
	var err error

	if ifModifiedSince != "" {
		ifModifiedSinceInt, err = strconv.ParseInt(ifModifiedSince, 10, 64)
		if err != nil || ifModifiedSinceInt < 0 {
			http.Error(w, "Invalid ifModifiedSince", http.StatusBadRequest)
			return
		}
	} else {
		ifModifiedSinceInt = 0
	}

	if musicFolderId != "" {
		musicFolderIdInt, err := strconv.Atoi(musicFolderId)
		if err != nil || musicFolderIdInt < 0 {
			http.Error(w, "Invalid musicFolderId", http.StatusBadRequest)
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

	queryMusicFolderInt64s := []int64{}
	// if the musicFolderId param is valid, use it - otherwise, use the user's folders
	if musicFolderIdInt != 0 {
		queryMusicFolderInt64s = append(queryMusicFolderInt64s, int64(musicFolderIdInt))
	} else {
		queryMusicFolderInt64s = append(queryMusicFolderInt64s, logic.IntSliceToInt64Slice(requestUser.Folders)...)
	}

	indexes, err := database.GetIndexes(ctx, queryMusicFolderInt64s, ifModifiedSinceInt)
	if err != nil {
		logger.Printf("Error querying database in GetIndexes: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to query database", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx, false)
	response.SubsonicResponse.Indexes = &indexes

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
