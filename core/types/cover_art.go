package types

import "encoding/xml"

type CoverArt struct {
	Id   string `xml:"id,attr" json:"id"`
	Size int64  `xml:"size,attr" json:"size"`
}

type SubsonicCoverArt struct {
	XMLName       xml.Name       `xml:"subsonic-response" json:"-"`
	Xmlns         string         `xml:"xmlns,attr" json:"-"`
	Status        string         `xml:"status,attr" json:"status"`
	Version       string         `xml:"version,attr" json:"version"`
	Type          string         `xml:"type,attr" json:"type"`
	ServerVersion string         `xml:"serverVersion,attr" json:"serverVersion"`
	OpenSubsonic  bool           `xml:"openSubsonic,attr" json:"openSubsonic"`
	Error         *SubsonicError `xml:"error,omitempty" json:"error,omitempty"`
	CoverArt      *CoverArt      `xml:"coverArt" json:"coverArt"`
}

type SubsonicCoverArtResponse struct {
	SubsonicResponse SubsonicCoverArt `json:"subsonic-response"`
}
