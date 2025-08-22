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
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"

	"github.com/NYTimes/gziphandler"
	"github.com/rs/cors"
)

//go:embed all:frontend/dist
var dist embed.FS
var distSubFS fs.FS
var err error

func StartServer() *http.Server {
	router := http.NewServeMux()

	distSubFS, err = fs.Sub(dist, "frontend/dist")
	if err != nil {
		log.Fatal("Failed to create sub filesystem:", err)
	}

	// authenticated routes
	router.Handle("/api/albums", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbums)))
	router.Handle("/api/albums/{musicBrainzAlbumId}/tracks", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbumTracks)))
	router.Handle("/api/tracks", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetTracks)))
	router.Handle("/api/tracks/{musicBrainzTrackId}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetTrack)))
	router.Handle("/api/genres/tracks", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetTracksByGenre)))
	router.Handle("/api/search", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSearchMetadata)))
	router.Handle("/api/playcounts", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetPlaycounts)))

	router.Handle("GET /share/img/{imageId}", http.HandlerFunc(handlers.HandleGetShareImg))

	// OpenSubsonic routes
	/// System
	router.Handle("/rest/ping.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandlePing)))
	router.Handle("/rest/getLicense.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleLicense)))
	router.Handle("/rest/getOpenSubsonicExtensions.view", http.HandlerFunc(handlers.HandleOpenSubsonicExtensions))
	router.Handle("/rest/tokenInfo.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleTokenInfo)))
	/// Browsing
	router.Handle("/rest/getMusicFolders.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetMusicFolders)))
	router.Handle("/rest/getIndexes.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetIndexes)))
	router.Handle("/rest/getMusicDirectory.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetMusicDirectory)))
	router.Handle("/rest/getArtists.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtists)))
	router.Handle("/rest/getArtist.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtist)))
	router.Handle("/rest/getAlbum.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbum)))
	router.Handle("/rest/getSong.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetSong)))
	router.Handle("/rest/getVideos.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetVideos)))
	router.Handle("/rest/getVideoInfo.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetVideoInfo)))
	router.Handle("/rest/getGenres.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetGenres)))
	// Album/song lists
	router.Handle("/rest/getSongsByGenre.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetSongsByGenre)))
	router.Handle("/rest/getNowPlaying.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetNowPlaying)))
	// Media retrieval
	router.Handle("/rest/stream.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleStream)))
	router.Handle("/rest/download.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDownload)))
	router.Handle("/rest/getCaptions.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetCaptions)))
	router.Handle("/rest/getCoverArt.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetCoverArt)))
	router.Handle("/rest/getArtistArt.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtistArt)))
	router.Handle("/rest/getLyrics.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetLyrics)))
	router.Handle("/rest/getLyricsBySongId.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetLyricsBySongId)))
	router.Handle("/rest/getAvatar.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAvatar)))
	router.Handle("/rest/createAvatar.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreateAvatar)))
	router.Handle("/rest/updateAvatar.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUpdateAvatar)))
	router.Handle("/rest/deleteAvatar.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteAvatar)))
	// Media annotation
	router.Handle("/rest/star.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleStar)))
	router.Handle("/rest/unstar.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUnStar)))
	router.Handle("/rest/setRating.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSetRating)))
	router.Handle("/rest/scrobble.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleScrobble)))
	// Jukebox
	router.Handle("/rest/jukeboxControl.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleJukeboxControl)))
	// Chat
	router.Handle("/rest/getChatMessages.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetChatMessages)))
	router.Handle("/rest/addChatMessage.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleAddChatMessage)))
	// User management
	router.Handle("/rest/getUser.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetUser)))
	router.Handle("/rest/getUsers.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetUsers)))
	router.Handle("/rest/createUser.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreateUser)))
	router.Handle("/rest/updateUser.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUpdateUser)))
	router.Handle("/rest/deleteUser.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteUser)))
	router.Handle("/rest/changePassword.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleChangePassword)))
	// Media library scanning
	router.Handle("/rest/getScanStatus.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetScanStatus)))
	router.Handle("/rest/startScan.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleStartScan)))
	// server 404
	router.Handle("/rest/{unknownEndpoint}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleNotFound)))

	router.HandleFunc("/", handleFrontend)

	handler := cors.AllowAll().Handler(
		gziphandler.GzipHandler(router),
	)

	var serverAddress string
	if config.IsLocalDevEnv() {
		serverAddress = "localhost:8080"
		logger.Println("Application running at http://localhost:8080")
	} else {
		serverAddress = ":8080"
		logger.Println("Application running at :8080")
	}

	server := &http.Server{
		Addr:    serverAddress,
		Handler: handler,
	}
	return server
}

func handleFrontend(w http.ResponseWriter, r *http.Request) {
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
