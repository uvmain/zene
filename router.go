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

type CaseInsensitiveMux struct {
	mux *http.ServeMux
}

func (c *CaseInsensitiveMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// lowercase the url.path for case-insensitive matching
	r.URL.Path = strings.ToLower(r.URL.Path)
	// if path ends with .view, remove the .view suffix for matching
	r.URL.Path = strings.TrimSuffix(r.URL.Path, ".view")
	c.mux.ServeHTTP(w, r)
}

func NewCaseInsensitiveMux() *CaseInsensitiveMux {
	return &CaseInsensitiveMux{mux: http.NewServeMux()}
}

func (c *CaseInsensitiveMux) Handle(pattern string, handler http.Handler) {
	// normalize the registered pattern
	c.mux.Handle(strings.ToLower(pattern), handler)
}

func (c *CaseInsensitiveMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	// normalize the registered pattern
	c.mux.HandleFunc(strings.ToLower(pattern), handler)
}

func StartServer() *http.Server {

	router := NewCaseInsensitiveMux()
	// all registered paths should be lowercase

	distSubFS, err = fs.Sub(dist, "frontend/dist")
	if err != nil {
		log.Fatal("Failed to create sub filesystem:", err)
	}

	// unauthenticated routes
	router.Handle("GET /share/img/{imageId}", http.HandlerFunc(handlers.HandleGetShareImg))

	/* cSpell:disable */

	// OpenSubsonic routes
	/// System
	router.Handle("/rest/ping", http.HandlerFunc(handlers.HandlePing))
	router.Handle("/rest/getlicense", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleLicense)))
	router.Handle("/rest/getopensubsonicextensions", http.HandlerFunc(handlers.HandleOpenSubsonicExtensions))
	router.Handle("/rest/tokeninfo", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleTokenInfo)))
	/// Browsing
	router.Handle("/rest/getmusicfolders", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetMusicFolders)))
	router.Handle("/rest/getindexes", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetIndexes)))
	router.Handle("/rest/getmusicdirectory", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetMusicDirectory)))
	router.Handle("/rest/getgenres", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetGenres)))
	router.Handle("/rest/getartists", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtists)))
	router.Handle("/rest/getartist", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtist)))
	router.Handle("/rest/getalbum", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbum)))
	router.Handle("/rest/getsong", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetSong)))
	router.Handle("/rest/getvideos", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetVideos)))
	router.Handle("/rest/getvideoinfo", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetVideoInfo)))
	router.Handle("/rest/getartistinfo", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtistInfo)))
	router.Handle("/rest/getartistinfo2", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtistInfo)))
	router.Handle("/rest/getalbuminfo", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbumInfo)))
	router.Handle("/rest/getalbuminfo2", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbumInfo)))
	router.Handle("/rest/getsimilarsongs", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetSimilarSongs)))
	router.Handle("/rest/getsimilarsongs2", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetSimilarSongs)))
	router.Handle("/rest/gettopsongs", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetTopSongs)))
	// Album/song lists
	router.Handle("/rest/getalbumlist", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbumList)))
	router.Handle("/rest/getalbumlist2", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbumList)))
	router.Handle("/rest/getrandomsongs", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetRandomSongs)))
	router.Handle("/rest/getsongsbygenre", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetSongsByGenre)))
	router.Handle("/rest/getnowplaying", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetNowPlaying)))
	router.Handle("/rest/getstarred", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetStarred)))
	router.Handle("/rest/getstarred2", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetStarred)))
	// Searching
	router.Handle("/rest/search", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSearch)))
	router.Handle("/rest/search2", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSearch)))
	router.Handle("/rest/search3", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSearch)))
	// Playlists
	router.Handle("/rest/getplaylists", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetPlaylists)))
	router.Handle("/rest/getplaylist", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetPlaylist)))
	router.Handle("/rest/createplaylist", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreatePlaylist)))
	router.Handle("/rest/updateplaylist", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUpdatePlaylist)))
	router.Handle("/rest/deleteplaylist", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeletePlaylist)))
	// Media retrieval
	router.Handle("/rest/stream", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleStream)))
	router.Handle("/rest/download", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDownload)))
	router.Handle("/rest/getcaptions", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetCaptions)))
	router.Handle("/rest/getcoverart", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetCoverArt)))
	router.Handle("/rest/getlyrics", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetLyrics)))
	router.Handle("/rest/getlyricsbysongid", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetLyricsBySongId)))
	router.Handle("/rest/getavatar", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAvatar)))
	router.Handle("/rest/createavatar", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreateAvatar)))
	router.Handle("/rest/updateavatar", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUpdateAvatar)))
	router.Handle("/rest/deleteavatar", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteAvatar)))
	// Media annotation
	router.Handle("/rest/star", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleStar)))
	router.Handle("/rest/unstar", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUnStar)))
	router.Handle("/rest/setrating", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSetRating)))
	router.Handle("/rest/scrobble", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleScrobble)))
	// Jukebox
	router.Handle("/rest/jukeboxcontrol", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleJukeboxControl)))
	// Internet radio
	router.Handle("/rest/getinternetradiostations", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetInternetRadioStations)))
	router.Handle("/rest/createinternetradiostation", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreateInternetRadioStation)))
	router.Handle("/rest/updateinternetradiostation", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUpdateInternetRadioStation)))
	router.Handle("/rest/deleteinternetradiostation", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteInternetRadioStation)))
	// Chat
	router.Handle("/rest/getchatmessages", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetChatMessages)))
	router.Handle("/rest/addchatmessage", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleAddChatMessage)))
	// User management
	router.Handle("/rest/getuser", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetUser)))
	router.Handle("/rest/getusers", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetUsers)))
	router.Handle("/rest/createuser", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreateUser)))
	router.Handle("/rest/updateuser", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUpdateUser)))
	router.Handle("/rest/deleteuser", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteUser)))
	router.Handle("/rest/changepassword", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleChangePassword)))
	router.Handle("/rest/createapikey", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreateApiKey)))
	router.Handle("/rest/getapikeys", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetApiKeys)))
	router.Handle("/rest/deleteapikey", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteApiKey)))
	// Bookmarks
	router.Handle("/rest/createbookmark", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreateBookmark)))
	router.Handle("/rest/getbookmarks", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetBookmarks)))
	router.Handle("/rest/deletebookmark", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteBookmark)))
	router.Handle("/rest/saveplayqueue", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSaveOrClearPlayqueue)))
	router.Handle("/rest/saveplayqueuebyindex", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSaveOrClearPlayqueue)))
	router.Handle("/rest/getplayqueuebyindex", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetPlayqueueByIndex)))
	// Media library scanning
	router.Handle("/rest/getscanstatus", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetScanStatus)))
	router.Handle("/rest/startscan", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleStartScan)))
	// server 404
	router.Handle("/rest/{unknownEndpoint}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleNotFound)))

	/* cSpell:enable */

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
