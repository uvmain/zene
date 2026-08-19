package types

type ShareRow struct {
	Id          int             `json:"id" xml:"id,attr"`
	Url         string          `json:"url" xml:"url,attr"`
	Description string          `json:"description,omitempty" xml:"description,attr,omitempty"`
	Username    string          `json:"username" xml:"username,attr"`
	Created     string          `json:"created" xml:"created,attr"`
	VisitCount  int             `json:"visitCount" xml:"visit_count,attr"`
	Entries     []SubsonicChild `json:"entry,omitempty" xml:"entry,omitempty"`
}

type Shares struct {
	Share []ShareRow `json:"share" xml:"share"`
}

type DbUserShare struct {
	ShareId       int
	Token         string
	Description   string
	ShareCreated  string
	Expires       string
	VisitCount    int
	ShareOwner    string
	TrackId       string
	AlbumId       string
	Title         string
	Album         string
	Artist        string
	TrackNumber   int
	Year          int
	Genre         string
	CoverArt      string
	Size          int
	Duration      int
	Bitrate       int
	Path          string
	DateAdded     string
	DiscNumber    int
	ArtistId      string
	AlbumArtist   string
	BitDepth      int
	SampleRate    int
	Channels      int
	UserRating    int
	AverageRating float64
	PlayCount     int
	LastPlayed    string
	DateStarred   string
	AlbumArtistId string
}
