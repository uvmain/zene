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

	router.HandleFunc("/", handleFrontend)
	// authenticated routes
	router.Handle("GET /api/artists", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtists)))                                   // returns []types.ArtistResponse; query params: search=searchTerm, recent=true, random=false, limit=10, offset=10
	router.Handle("GET /api/artists/{musicBrainzArtistId}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtist)))              // returns types.ArtistResponse
	router.Handle("GET /api/artists/{musicBrainzArtistId}/tracks", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtistTracks))) // returns []types.MetadataWithPlaycounts; query params: recent=true, random=false, limit=10, offset=10
	router.Handle("GET /api/artists/{musicBrainzArtistId}/art", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtistArt)))       // returns image/jpeg blob
	router.Handle("GET /api/artists/{musicBrainzArtistId}/albums", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtistAlbums))) // returns []types.AlbumsResponse
	router.Handle("GET /api/albums", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbums)))                                     // returns []types.AlbumsResponse; query params: recent=true, random=false, limit=10
	router.Handle("GET /api/albums/{musicBrainzAlbumId}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbum)))                 // returns types.AlbumsResponse
	router.Handle("GET /api/albums/{musicBrainzAlbumId}/art", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbumArt)))          // returns image/jpeg blob
	router.Handle("GET /api/albums/{musicBrainzAlbumId}/tracks", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbumTracks)))    // returns []types.MetadataWithPlaycounts
	router.Handle("GET /api/tracks", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetTracks)))                                     // returns []types.Metadata; query params: recent=true, random=false, limit=10
	router.Handle("GET /api/tracks/{musicBrainzTrackId}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetTrack)))                 // returns types.MetadataWithPlaycounts
	router.Handle("GET /api/tracks/{musicBrainzTrackId}/lyrics", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetTrackLyrics)))    // returns types.Lyrics
	router.Handle("GET /api/tracks/{musicBrainzTrackId}/download", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDownloadTrack)))   // returns blob
	router.Handle("GET /api/tracks/{musicBrainzTrackId}/stream", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleStreamTrack)))       // returns blob range
	router.Handle("GET /api/genres", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetGenres)))                                     // query params: search=searchTerm
	router.Handle("GET /api/genres/tracks", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetTracksByGenre)))                       // query params: genres=genre1,genre2 condition=and|or
	router.Handle("GET /api/search", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSearchMetadata)))                                // query params: search=searchTerm
	router.Handle("GET /api/user", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetCurrentUser)))                                  // return types.User - current user
	router.Handle("GET /api/playcounts", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetPlaycounts)))                             // return []types.Playcount; query params: user_id=1, musicbrainz_track_id=musicBrainzTrackId
	router.Handle("POST /api/playcounts", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUpsertPlaycount)))                          // return handlers.StandardResponse; form body: user_id=1, musicbrainz_track_id=musicBrainzTrackId

	// admin routes
	router.Handle("POST /api/scan", auth.AuthMiddleware(http.HandlerFunc(handlers.HandlePostScan)))                   // triggers a scan of the music library if one is not already running
	router.Handle("GET /api/users", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAllUsers)))                // return []types.User - all users
	router.Handle("POST /api/users", auth.AuthMiddleware(http.HandlerFunc(handlers.HandlePostNewUser)))               // return userId int64
	router.Handle("GET /api/users/{userId}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetUserById)))       // return types.User - user by ID
	router.Handle("PATCH /api/users/{userId}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandlePatchUserById)))   // return userId int64
	router.Handle("DELETE /api/users/{userId}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteUserById))) // return { Status: string }

	// OpenSubsonic routes
	/// System
	router.Handle("/rest/ping.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandlePing)))                   // returns types.SubsonicResponse
	router.Handle("/rest/getLicense.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleLicense)))          // returns types.SubsonicLicenseResponse
	router.Handle("/rest/getOpenSubsonicExtensions.view", http.HandlerFunc(handlers.HandleOpenSubsonicExtensions)) // returns types.SubsonicOpenSubsonicExtensionsResponse
	router.Handle("/rest/tokenInfo.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleTokenInfo)))         // returns types.SubsonicTokenInfoResponse
	/// Browsing
	router.Handle("/rest/getMusicFolders.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetMusicFolders))) // returns types.SubsonicMusicFoldersResponse
	// Media retrieval
	router.Handle("/rest/getAvatar.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAvatar)))       // returns Image blob or types.SubsonicResponse error
	router.Handle("/rest/createAvatar.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreateAvatar))) // returns types.SubsonicResponse - not in the OpenSubsonic API spec
	router.Handle("/rest/updateAvatar.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUpdateAvatar))) // returns types.SubsonicResponse - not in the OpenSubsonic API spec
	router.Handle("/rest/deleteAvatar.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteAvatar))) // returns types.SubsonicResponse - not in the OpenSubsonic API spec
	// Chat
	router.Handle("/rest/getChatMessages.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetChatMessages))) // returns types.SubsonicChatMessagesResponse
	router.Handle("/rest/addChatMessage.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleAddChatMessage)))   // returns types.SubsonicResponse
	// User management
	router.Handle("/rest/getUser.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetUser)))               // returns types.SubsonicUserResponse
	router.Handle("/rest/getUsers.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetUsers)))             // returns types.SubsonicUsersResponse
	router.Handle("/rest/createUser.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreateUser)))         // returns types.SubsonicResponse
	router.Handle("/rest/updateUser.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUpdateUser)))         // returns types.SubsonicResponse
	router.Handle("/rest/deleteUser.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteUser)))         // returns types.SubsonicResponse
	router.Handle("/rest/changePassword.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleChangePassword))) // returns types.SubsonicResponse
	// Media library scanning
	router.Handle("/rest/getScanStatus.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetScanStatus))) // returns types.SubsonicScanStatusResponse
	router.Handle("/rest/startScan.view", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleStartScan)))         // returns types.SubsonicScanStatusResponse

	handler := cors.AllowAll().Handler(router)

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
