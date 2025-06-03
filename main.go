package main

import (
	"zene/core/art"
	"zene/core/config"
	"zene/core/database"
	"zene/core/net"
	"zene/core/scanner"

	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	config.LoadConfig()

	database.Initialise()
	defer database.CloseDatabase()

	art.Initialise()

	go scanner.RunScan()

	StartServer()
}

//go:embed all:frontend/dist
var dist embed.FS

func StartServer() {
	router := http.NewServeMux()

	dist, err := fs.Sub(dist, "frontend/dist")
	if err != nil {
		log.Fatalf("sub error")
		return
	}
	router.Handle("GET /", http.FileServer(http.FS(dist)))

	router.HandleFunc("GET /api/files", net.HandleGetFiles)                                    // returns []types.FilesRow
	router.HandleFunc("GET /api/files/{fileId}", net.HandleGetFile)                            // returns types.FilesRow
	router.HandleFunc("GET /api/files/{fileId}/download", net.HandleDownloadFile)              // returns blob
	router.HandleFunc("GET /api/files/{fileId}/stream", net.HandleStreamFile)                  // returns blob
	router.HandleFunc("GET /api/artists", net.HandleGetArtists)                                // returns []types.ArtistResponse; query params: search=searchTerm, recent=true, random=false, limit=10
	router.HandleFunc("GET /api/artists/{musicBrainzArtistId}", net.HandleGetArtist)           // returns types.ArtistResponse
	router.HandleFunc("GET /api/artists/{musicBrainzArtistId}/art", net.GetArtistArt)          // returns image/jpeg blob
	router.HandleFunc("GET /api/albums", net.HandleGetAlbums)                                  // returns []types.AlbumsResponse; query params: recent=true, random=false, limit=10
	router.HandleFunc("GET /api/albums/{musicBrainzAlbumId}", net.HandleGetAlbum)              // returns types.AlbumsResponse
	router.HandleFunc("GET /api/albums/{musicBrainzAlbumId}/art", net.HandleGetAlbumArt)       // returns image/jpeg blob
	router.HandleFunc("GET /api/albums/{musicBrainzAlbumId}/tracks", net.HandleGetAlbumTracks) // returns []types.TrackMetadata
	router.HandleFunc("GET /api/tracks", net.HandleGetTracks)                                  // returns []types.TrackMetadata; query params: recent=true, random=false, limit=10
	router.HandleFunc("GET /api/tracks/{musicBrainzTrackId}", net.HandleGetTrack)              // returns types.TrackMetadata
	router.HandleFunc("GET /api/genres", net.HandleGetGenres)                                  // query params: search=searchTerm
	router.HandleFunc("POST /api/scan", net.HandlePostScan)                                    //
	router.HandleFunc("GET /api/search", net.HandleSearchMetadata)                             // query params: search=searchTerm

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
