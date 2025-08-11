package handlers

import (
	"net/http"
)

func HandleUpdateAvatar(w http.ResponseWriter, r *http.Request) {
	HandleCreateAvatar(w, r)
}
