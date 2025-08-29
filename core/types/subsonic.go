package types

import (
	"encoding/xml"
)

type SubsonicStandard struct {
	XMLName                xml.Name                  `xml:"subsonic-response" json:"-"`
	Xmlns                  string                    `xml:"xmlns,attr" json:"-"`
	Status                 string                    `xml:"status,attr" json:"status"`
	Version                string                    `xml:"version,attr" json:"version"`
	Type                   string                    `xml:"type,attr" json:"type"`
	ServerVersion          string                    `xml:"serverVersion,attr" json:"serverVersion"`
	OpenSubsonic           bool                      `xml:"openSubsonic,attr" json:"openSubsonic"`
	Error                  *SubsonicError            `xml:"error,omitempty" json:"error,omitempty"`
	License                *LicenseInfo              `xml:"license,omitempty" json:"license,omitempty"`
	ScanStatus             *ScanStatus               `xml:"scanStatus,omitempty" json:"scanStatus,omitempty"`
	OpenSubsonicExtensions []*OpenSubsonicExtensions `xml:"openSubsonicExtensions,omitempty" json:"openSubsonicExtensions,omitempty"`
	User                   *SubsonicUser             `xml:"user,omitempty" json:"user,omitempty"`
	Users                  *SubsonicUsers            `xml:"users,omitempty" json:"users,omitempty"`
	Lyrics                 *SubsonicLyrics           `xml:"lyrics,omitempty" json:"lyrics,omitempty"`
	LyricsList             *SubsonicLyricsList       `xml:"lyricsList,omitempty" json:"lyricsList,omitempty"`
	MusicFolders           *MusicFolders             `xml:"musicFolders,omitempty" json:"musicFolders,omitempty"`
	Genres                 *Genres                   `xml:"genres,omitempty" json:"genres,omitempty"`
	CoverArt               *CoverArt                 `xml:"coverArt,omitempty" json:"coverArt,omitempty"`
	ChatMessages           *ChatMessages             `xml:"chatMessages,omitempty" json:"chatMessages,omitempty"`
	TokenInfo              *TokenInfo                `xml:"tokenInfo,omitempty" json:"tokenInfo,omitempty"`
	Indexes                *SubsonicIndexes          `xml:"indexes,omitempty" json:"indexes,omitempty"`
	Song                   *SubsonicChild            `xml:"song,omitempty" json:"song,omitempty"`
	RandomSongs            *RandomSongs              `xml:"randomSongs,omitempty" json:"randomSongs,omitempty"`
	Starred                *Starred                  `xml:"starred,omitempty" json:"starred,omitempty"`
	Starred2               *Starred2                 `xml:"starred2,omitempty" json:"starred2,omitempty"`
	SongsByGenre           *SongsByGenre             `xml:"songsByGenre,omitempty" json:"songsByGenre,omitempty"`
	MusicDirectory         *SubsonicDirectory        `xml:"directory,omitempty" json:"directory,omitempty"`
	NowPlaying             *SubsonicNowPlaying       `xml:"nowPlaying,omitempty" json:"nowPlaying,omitempty"`
	Artist                 *Artist                   `xml:"artist,omitempty" json:"artist,omitempty"`
	Artists                *SubsonicArtistsWrapper   `xml:"artists,omitempty" json:"artists,omitempty"`
	Album                  *AlbumId3                 `xml:"album,omitempty" json:"album,omitempty"`
	ArtistInfo             *ArtistInfo               `xml:"artistInfo,omitempty" json:"artistInfo,omitempty"`
	ArtistInfo2            *ArtistInfo               `xml:"artistInfo2,omitempty" json:"artistInfo2,omitempty"`
	AlbumInfo              *AlbumInfo                `xml:"albumInfo,omitempty" json:"albumInfo,omitempty"`
	AlbumInfo2             *AlbumInfo                `xml:"albumInfo2,omitempty" json:"albumInfo2,omitempty"`
	TopSongs               *TopSongs                 `xml:"topSongs,omitempty" json:"topSongs,omitempty"`
	SearchResult2          *SearchResult2            `xml:"searchResult2,omitempty" json:"searchResult2,omitempty"`
	SearchResult3          *SearchResult3            `xml:"searchResult3,omitempty" json:"searchResult3,omitempty"`
	AlbumList              *AlbumList                `xml:"albumList,omitempty" json:"albumList,omitempty"`
	AlbumList2             *AlbumList2               `xml:"albumList2,omitempty" json:"albumList2,omitempty"`
	SimilarSongs           *SimilarSongs             `xml:"similarSongs,omitempty" json:"similarSongs,omitempty"`
	SimilarSongs2          *SimilarSongs2            `xml:"similarSongs2,omitempty" json:"similarSongs2,omitempty"`
	Playlist               *PlaylistRow              `xml:"playlist,omitempty" json:"playlist,omitempty"`
	Playlists              *Playlists                `xml:"playlists,omitempty" json:"playlists,omitempty"`
}

type SubsonicResponse struct {
	SubsonicResponse SubsonicStandard `json:"subsonic-response"`
}

// SubsonicError represents a Subsonic API error
type SubsonicError struct {
	Code    int    `xml:"code,attr" json:"code"`
	Message string `xml:"message,attr" json:"message"`
	HelpUrl string `xml:"helpUrl,attr,omitempty" json:"helpUrl,omitempty"`
}

// Subsonic error codes as defined in the OpenSubsonic API specification
const (
	ErrorGeneric                   = 0
	ErrorMissingParameter          = 10
	ErrorIncompatibleVersion       = 20
	ErrorIncompatibleClient        = 30
	ErrorWrongCredentials          = 40
	ErrorTokenAuthNotSupported     = 41
	ErrorAuthMechanismNotSupported = 42
	ErrorTooManyAuthMechanisms     = 43
	ErrorInvalidApiKey             = 44
	ErrorNotAuthorized             = 50
	ErrorTrialExpired              = 60
	ErrorDataNotFound              = 70
)

type LicenseInfo struct {
	Valid          bool   `xml:"valid,attr" json:"valid"`
	Email          string `xml:"email,attr,omitempty" json:"email,omitempty"`
	LicenseExpires string `xml:"licenseExpires,attr,omitempty" json:"licenseExpires,omitempty"`
	TrialExpires   string `xml:"trialExpires,attr,omitempty" json:"trialExpires,omitempty"`
}

type TokenInfo struct {
	Username string `xml:"username" json:"username"`
}

type ItemDate struct {
	Year  int `xml:"year,attr,omitempty" json:"year,omitempty"`
	Month int `xml:"month,attr,omitempty" json:"month,omitempty"`
	Day   int `xml:"day,attr,omitempty" json:"day,omitempty"`
}

type TopSongs struct {
	Songs []SubsonicChild `xml:"song" json:"song"`
}

type Starred struct {
	Artists []Artist        `xml:"artist" json:"artist"`
	Albums  []AlbumId3      `xml:"album" json:"album"`
	Songs   []SubsonicChild `xml:"song" json:"song"`
}

type Starred2 struct {
	Artists []Artist        `xml:"artist" json:"artist"`
	Albums  []AlbumId3      `xml:"album" json:"album"`
	Songs   []SubsonicChild `xml:"song" json:"song"`
}

type SimilarSongs struct {
	Songs []SubsonicChild `xml:"song" json:"song"`
}

type SimilarSongs2 struct {
	Songs []SubsonicChild `xml:"song" json:"song"`
}
