package handlers

import (
	"net/http"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetBookmarks(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]

	ctx := r.Context()

	bookmarks, err := database.GetBookmarks(ctx)
	if err != nil {
		logger.Printf("Error getting bookmarks: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to get bookmarks", "")
		return
	}
	if bookmarks == nil {
		bookmarks = []types.Bookmark{}
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	response.SubsonicResponse.Bookmarks = &types.Bookmarks{}
	response.SubsonicResponse.Bookmarks.Bookmarks = bookmarks

	net.WriteSubsonicResponse(w, r, response, format)
}
