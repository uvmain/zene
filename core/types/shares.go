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
