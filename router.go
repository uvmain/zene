package main

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
	"zene/core/auth"
	"zene/core/config"
	"zene/core/logic"
	"zene/core/net"

	"github.com/rs/cors"
)

//go:embed all:frontend/dist
var dist embed.FS
var distSubFS fs.FS
var err error

func StartServer() {
	router := http.NewServeMux()

	distSubFS, err = fs.Sub(dist, "frontend/dist")
	if err != nil {
		log.Fatal("Failed to create sub filesystem:", err)
	}

	router.HandleFunc("/", HandleFrontend)

	// auth
	router.HandleFunc("POST /api/login", auth.LoginHandler)
	router.HandleFunc("GET /api/logout", auth.LogoutHandler)
	router.HandleFunc("GET /api/check-session", auth.CheckSessionHandler)

	// authenticated routes
	router.Handle("GET /api/artists", auth.AuthMiddleware(http.HandlerFunc(net.HandleGetArtists)))                                   // returns []types.ArtistResponse; query params: search=searchTerm, recent=true, random=false, limit=10, offset=10
	router.Handle("GET /api/artists/{musicBrainzArtistId}", auth.AuthMiddleware(http.HandlerFunc(net.HandleGetArtist)))              // returns types.ArtistResponse
	router.Handle("GET /api/artists/{musicBrainzArtistId}/tracks", auth.AuthMiddleware(http.HandlerFunc(net.HandleGetArtistTracks))) // returns []types.TrackMetadata; query params: recent=true, random=false, limit=10, offset=10
	router.Handle("GET /api/artists/{musicBrainzArtistId}/art", auth.AuthMiddleware(http.HandlerFunc(net.HandleGetArtistArt)))       // returns image/jpeg blob
	router.Handle("GET /api/albums", auth.AuthMiddleware(http.HandlerFunc(net.HandleGetAlbums)))                                     // returns []types.AlbumsResponse; query params: recent=true, random=false, limit=10
	router.Handle("GET /api/albums/{musicBrainzAlbumId}", auth.AuthMiddleware(http.HandlerFunc(net.HandleGetAlbum)))                 // returns types.AlbumsResponse
	router.Handle("GET /api/albums/{musicBrainzAlbumId}/art", auth.AuthMiddleware(http.HandlerFunc(net.HandleGetAlbumArt)))          // returns image/jpeg blob
	router.Handle("GET /api/albums/{musicBrainzAlbumId}/tracks", auth.AuthMiddleware(http.HandlerFunc(net.HandleGetAlbumTracks)))    // returns []types.TrackMetadata
	router.Handle("GET /api/tracks", auth.AuthMiddleware(http.HandlerFunc(net.HandleGetTracks)))                                     // returns []types.TrackMetadata; query params: recent=true, random=false, limit=10
	router.Handle("GET /api/tracks/{musicBrainzTrackId}", auth.AuthMiddleware(http.HandlerFunc(net.HandleGetTrack)))                 // returns types.TrackMetadata
	router.Handle("GET /api/tracks/{musicBrainzTrackId}/download", auth.AuthMiddleware(http.HandlerFunc(net.HandleDownloadTrack)))   // returns blob
	router.Handle("GET /api/tracks/{musicBrainzTrackId}/stream", auth.AuthMiddleware(http.HandlerFunc(net.HandleStreamTrack)))       // returns blob
	router.Handle("GET /api/genres", auth.AuthMiddleware(http.HandlerFunc(net.HandleGetGenres)))                                     // query params: search=searchTerm
	// router.Handle("POST /api/scan", auth.AuthMiddleware(http.HandlerFunc(net.HandlePostScan)))                                       //
	router.Handle("GET /api/search", auth.AuthMiddleware(http.HandlerFunc(net.HandleSearchMetadata))) // query params: search=searchTerm

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

func HandleFrontend(w http.ResponseWriter, r *http.Request) {
	bootTime := logic.GetBootTime().Truncate(time.Second).UTC()

	cleanPath := path.Clean(r.URL.Path)
	if cleanPath == "/" {
		cleanPath = "/index.html"
	} else {
		cleanPath = strings.TrimPrefix(cleanPath, "/")
	}

	// static content
	file, err := distSubFS.Open(cleanPath)
	if err == nil {
		defer file.Close()

		if net.IfModifiedResponse(w, r, bootTime) {
			return
		}

		http.ServeContent(w, r, cleanPath, bootTime, file.(io.ReadSeeker))
		return
	}

	// serve index.html for vue-router content
	indexFile, err := distSubFS.Open("index.html")
	if err != nil {
		http.Error(w, "index.html not found", http.StatusNotFound)
		return
	}
	defer indexFile.Close()

	if net.IfModifiedResponse(w, r, bootTime) {
		return
	}

	http.ServeContent(w, r, "index.html", bootTime, indexFile.(io.ReadSeeker))
}
