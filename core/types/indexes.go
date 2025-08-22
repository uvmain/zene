package types

type Shortcut struct {
	Id   string `xml:"id,attr" json:"id"`
	Name string `xml:"name,attr" json:"name"`
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
