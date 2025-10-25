package main

import (
	"embed"
	"io"
	"io/fs"
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

	"github.com/go-swiss/compress"
	"github.com/rs/cors"
)

//go:embed all:frontend/dist
var dist embed.FS
var distSubFS fs.FS

func init() {
	var err error
	distSubFS, err = fs.Sub(dist, "frontend/dist")
	if err != nil {
		panic("Failed to create sub filesystem: " + err.Error())
	}
}

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
	// lowercase the registered pattern
	c.mux.Handle(strings.ToLower(pattern), handler)
}

func (c *CaseInsensitiveMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	// lowercase the registered pattern
	c.mux.HandleFunc(strings.ToLower(pattern), handler)
}

func StartServer() *http.Server {

	// API router (case-insensitive)
	apiRouter := NewCaseInsensitiveMux()
	// all registered API paths should be lowercase
	apiRouter.Handle("/share/img/{image_id}", http.HandlerFunc(handlers.HandleGetShareImg))
	apiRouter.Handle("/rest/getalbumarts", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbumArts)))
	apiRouter.Handle("/rest/updatealbumart", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUpdateAlbumArt)))
	apiRouter.Handle("/rest/getartistlist", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtistList)))
	/* cSpell:disable */

	// System
	apiRouter.Handle("/rest/ping", http.HandlerFunc(handlers.HandlePing))
	apiRouter.Handle("/rest/getlicense", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleLicense)))
	apiRouter.Handle("/rest/getopensubsonicextensions", http.HandlerFunc(handlers.HandleOpenSubsonicExtensions))
	apiRouter.Handle("/rest/tokeninfo", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleTokenInfo)))
	// Browsing
	apiRouter.Handle("/rest/getmusicfolders", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetMusicFolders)))
	apiRouter.Handle("/rest/getindexes", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetIndexes)))
	apiRouter.Handle("/rest/getmusicdirectory", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetMusicDirectory)))
	apiRouter.Handle("/rest/getgenres", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetGenres)))
	apiRouter.Handle("/rest/getartists", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtists)))
	apiRouter.Handle("/rest/getartist", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtist)))
	apiRouter.Handle("/rest/getalbum", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbum)))
	apiRouter.Handle("/rest/getsong", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetSong)))
	apiRouter.Handle("/rest/getvideos", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetVideos)))
	apiRouter.Handle("/rest/getvideoinfo", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetVideoInfo)))
	apiRouter.Handle("/rest/getartistinfo", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtistInfo)))
	apiRouter.Handle("/rest/getartistinfo2", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetArtistInfo)))
	apiRouter.Handle("/rest/getalbuminfo", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbumInfo)))
	apiRouter.Handle("/rest/getalbuminfo2", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbumInfo)))
	apiRouter.Handle("/rest/getsimilarsongs", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetSimilarSongs)))
	apiRouter.Handle("/rest/getsimilarsongs2", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetSimilarSongs)))
	apiRouter.Handle("/rest/gettopsongs", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetTopSongs)))
	// Album/song lists
	apiRouter.Handle("/rest/getalbumlist", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbumList)))
	apiRouter.Handle("/rest/getalbumlist2", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAlbumList)))
	apiRouter.Handle("/rest/getrandomsongs", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetRandomSongs)))
	apiRouter.Handle("/rest/getsongsbygenre", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetSongsByGenre)))
	apiRouter.Handle("/rest/getnowplaying", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetNowPlaying)))
	apiRouter.Handle("/rest/getstarred", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetStarred)))
	apiRouter.Handle("/rest/getstarred2", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetStarred)))
	// Searching
	apiRouter.Handle("/rest/search", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSearch)))
	apiRouter.Handle("/rest/search2", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSearch)))
	apiRouter.Handle("/rest/search3", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSearch)))
	// Playlists
	apiRouter.Handle("/rest/getplaylists", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetPlaylists)))
	apiRouter.Handle("/rest/getplaylist", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetPlaylist)))
	apiRouter.Handle("/rest/createplaylist", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreatePlaylist)))
	apiRouter.Handle("/rest/updateplaylist", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUpdatePlaylist)))
	apiRouter.Handle("/rest/deleteplaylist", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeletePlaylist)))
	// Media retrieval
	apiRouter.Handle("/rest/stream", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleStream)))
	apiRouter.Handle("/rest/download", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDownload)))
	apiRouter.Handle("/rest/getcaptions", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetCaptions)))
	apiRouter.Handle("/rest/getcoverart", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetCoverArt)))
	apiRouter.Handle("/rest/getlyrics", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetLyrics)))
	apiRouter.Handle("/rest/getlyricsbysongid", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetLyricsBySongId)))
	apiRouter.Handle("/rest/getavatar", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetAvatar)))
	apiRouter.Handle("/rest/createavatar", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreateAvatar)))
	apiRouter.Handle("/rest/updateavatar", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUpdateAvatar)))
	apiRouter.Handle("/rest/deleteavatar", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteAvatar)))
	// Media annotation
	apiRouter.Handle("/rest/star", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleStar)))
	apiRouter.Handle("/rest/unstar", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUnStar)))
	apiRouter.Handle("/rest/setrating", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSetRating)))
	apiRouter.Handle("/rest/scrobble", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleScrobble)))
	// Sharing
	// Podcast
	apiRouter.Handle("/rest/getpodcasts", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetPodcasts)))
	apiRouter.Handle("/rest/getnewestpodcasts", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetNewestPodcasts)))
	apiRouter.Handle("/rest/getpodcastepisode", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetPodcastEpisode)))
	apiRouter.Handle("/rest/refreshpodcast", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleRefreshPodcast)))
	apiRouter.Handle("/rest/refreshpodcasts", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleRefreshPodcasts)))
	apiRouter.Handle("/rest/createpodcastchannel", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreatePodcastChannel)))
	apiRouter.Handle("/rest/deletepodcastchannel", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeletePodcastChannel)))
	apiRouter.Handle("/rest/deletepodcastepisode", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeletePodcastEpisode)))
	apiRouter.Handle("/rest/downloadpodcastepisode", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDownloadPodcastEpisode)))
	// Jukebox
	apiRouter.Handle("/rest/jukeboxcontrol", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleJukeboxControl)))
	// Internet radio
	apiRouter.Handle("/rest/getinternetradiostations", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetInternetRadioStations)))
	apiRouter.Handle("/rest/createinternetradiostation", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreateInternetRadioStation)))
	apiRouter.Handle("/rest/updateinternetradiostation", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUpdateInternetRadioStation)))
	apiRouter.Handle("/rest/deleteinternetradiostation", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteInternetRadioStation)))
	// Chat
	apiRouter.Handle("/rest/getchatmessages", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetChatMessages)))
	apiRouter.Handle("/rest/addchatmessage", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleAddChatMessage)))
	// User management
	apiRouter.Handle("/rest/getuser", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetUser)))
	apiRouter.Handle("/rest/getusers", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetUsers)))
	apiRouter.Handle("/rest/createuser", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreateUser)))
	apiRouter.Handle("/rest/updateuser", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleUpdateUser)))
	apiRouter.Handle("/rest/deleteuser", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteUser)))
	apiRouter.Handle("/rest/changepassword", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleChangePassword)))
	apiRouter.Handle("/rest/createapikey", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreateApiKey)))
	apiRouter.Handle("/rest/getapikeys", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetApiKeys)))
	apiRouter.Handle("/rest/deleteapikey", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteApiKey)))
	// Bookmarks
	apiRouter.Handle("/rest/createbookmark", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCreateBookmark)))
	apiRouter.Handle("/rest/getbookmarks", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetBookmarks)))
	apiRouter.Handle("/rest/deletebookmark", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteBookmark)))
	apiRouter.Handle("/rest/saveplayqueue", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSaveOrClearPlayqueue)))
	apiRouter.Handle("/rest/saveplayqueuebyindex", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleSaveOrClearPlayqueue)))
	apiRouter.Handle("/rest/getplayqueue", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetPlayqueue)))
	apiRouter.Handle("/rest/getplayqueuebyindex", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetPlayqueueByIndex)))
	apiRouter.Handle("/rest/getscanstatus", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleGetScanStatus)))
	// Media library scanning
	apiRouter.Handle("/rest/startscan", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleStartScan)))
	apiRouter.Handle("/rest/{unknownEndpoint}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleNotFound)))
	/* cSpell:enable */

	// frontend router (case-sensitive, for static files and SPA)
	frontendMux := http.NewServeMux()
	frontendMux.HandleFunc("/", handleFrontend)

	// main handler
	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// backend routes (case-insensitive)
		if strings.HasPrefix(strings.ToLower(r.URL.Path), "/rest/") || strings.HasPrefix(strings.ToLower(r.URL.Path), "/share/") {
			apiRouter.ServeHTTP(w, r)
			return
		}
		// frontend routes
		frontendMux.ServeHTTP(w, r)
	})

	handler := cors.AllowAll().Handler(
		compress.Middleware(mainHandler),
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
