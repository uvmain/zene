package handlers

import (
	"encoding/json"
	"strconv"

	"net/http"
	"zene/core/butterchurn"
	"zene/core/net"
	"zene/core/types"
)

func HandleGetButterchurnPresets(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	random := form["random"]
	count := form["count"]

	if random == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "random parameter is required", "")
		return
	}

	var countInt int
	if count == "" {
		countInt = 0
	} else {
		countInt, _ = strconv.Atoi(count)
	}

	presets, err := butterchurn.GetPresets(countInt, random == "true")
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Error loading presets", "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(presets); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
