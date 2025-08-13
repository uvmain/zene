package types

type SubsonicLyrics struct {
	Artist string `json:"artist" xml:"artist,attr"`
	Title  string `json:"title" xml:"title,attr"`
	Value  string `json:"value" xml:"value,attr"`
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

type LyricsDatabaseRow struct {
	MusicBrainzTrackID string `json:"musicBrainzTrackId"`
	PlainLyrics        string `json:"plainLyrics"`
	SyncedLyrics       string `json:"syncedLyrics"`
}
