package types

type SearchResult2 struct {
	Artists []Artist        `xml:"artist" json:"artist"`
	Albums  []AlbumId3      `xml:"album" json:"album"`
	Songs   []SubsonicChild `xml:"song" json:"song"`
}

type SearchResult3 struct {
	Artists []Artist        `xml:"artist" json:"artist"`
	Albums  []AlbumId3      `xml:"album" json:"album"`
	Songs   []SubsonicChild `xml:"song" json:"song"`
}
