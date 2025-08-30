package handlers

import (
	"cmp"
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleStar(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	metadataId := cmp.Or(form["id"], form["albumid"], form["artistid"])

	ctx := r.Context()

	user, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to add chat messages", "")
		return
	}

	if metadataId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "one of id, albumId, or artistId is required", "")
		return
	}

	err = database.UpsertUserStar(ctx, user.Id, metadataId)
	if err != nil {
		logger.Printf("Error inserting user star for user %d: %v", user.Id, err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to add user star", "")
		return
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
