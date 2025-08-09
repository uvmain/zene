package types

import "encoding/xml"

type MusicFolders struct {
	MusicFolder []MusicFolder `xml:"musicFolder" json:"musicFolder"`
}

type MusicFolder struct {
	Id   int    `xml:"id,attr" json:"id"`
	Name string `xml:"name,attr" json:"name"`
}

type SubsonicMusicFolders struct {
	XMLName       xml.Name       `xml:"subsonic-response" json:"-"`
	Xmlns         string         `xml:"xmlns,attr" json:"-"`
	Status        string         `xml:"status,attr" json:"status"`
	Version       string         `xml:"version,attr" json:"version"`
	Type          string         `xml:"type,attr" json:"type"`
	ServerVersion string         `xml:"serverVersion,attr" json:"serverVersion"`
	OpenSubsonic  bool           `xml:"openSubsonic,attr" json:"openSubsonic"`
	Error         *SubsonicError `xml:"error,omitempty" json:"error,omitempty"`
	MusicFolders  *MusicFolders  `xml:"musicFolders" json:"musicFolders"`
}

type SubsonicMusicFoldersResponse struct {
	SubsonicResponse SubsonicMusicFolders `json:"subsonic-response"`
}
