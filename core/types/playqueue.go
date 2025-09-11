package types

type Playqueue struct {
	Current   string          `json:"current" xml:"current,attr"`
	Position  int             `json:"position" xml:"position,attr"`
	Username  string          `json:"username" xml:"username,attr"`
	Changed   string          `json:"changed" xml:"changed,attr"`
	ChangedBy string          `json:"changedBy" xml:"changedBy,attr"`
	Entry     []SubsonicChild `json:"entry" xml:"entry"`
}

type PlayqueueByIndex struct {
	CurrentIndex int             `json:"currentIndex" xml:"currentIndex,attr"`
	Position     int             `json:"position" xml:"position,attr"`
	Username     string          `json:"username" xml:"username,attr"`
	Changed      string          `json:"changed" xml:"changed,attr"`
	ChangedBy    string          `json:"changedBy" xml:"changedBy,attr"`
	Entry        []SubsonicChild `json:"entry" xml:"entry"`
}

type PlayqueueRowParsed struct {
	Username     string
	Changed      string
	ChangedBy    string
	Position     int
	TrackIds     []string
	CurrentIndex int
}
