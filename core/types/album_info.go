package types

type AlbumInfo struct {
	Notes          string `json:"notes,omitempty" xml:"notes,attr,omitempty"`
	MusicBrainzId  string `json:"musicBrainzId,omitempty" xml:"musicBrainzId,attr,omitempty"`
	LastFmUrl      string `json:"lastFmUrl,omitempty" xml:"lastFmUrl,attr,omitempty"`
	SmallImageUrl  string `json:"smallImageUrl,omitempty" xml:"smallImageUrl,attr,omitempty"`
	MediumImageUrl string `json:"mediumImageUrl,omitempty" xml:"mediumImageUrl,attr,omitempty"`
	LargeImageUrl  string `json:"largeImageUrl,omitempty" xml:"largeImageUrl,attr,omitempty"`
}
