package main

import (
	"log"
	"net/http"
	"time"
	"zene/config"
	"zene/handlers"

	"github.com/rs/cors"
)

func enableCdnCaching(w http.ResponseWriter) {
	expiryDate := time.Now().AddDate(1, 0, 0)
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("Expires", expiryDate.String())
}

func StartServer() {
	router := http.NewServeMux()

	distDir := http.Dir("../dist")
	fileServer := http.FileServer(distDir)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := distDir.Open(r.URL.Path); err == nil {
			enableCdnCaching(w)
			fileServer.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, "../dist/index.html")
	})

	router.HandleFunc("/api/files", handlers.HandleGetAllFiles)
	router.HandleFunc("/api/file", handlers.HandleGetFileByName)
	router.HandleFunc("/api/artists", handlers.HandleGetArtists)

	handler := cors.AllowAll().Handler(router)

	var serverAddress string
	if config.IsLocalDevEnv() {
		serverAddress = "localhost:8080"
		log.Println("Application running at http://localhost:8080")
	} else {
		serverAddress = ":8080"
		log.Println("Application running at :8080")
	}

	http.ListenAndServe(serverAddress, handler)
}
