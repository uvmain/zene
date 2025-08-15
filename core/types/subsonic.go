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
}

type SubsonicResponse struct {
	SubsonicResponse SubsonicStandard `json:"subsonic-response"`
}

// SubsonicError represents a Subsonic API error
type SubsonicError struct {
	Code    int64  `xml:"code,attr" json:"code"`
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

type SubsonicArtist struct {
	Id             string   `xml:"id,attr" json:"id"`
	Name           string   `xml:"name,attr" json:"name"`
	CoverArt       string   `xml:"coverArt,attr" json:"coverArt"`
	AlbumCount     int64    `xml:"albumCount,attr" json:"albumCount"`
	ArtistImageUrl string   `xml:"artistImageUrl,attr,omitempty" json:"artistImageUrl,omitempty"`
	Starred        string   `xml:"starred,attr" json:"starred"`
	UserRating     int64    `xml:"userRating,attr" json:"userRating"`
	AverageRating  float64  `xml:"averageRating,attr" json:"averageRating"`
	MusicBrainzId  string   `xml:"musicBrainzId,attr" json:"musicBrainzId"`
	Roles          []string `xml:"roles>role" json:"roles"`
}

type SubsonicAlbumID3 struct {
	ID           string `xml:"id,attr" json:"id"`
	Name         string `xml:"name,attr" json:"name"`
	Version      string `xml:"version,attr" json:"version"`
	Artist       string `xml:"artist,attr" json:"artist"`
	Year         int64  `xml:"year,attr" json:"year"`
	CoverArt     string `xml:"coverArt,attr" json:"coverArt"`
	Starred      string `xml:"starred,attr" json:"starred"`
	Duration     int64  `xml:"duration,attr" json:"duration"`
	PlayCount    int64  `xml:"playCount,attr" json:"playCount"`
	Genre        string `xml:"genre,attr" json:"genre"`
	Created      string `xml:"created,attr" json:"created"`
	ArtistID     string `xml:"artistId,attr" json:"artistId"`
	SongCount    int64  `xml:"songCount,attr" json:"songCount"`
	Played       string `xml:"played,attr" json:"played"`
	UserRating   int64  `xml:"userRating,attr" json:"userRating"`
	RecordLabels []struct {
		Name string `xml:"name,attr" json:"name"`
	} `xml:"recordLabels" json:"recordLabels"`
	MusicBrainzID string `xml:"musicBrainzId,attr" json:"musicBrainzId"`
	Genres        []struct {
		Name string `xml:"name,attr" json:"name"`
	} `xml:"genres" json:"genres"`
	Artists []struct {
		ID   string `xml:"id,attr" json:"id"`
		Name string `xml:"name,attr" json:"name"`
	} `xml:"artists" json:"artists"`
	DisplayArtist       string   `xml:"displayArtist,attr" json:"displayArtist"`
	ReleaseTypes        []string `xml:"releaseTypes,attr" json:"releaseTypes"`
	Moods               []string `xml:"moods,attr" json:"moods"`
	SortName            string   `xml:"sortName,attr" json:"sortName"`
	OriginalReleaseDate struct {
		Year  int64 `xml:"year,attr" json:"year"`
		Month int64 `xml:"month,attr" json:"month"`
		Day   int64 `xml:"day,attr" json:"day"`
	} `xml:"originalReleaseDate" json:"originalReleaseDate"`
	ReleaseDate struct {
		Year  int64 `xml:"year,attr" json:"year"`
		Month int64 `xml:"month,attr" json:"month"`
		Day   int64 `xml:"day,attr" json:"day"`
	} `xml:"releaseDate" json:"releaseDate"`
	IsCompilation  bool   `xml:"isCompilation,attr" json:"isCompilation"`
	ExplicitStatus string `xml:"explicitStatus,attr" json:"explicitStatus"`
	DiscTitles     []struct {
		Disc  int64  `xml:"disc,attr" json:"disc"`
		Title string `xml:"title,attr" json:"title"`
	} `xml:"discTitles" json:"discTitles"`
}

type SubsonicSong struct {
	ID            string   `xml:"id,attr" json:"id"`
	Parent        string   `xml:"parent,attr" json:"parent"`
	IsDir         bool     `xml:"isDir,attr" json:"isDir"`
	Title         string   `xml:"title,attr" json:"title"`
	Album         string   `xml:"album,attr" json:"album"`
	Artist        string   `xml:"artist,attr" json:"artist"`
	Track         int64    `xml:"track,attr" json:"track"`
	Year          int64    `xml:"year,attr" json:"year"`
	Genre         string   `xml:"genre,attr" json:"genre"`
	CoverArt      string   `xml:"coverArt,attr" json:"coverArt"`
	Size          int64    `xml:"size,attr" json:"size"`
	ContentType   string   `xml:"contentType,attr" json:"contentType"`
	Suffix        string   `xml:"suffix,attr" json:"suffix"`
	Starred       string   `xml:"starred,attr" json:"starred"`
	Duration      int64    `xml:"duration,attr" json:"duration"`
	BitRate       int64    `xml:"bitRate,attr" json:"bitRate"`
	Path          string   `xml:"path,attr" json:"path"`
	DiscNumber    int64    `xml:"discNumber,attr" json:"discNumber"`
	Created       string   `xml:"created,attr" json:"created"`
	AlbumID       string   `xml:"albumId,attr" json:"albumId"`
	ArtistID      string   `xml:"artistId,attr" json:"artistId"`
	Type          string   `xml:"type,attr" json:"type"`
	IsVideo       bool     `xml:"isVideo,attr" json:"isVideo"`
	Bpm           int64    `xml:"bpm,attr" json:"bpm"`
	Comment       string   `xml:"comment,attr" json:"comment"`
	SortName      string   `xml:"sortName,attr" json:"sortName"`
	MediaType     string   `xml:"mediaType,attr" json:"mediaType"`
	MusicBrainzID string   `xml:"musicBrainzId,attr" json:"musicBrainzId"`
	Isrc          []string `xml:"isrc,attr" json:"isrc"`
	Genres        []struct {
		Name string `xml:"name,attr" json:"name"`
	} `xml:"genres" json:"genres"`
	ReplayGain struct {
		TrackGain int64 `xml:"trackGain,attr" json:"trackGain"`
		AlbumGain int64 `xml:"albumGain,attr" json:"albumGain"`
		TrackPeak int64 `xml:"trackPeak,attr" json:"trackPeak"`
		AlbumPeak int64 `xml:"albumPeak,attr" json:"albumPeak"`
	} `xml:"replayGain" json:"replayGain"`
	ChannelCount int64    `xml:"channelCount,attr" json:"channelCount"`
	SamplingRate int64    `xml:"samplingRate,attr" json:"samplingRate"`
	BitDepth     int64    `xml:"bitDepth,attr" json:"bitDepth"`
	Moods        []string `xml:"moods" json:"moods"`
	Artists      []struct {
		ID   string `xml:"id,attr" json:"id"`
		Name string `xml:"name,attr" json:"name"`
	} `xml:"artists" json:"artists"`
	DisplayArtist string `xml:"displayArtist,attr" json:"displayArtist"`
	AlbumArtists  []struct {
		ID   string `xml:"id,attr" json:"id"`
		Name string `xml:"name,attr" json:"name"`
	} `xml:"albumArtists" json:"albumArtists"`
	DisplayAlbumArtist string                `xml:"displayAlbumArtist,attr" json:"displayAlbumArtist"`
	Contributors       []SubsonicContributor `xml:"contributors" json:"contributors"`
	DisplayComposer    string                `xml:"displayComposer,attr" json:"displayComposer"`
	ExplicitStatus     string                `xml:"explicitStatus,attr" json:"explicitStatus"`
}

type SubsonicContributor struct {
	Role   string `xml:"role,attr" json:"role"`
	Artist struct {
		ID   string `xml:"id,attr" json:"id"`
		Name string `xml:"name,attr" json:"name"`
	} `xml:"artist" json:"artist"`
}
