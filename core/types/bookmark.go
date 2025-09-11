package types

type Bookmark struct {
	Entry    SubsonicChild `json:"entry" xml:"entry"`
	Position int           `json:"position" xml:"position,attr"`
	Username string        `json:"username" xml:"username,attr"`
	Comment  string        `json:"comment" xml:"comment,attr"`
	Created  string        `json:"created" xml:"created,attr"`
	Changed  string        `json:"changed" xml:"changed,attr"`
}

type Bookmarks struct {
	Bookmarks []Bookmark `json:"bookmark" xml:"bookmark"`
}
