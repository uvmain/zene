package types

type MusicFolders struct {
	MusicFolder []MusicFolder `xml:"musicFolder" json:"musicFolder"`
}

type MusicFolder struct {
	Id   int    `xml:"id,attr" json:"id"`
	Name string `xml:"name,attr" json:"name"`
}
