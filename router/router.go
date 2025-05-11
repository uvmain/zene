package router

import (
	"net/http"
	"zene/handlers"
)

func CreateRoutes() {
	http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGetAllFiles(w)
	})
	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGetFileByName(w, r)
	})
}
