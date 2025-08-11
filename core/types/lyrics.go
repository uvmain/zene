package types

import "encoding/xml"

type SubsonicLyrics struct {
	Artist string `json:"artist" xml:"artist,attr"`
	Title  string `json:"title" xml:"title,attr"`
	Value  string `json:"value" xml:"value,attr"`
}

type SubsonicLyricsResponse struct {
	XMLName       xml.Name        `xml:"subsonic-response" json:"-"`
	Xmlns         string          `xml:"xmlns,attr" json:"-"`
	Status        string          `xml:"status,attr" json:"status"`
	Version       string          `xml:"version,attr" json:"version"`
	Type          string          `xml:"type,attr" json:"type"`
	ServerVersion string          `xml:"serverVersion,attr" json:"serverVersion"`
	OpenSubsonic  bool            `xml:"openSubsonic,attr" json:"openSubsonic"`
	Error         *SubsonicError  `xml:"error,omitempty" json:"error,omitempty"`
	Lyrics        *SubsonicLyrics `xml:"lyrics" json:"lyrics"`
}

type SubsonicLyricsResponseWrapper struct {
	SubsonicResponse SubsonicLyricsResponse `json:"subsonic-response"`
}

type StructuredLyricsLine struct {
	Start int    `xml:"start,attr" json:"start,omitempty"`
	Value string `xml:"value,attr" json:"value"`
}

type StructuredLyrics struct {
	DisplayArtist string                 `xml:"displayArtist,attr" json:"displayArtist"`
	DisplayTitle  string                 `xml:"displayTitle,attr" json:"displayTitle"`
	Lang          string                 `xml:"lang,attr" json:"lang"`
	Offset        int                    `xml:"offset,attr" json:"offset"`
	Synced        bool                   `xml:"synced,attr" json:"synced"`
	Line          []StructuredLyricsLine `xml:"line" json:"line"`
}

type SubsonicLyricsList struct {
	StructuredLyrics []StructuredLyrics `xml:"structuredLyrics,attr" json:"structuredLyrics"`
}

type SubsonicLyricsListResponse struct {
	XMLName       xml.Name            `xml:"subsonic-response" json:"-"`
	Xmlns         string              `xml:"xmlns,attr" json:"-"`
	Status        string              `xml:"status,attr" json:"status"`
	Version       string              `xml:"version,attr" json:"version"`
	Type          string              `xml:"type,attr" json:"type"`
	ServerVersion string              `xml:"serverVersion,attr" json:"serverVersion"`
	OpenSubsonic  bool                `xml:"openSubsonic,attr" json:"openSubsonic"`
	Error         *SubsonicError      `xml:"error,omitempty" json:"error,omitempty"`
	LyricsList    *SubsonicLyricsList `xml:"lyricsList" json:"lyricsList"`
}

type SubsonicLyricsListResponseWrapper struct {
	SubsonicResponse SubsonicLyricsListResponse `json:"subsonic-response"`
}

type LyricsDatabaseRow struct {
	MusicBrainzTrackID string `json:"musicBrainzTrackId"`
	PlainLyrics        string `json:"plainLyrics"`
	SyncedLyrics       string `json:"syncedLyrics"`
}
