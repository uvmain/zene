package types

type Genre struct {
	SongCount  int64  `xml:"songCount,attr" json:"song_count"`
	AlbumCount int64  `xml:"albumCount,attr" json:"album_count"`
	Value      string `xml:"value,attr" json:"value"`
}

type Genres struct {
	Genre []Genre `xml:"genre" json:"genre"`
}
