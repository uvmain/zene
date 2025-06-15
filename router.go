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
	"zene/core/handlers"
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
	router.Handle("GET /api/artists", auth.AdminAuthMiddleware(http.HandlerFunc(handlers.HandleGetArtists)))                              // returns []types.ArtistResponse; query params: search=searchTerm, recent=true, random=false, limit=10, offset=10
	router.Handle("GET /api/artists/{musicBrainzArtistId}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtist)))              // returns types.ArtistResponse
	router.Handle("GET /api/artists/{musicBrainzArtistId}/tracks", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtistTracks))) // returns []types.Metadata; query params: recent=true, random=false, limit=10, offset=10
	router.Handle("GET /api/artists/{musicBrainzArtistId}/art", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtistArt)))       // returns image/jpeg blob
	router.Handle("GET /api/artists/{musicBrainzArtistId}/albums", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtistAlbums))) // returns []types.AlbumsResponse
	router.Handle("GET /api/albums", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbums)))                                     // returns []types.AlbumsResponse; query params: recent=true, random=false, limit=10
	router.Handle("GET /api/albums/{musicBrainzAlbumId}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbum)))                 // returns types.AlbumsResponse
	router.Handle("GET /api/albums/{musicBrainzAlbumId}/art", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbumArt)))          // returns image/jpeg blob
	router.Handle("GET /api/albums/{musicBrainzAlbumId}/tracks", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbumTracks)))    // returns []types.Metadata
	router.Handle("GET /api/tracks", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetTracks)))                                     // returns []types.Metadata; query params: recent=true, random=false, limit=10
	router.Handle("GET /api/tracks/{musicBrainzTrackId}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetTrack)))                 // returns types.Metadata
	router.Handle("GET /api/tracks/{musicBrainzTrackId}/download", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDownloadTrack)))   // returns blob
	router.Handle("GET /api/tracks/{musicBrainzTrackId}/stream", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleStreamTrack)))       // returns blob range
	router.Handle("GET /api/genres", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetGenres)))                                     // query params: search=searchTerm
	router.Handle("GET /api/search", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSearchMetadata)))                                // query params: search=searchTerm
	router.Handle("GET /api/user", auth.AdminAuthMiddleware(http.HandlerFunc(handlers.HandleGetCurrentUser)))                             // return types.User - current user

	// admin routes
	router.Handle("POST /api/scan", auth.AdminAuthMiddleware(http.HandlerFunc(handlers.HandlePostScan)))                   // triggers a scan of the music library if one is not already running
	router.Handle("GET /api/users", auth.AdminAuthMiddleware(http.HandlerFunc(handlers.HandleGetAllUsers)))                // return []types.User - all users
	router.Handle("POST /api/users", auth.AdminAuthMiddleware(http.HandlerFunc(handlers.HandlePostNewUser)))               // return userId int64
	router.Handle("GET /api/users/{userId}", auth.AdminAuthMiddleware(http.HandlerFunc(handlers.HandleGetUserById)))       // return types.User - user by ID
	router.Handle("PATCH /api/users/{userId}", auth.AdminAuthMiddleware(http.HandlerFunc(handlers.HandlePatchUserById)))   // return userId int64
	router.Handle("DELETE /api/users/{userId}", auth.AdminAuthMiddleware(http.HandlerFunc(handlers.HandleDeleteUserById))) // return { Status: string }

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
