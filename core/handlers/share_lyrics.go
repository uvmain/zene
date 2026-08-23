package handlers

import (
	"encoding/json"
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/lyrics"
	"zene/core/net"
)

func HandleShareLyrics(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	token := r.PathValue("share_token")
	form := net.NormalisedForm(r, w)
	mediaId := form["mediaid"]

	ctx := r.Context()

	_, tokenIsValid := database.ValidateTokenAndMediaId(ctx, token, mediaId)
	if !tokenIsValid {
		logger.Printf("Error validating share token %s", token)
		http.Error(w, "invalid token and media_id combo", http.StatusBadRequest)
		return
	}

	if mediaId == "" {
		http.Error(w, "mediaid parameter is required", http.StatusBadRequest)
		return
	}

	lyricsData, err := lyrics.GetLyricsForMusicBrainzTrackId(ctx, mediaId)
	if err != nil {
		http.Error(w, "Error fetching lyrics", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(lyricsData.SyncedLyrics); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
