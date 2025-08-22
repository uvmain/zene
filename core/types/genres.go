package types

type Genre struct {
	SongCount  int    `xml:"songCount,attr" json:"song_count"`
	AlbumCount int    `xml:"albumCount,attr" json:"album_count"`
	Value      string `xml:"value,attr" json:"value"`
}

type Genres struct {
	Genre []Genre `xml:"genre" json:"genre"`
}

type ItemGenre struct {
	Name string `xml:"name,attr" json:"name"`
}
