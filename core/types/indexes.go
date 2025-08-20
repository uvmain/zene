package types

type Shortcut struct {
	Id   string `xml:"id,attr" json:"id"`
	Name string `xml:"name,attr" json:"name"`
}

type Artist struct {
	Id             string  `xml:"id,attr" json:"id"`
	Name           string  `xml:"name,attr" json:"name"`
	CoverArt       string  `xml:"coverArt,attr" json:"coverArt"`
	ArtistImageUrl string  `xml:"artistImageUrl,attr" json:"artistImageUrl"`
	Starred        string  `xml:"starred,attr" json:"starred"`
	UserRating     int     `xml:"userRating,attr" json:"userRating"`
	AverageRating  float64 `xml:"averageRating,attr" json:"averageRating"`
}

type Index struct {
	Name   string   `xml:"name,attr" json:"name"`
	Artist []Artist `xml:"artist" json:"artist"`
}

type SubsonicIndexes struct {
	Indexes         *[]Index `xml:"index,omitempty" json:"index,omitempty"`
	LastModified    int      `xml:"lastModified,attr" json:"lastModified"`
	IgnoredArticles string   `xml:"ignoredArticles,attr" json:"ignoredArticles"`
}
