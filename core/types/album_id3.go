package types

type AlbumId3 struct {
	Id            string             `xml:"id,attr" json:"id"`
	Album         string             `xml:"album,attr,omitempty" json:"album,omitempty"`
	Title         string             `xml:"title,attr" json:"title"`
	Name          string             `xml:"name,attr" json:"name"`
	Artist        string             `xml:"artist,attr,omitempty" json:"artist,omitempty"`
	ArtistId      string             `xml:"artistId,attr,omitempty" json:"artistId,omitempty"`
	CoverArt      string             `xml:"coverArt,attr,omitempty" json:"coverArt,omitempty"`
	SongCount     int                `xml:"songCount,attr" json:"songCount"`
	Duration      int                `xml:"duration,attr" json:"duration"` // in seconds
	PlayCount     int                `xml:"playCount,attr" json:"playCount"`
	Created       string             `xml:"created,attr" json:"created"` // date_added, ISO 8601
	Starred       string             `xml:"starred,attr" json:"starred"`
	Year          int                `xml:"year,attr,omitempty" json:"year,omitempty"`
	Genre         string             `xml:"genre,attr,omitempty" json:"genre,omitempty"`
	Played        string             `xml:"played,attr,omitempty" json:"played,omitempty"` // last played, ISO 8601
	UserRating    int                `xml:"userRating,attr" json:"userRating"`
	RecordLabels  []ChildRecordLabel `xml:"recordLabels" json:"recordLabels"`
	MusicBrainzId string             `xml:"musicBrainzId,attr,omitempty" json:"musicBrainzId,omitempty"` // musicbrainz_album_id
	Genres        []ItemGenre        `xml:"genres" json:"genres"`
	DisplayArtist string             `xml:"displayArtist,attr,omitempty" json:"displayArtist,omitempty"`
	SortName      string             `xml:"sortName,attr,omitempty" json:"sortName,omitempty"`
	ReleaseDate   ItemDate           `xml:"releaseDate,omitempty" json:"releaseDate,omitempty"`
	Songs         []SubsonicChild    `xml:"song,omitempty" json:"song,omitempty"`
}

type AlbumList struct {
	Albums []AlbumId3 `xml:"album" json:"album"`
}

type AlbumList2 struct {
	Albums []AlbumId3 `xml:"album" json:"album"`
}
