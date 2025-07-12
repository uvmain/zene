package handlers

import (
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func HandleGetPlaycounts(w http.ResponseWriter, r *http.Request) {
	userIdParam := r.URL.Query().Get("user_id")
	musicbrainzTrackId := r.URL.Query().Get("musicbrainz_track_id")
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		userId = 0
	}

	response := &PlaycountsResponse{}

	rows, err := database.GetPlayCounts(r.Context(), musicbrainzTrackId, userId)
	if err != nil {
		handleErrorResponse(w, response, "Failed to query database", err, http.StatusBadRequest)
		return
	}

	response.Playcounts = make([]*types.Playcount, len(rows))
	for i := range rows {
		response.Playcounts[i] = &rows[i]
	}

	handleSuccessResponse(w, response)
}

func HandleUpsertPlaycount(w http.ResponseWriter, r *http.Request) {
	userIdParam := r.FormValue("user_id")
	musicbrainzTrackId := r.FormValue("musicbrainz_track_id")

	response := &StandardResponse{}

	var userId int64
	var err error

	if userIdParam != "" {
		userId, err = strconv.ParseInt(userIdParam, 10, 64)
		if err != nil {
			handleErrorResponse(w, response, "Invalid user_id parameter", err, http.StatusBadRequest)
			return
		}
	} else {
		userId, err = logic.GetUserIdFromContext(r.Context())
		if err != nil {
			handleErrorResponse(w, response, "Failed to get user_id from context", err, http.StatusUnauthorized)
			return
		}
	}

	logger.Printf("upserting playcount for user_id: %v, musicbrainz_track_id: %s", userId, musicbrainzTrackId)

	err = database.UpsertPlayCount(r.Context(), userId, musicbrainzTrackId)
	if err != nil {
		handleErrorResponse(w, response, "Failed to upsert playcount", err, http.StatusInternalServerError)
		return
	}

	handleSuccessResponse(w, response)
}
