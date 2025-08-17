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

type Child struct {
	Id                    string `xml:"id,attr" json:"id"`
	Parent                string `xml:"parent,attr" json:"parent"`
	Title                 string `xml:"title,attr" json:"title"`
	IsDir                 string `xml:"isDir,attr" json:"isDir"`
	Album                 string `xml:"album,attr" json:"album"`
	Artist                string `xml:"artist,attr" json:"artist"`
	Track                 string `xml:"track,attr" json:"track"`
	Year                  string `xml:"year,attr" json:"year"`
	Genre                 string `xml:"genre,attr" json:"genre"`
	CoverArt              string `xml:"coverArt,attr" json:"coverArt"`
	Size                  string `xml:"size,attr" json:"size"`
	ContentType           string `xml:"contentType,attr" json:"contentType"`
	Suffix                string `xml:"suffix,attr" json:"suffix"`
	Duration              string `xml:"duration,attr" json:"duration"`
	BitRate               string `xml:"bitRate,attr" json:"bitRate"`
	Path                  string `xml:"path,attr" json:"path"`
	TranscodedContentType string `xml:"transcodedContentType,omitempty" json:"transcodedContentType,omitempty"`
	TranscodedSuffix      string `xml:"transcodedSuffix,omitempty" json:"transcodedSuffix,omitempty"`
}

type Index struct {
	Name   string   `xml:"name,attr" json:"name"`
	Artist []Artist `xml:"artist" json:"artist"`
}

// type Indexes struct {
// 	Index    *[]Index    `xml:"index,omitempty" json:"index,omitempty"`
// 	Child    *[]Child    `xml:"child,omitempty" json:"child,omitempty"`
// 	Shortcut *[]Shortcut `xml:"shortcut,omitempty" json:"shortcut,omitempty"`
// }

type SubsonicIndexes struct {
	Indexes         *[]Index `xml:"index,omitempty" json:"index,omitempty"`
	LastModified    int      `xml:"lastModified,attr" json:"lastModified"`
	IgnoredArticles string   `xml:"ignoredArticles,attr" json:"ignoredArticles"`
}
