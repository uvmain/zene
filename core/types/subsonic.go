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
	Song                   *SubsonicSong             `xml:"song,omitempty" json:"song,omitempty"`
	SongsByGenre           *SongsByGenre             `xml:"songsByGenre,omitempty" json:"songsByGenre,omitempty"`
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
