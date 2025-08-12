package types

import "encoding/xml"

type Genre struct {
	SongCount  int64  `xml:"songCount,attr" json:"song_count"`
	AlbumCount int64  `xml:"albumCount,attr" json:"album_count"`
	Value      string `xml:"value,attr" json:"value"`
}

type Genres struct {
	Genre []Genre `xml:"genre" json:"genre"`
}

type SubsonicGenres struct {
	XMLName       xml.Name       `xml:"subsonic-response" json:"-"`
	Xmlns         string         `xml:"xmlns,attr" json:"-"`
	Status        string         `xml:"status,attr" json:"status"`
	Version       string         `xml:"version,attr" json:"version"`
	Type          string         `xml:"type,attr" json:"type"`
	ServerVersion string         `xml:"serverVersion,attr" json:"serverVersion"`
	OpenSubsonic  bool           `xml:"openSubsonic,attr" json:"openSubsonic"`
	Error         *SubsonicError `xml:"error,omitempty" json:"error,omitempty"`
	Genres        *Genres        `xml:"genres" json:"genres"`
}

type SubsonicGenresResponse struct {
	SubsonicResponse SubsonicGenres `json:"subsonic-response"`
}
