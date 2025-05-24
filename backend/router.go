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

	router.HandleFunc("GET /api/files", handlers.HandleGetFiles)                              // returns []types.FilesRow
	router.HandleFunc("GET /api/files/{fileId}", handlers.HandleGetFile)                      // returns types.FilesRow
	router.HandleFunc("GET /api/files/{fileId}/download", handlers.HandleDownloadFile)        // returns blob
	router.HandleFunc("GET /api/artists", handlers.HandleGetArtists)                          // returns []types.ArtistResponse; query params: search=searchTerm
	router.HandleFunc("GET /api/artists/{musicBrainzArtistId}", handlers.HandleGetArtist)     // returns types.ArtistResponse
	router.HandleFunc("GET /api/artists/{musicBrainzArtistId}/art", handlers.GetArtistArt)    // returns image/jpeg blob
	router.HandleFunc("GET /api/albums", handlers.HandleGetAlbums)                            // returns []types.AlbumsResponse; query params: recent=true, random=false, limit=10
	router.HandleFunc("GET /api/albums/{musicBrainzAlbumId}", handlers.HandleGetAlbum)        // returns types.AlbumsResponse
	router.HandleFunc("GET /api/albums/{musicBrainzAlbumId}/art", handlers.HandleGetAlbumArt) // returns image/jpeg blob
	router.HandleFunc("GET /api/tracks", handlers.HandleGetTracks)                            // returns []types.TrackMetadata; query params: recent=true, random=false, limit=10
	router.HandleFunc("GET /api/tracks/{musicBrainzTrackId}", handlers.HandleGetTrack)        // returns types.TrackMetadata
	router.HandleFunc("GET /api/genres", handlers.HandleGetGenres)                            // query params: search=searchTerm
	router.HandleFunc("POST /api/scan", handlers.HandlePostScan)                              //
	router.HandleFunc("GET /api/search", handlers.HandleSearchMetadata)                       // query params: search=searchTerm

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
