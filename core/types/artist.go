package types

type Artist struct {
	Id             string          `xml:"id,attr" json:"id"`
	Name           string          `xml:"name,attr" json:"name"`
	CoverArt       string          `xml:"coverArt,attr,omitempty" json:"coverArt,omitempty"`
	ArtistImageUrl string          `xml:"artistImageUrl,attr,omitempty" json:"artistImageUrl,omitempty"`
	AlbumCount     int             `xml:"albumCount,attr" json:"albumCount"`
	Starred        string          `xml:"starred,attr,omitempty" json:"starred,omitempty"`
	Album          []SubsonicChild `xml:"album,omitempty" json:"album,omitempty"`
	MusicBrainzId  string          `xml:"musicBrainzId,attr,omitempty" json:"musicBrainzId,omitempty"`
	SortName       string          `xml:"sortName,attr,omitempty" json:"sortName,omitempty"`
	UserRating     int             `xml:"userRating,attr,omitempty" json:"userRating"`
	AverageRating  float64         `xml:"averageRating,attr,omitempty" json:"averageRating"`
}

type SubsonicArtistWrapper struct {
	Artist Artist `xml:"artist,omitempty" json:"artist,omitempty"`
}

type SubsonicArtistsWrapper struct {
	Artists         *[]Index `xml:"index,omitempty" json:"index,omitempty"`
	IgnoredArticles string   `xml:"ignoredArticles,attr" json:"ignoredArticles"`
}

type ArtistInfo struct {
	Biography      string   `xml:"biography,attr,omitempty" json:"biography,omitempty"`
	MusicBrainzId  string   `xml:"musicBrainzId,attr,omitempty" json:"musicBrainzId,omitempty"`
	LastFmUrl      string   `xml:"lastFmUrl,attr,omitempty" json:"lastFmUrl,omitempty"`
	SmallImageUrl  string   `xml:"smallImageUrl,attr,omitempty" json:"smallImageUrl,omitempty"`
	MediumImageUrl string   `xml:"mediumImageUrl,attr,omitempty" json:"mediumImageUrl,omitempty"`
	LargeImageUrl  string   `xml:"largeImageUrl,attr,omitempty" json:"largeImageUrl,omitempty"`
	SimilarArtists []Artist `xml:"similarArtists>artist,omitempty" json:"similarArtists,omitempty"`
}
