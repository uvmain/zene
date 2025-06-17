package handlers

import (
	"net/http"
	"strconv"
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
