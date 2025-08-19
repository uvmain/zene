package types

type Directory struct {
	Id            string          `xml:"id,attr" json:"id"`
	Parent        string          `xml:"parent,attr" json:"parent"`
	Name          string          `xml:"name,attr" json:"name"`
	Starred       string          `xml:"starred,attr" json:"starred"`
	UserRating    int             `xml:"userRating,attr" json:"userRating"`
	AverageRating float64         `xml:"averageRating,attr" json:"averageRating"`
	PlayCount     int             `xml:"playCount,attr" json:"playCount"`
	CoverArt      string          `xml:"coverArt,attr" json:"coverArt"`
	SongCount     int             `xml:"songCount,attr" json:"songCount"`
	Child         []SubsonicChild `xml:"child" json:"child"`
}

type SubsonicDirectory struct {
	Directory *Directory `xml:"directory,omitempty" json:"directory,omitempty"`
}
