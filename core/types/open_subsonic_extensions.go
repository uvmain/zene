package types

type OpenSubsonicExtensions struct {
	Name     string `xml:"name" json:"name"`
	Versions []int  `xml:"versions" json:"versions"`
}
