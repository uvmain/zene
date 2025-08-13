package types

type CoverArt struct {
	Id   string `xml:"id,attr" json:"id"`
	Size int64  `xml:"size,attr" json:"size"`
}
