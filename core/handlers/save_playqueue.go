package handlers

import (
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleSaveOrClearPlayqueue(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	positionString := form["position"]
	current := form["current"]
	currentIndexString := form["currentindex"]
	client := form["c"]

	ctx := r.Context()

	_, idArray, err := net.ParseDuplicateFormKeys(r, "id", false)
	if err != nil {
		logger.Printf("Error parsing id parameters: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "error parsing id parameters", "")
		return
	}

	if len(idArray) == 0 && (current != "" || currentIndexString != "") {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "at least one id parameter is required if current or currentIndex is set", "")
		return
	}

	var position int
	if positionString != "" {
		position, err = strconv.Atoi(positionString)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "position parameter should be an integer", "")
			return
		}
	}

	var indexInt int
	if current != "" {
		found := false
		for index, id := range idArray {
			if id == current {
				found = true
				indexInt = index
				break
			}
		}
		if !found {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "current parameter must be one of the id parameters", "")
			return
		}
	} else if currentIndexString != "" {
		indexInt, err = strconv.Atoi(currentIndexString)
		if err != nil {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "currentIndex parameter should be an integer", "")
			return
		}
		if indexInt < 0 || indexInt >= len(idArray) {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "currentIndex parameter is out of range", "")
			return
		}
	}

	if len(idArray) == 0 {
		err = database.ClearPlayqueue(ctx)
		if err != nil {
			logger.Printf("Error clearing playqueue: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to clear playqueue", "")
			return
		}
	} else if len(idArray) > 0 {
		err = database.UpsertPlayqueue(ctx, idArray, indexInt, position, client)
		if err != nil {
			logger.Printf("Error inserting playqueue: %v", err)
			net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Failed to add playqueue", "")
			return
		}
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
