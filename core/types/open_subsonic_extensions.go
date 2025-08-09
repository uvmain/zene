package types

import "encoding/xml"

type OpenSubsonicExtensions struct {
	Name     string `xml:"name" json:"name"`
	Versions []int  `xml:"versions" json:"versions"`
}

type SubsonicOpenSubsonicExtensions struct {
	XMLName                xml.Name                  `xml:"subsonic-response" json:"-"`
	Xmlns                  string                    `xml:"xmlns,attr" json:"-"`
	Status                 string                    `xml:"status,attr" json:"status"`
	Version                string                    `xml:"version,attr" json:"version"`
	Type                   string                    `xml:"type,attr" json:"type"`
	ServerVersion          string                    `xml:"serverVersion,attr" json:"serverVersion"`
	OpenSubsonic           bool                      `xml:"openSubsonic,attr" json:"openSubsonic"`
	Error                  *SubsonicError            `xml:"error,omitempty" json:"error,omitempty"`
	OpenSubsonicExtensions []*OpenSubsonicExtensions `xml:"openSubsonicExtensions" json:"openSubsonicExtensions"`
}

type SubsonicOpenSubsonicExtensionsResponse struct {
	SubsonicResponse SubsonicOpenSubsonicExtensions `json:"subsonic-response"`
}
