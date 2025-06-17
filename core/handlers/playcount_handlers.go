package handlers

import (
	"net/http"
	"strconv"
	"zene/core/auth"
	"zene/core/database"
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

	rows, err := database.GetPlaycounts(r.Context(), musicbrainzTrackId, userId)
	if err != nil {
		handleErrorResponse(w, response, "Failed to query database", err, http.StatusBadRequest)
		return
	}

	response.Playcounts = make([]*types.Playcount, len(rows))
	for i := range rows {
		response.Playcounts[i] = &rows[i]
	}

	handleSuccessResponse(w, response)
	return
}

func HandleUpsertPlaycount(w http.ResponseWriter, r *http.Request) {
	userIdParam := r.URL.Query().Get("user_id")
	musicbrainzTrackId := r.URL.Query().Get("musicbrainz_track_id")

	response := &StandardResponse{}

	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		user, _, err := auth.GetUserFromRequest(r)
		if err != nil {
			handleErrorResponse(w, response, "Failed to get user_id", err, http.StatusInternalServerError)
		} else {
			userId = user.Id
		}
	}

	err = database.UpsertPlaycount(r.Context(), userId, musicbrainzTrackId)
	if err != nil {
		handleErrorResponse(w, response, "Failed to upsert playcount", err, http.StatusInternalServerError)
		return
	}

	handleSuccessResponse(w, response)
	return
}
