package handlers

import (
	"html"
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleCreateBookmark(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	id := form["id"]
	positionString := form["position"]
	comment := form["comment"]

	ctx := r.Context()

	if id == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "id parameter is required", "")
		return
	}

	if positionString == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "position parameter is required", "")
		return
	}

	var position int
	var err error
	if positionString != "" {
		position, err = strconv.Atoi(positionString)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "position parameter should be an integer", "")
			return
		}
	}

	comment = html.EscapeString(comment)

	validMusicbrainzId, validMetadata, err := database.IsValidMetadataId(ctx, id)
	if err != nil {
		logger.Printf("Error checking metadata ID: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to check metadata ID", "")
		return
	} else if !validMusicbrainzId || !validMetadata.MusicbrainzTrackId {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "id is not a valid musicbrainz track ID", "")
		return
	}

	logger.Printf("Creating bookmark for track ID %s at position %d with comment '%s'", id, position, comment)

	err = database.UpsertBookmark(ctx, id, position, comment)
	if err != nil {
		logger.Printf("Error inserting bookmark: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to add bookmark", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
