package main

import (
	"log"
	"net/http"
	"zene/config"
	"zene/handlers"
	"zene/net"

	"github.com/rs/cors"
)

func StartServer() {
	router := http.NewServeMux()

	distDir := http.Dir("../dist")
	fileServer := http.FileServer(distDir)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := distDir.Open(r.URL.Path); err == nil {
			net.EnableCdnCaching(w)
			fileServer.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, "../dist/index.html")
	})

	router.HandleFunc("GET /api/files", handlers.HandleGetAllFiles)                                         //
	router.HandleFunc("GET /api/files/{fileId}", handlers.HandleGetFileById)                                //
	router.HandleFunc("GET /api/artists", handlers.HandleGetArtists)                                        // query params: search=searchTerm
	router.HandleFunc("GET /api/artists/{musicBrainzArtistId}", handlers.HandleGetArtist)                   //
	router.HandleFunc("GET /api/artists/{musicBrainzArtistId}/art", handlers.GetArtistArt)                  //
	router.HandleFunc("GET /api/albums", handlers.HandleGetAlbums)                                          // query params: recent=true, random=false, limit=10
	router.HandleFunc("GET /api/metadata", handlers.HandleGetMetadata)                                      // query params: recent=true, random=false, limit=10
	router.HandleFunc("GET /api/genres", handlers.HandleGetUniqueGenres)                                    // query params: search=searchTerm
	router.HandleFunc("GET /api/scan", handlers.HandlePostScan)                                             //
	router.HandleFunc("GET /api/search", handlers.HandleSearchMetadata)                                     // query params: search=searchTerm
	router.HandleFunc("GET /api/art/albums/{musicBrainzAlbumId}", handlers.GetAlbumArtByMusicBrainzAlbumId) //

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
