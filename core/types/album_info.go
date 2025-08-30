package types

type AlbumInfo struct {
	Notes          string `json:"notes,omitempty" xml:"notes,omitempty"`
	MusicBrainzId  string `json:"musicBrainzId,omitempty" xml:"musicBrainzId,omitempty"`
	LastFmUrl      string `json:"lastFmUrl,omitempty" xml:"lastFmUrl,omitempty"`
	SmallImageUrl  string `json:"smallImageUrl,omitempty" xml:"smallImageUrl,omitempty"`
	MediumImageUrl string `json:"mediumImageUrl,omitempty" xml:"mediumImageUrl,omitempty"`
	LargeImageUrl  string `json:"largeImageUrl,omitempty" xml:"largeImageUrl,omitempty"`
}
