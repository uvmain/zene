package types

type BookmarkEntry struct {
	Id               string `json:"id" xml:"id,attr"`
	Parent           string `json:"parent" xml:"parent,attr"`
	IsDir            bool   `json:"isDir" xml:"isDir,attr"`
	Title            string `json:"title" xml:"title,attr"`
	Album            string `json:"album" xml:"album,attr"`
	Artist           string `json:"artist" xml:"artist,attr"`
	Track            int    `json:"track" xml:"track,attr"`
	Year             int    `json:"year" xml:"year,attr"`
	Genre            string `json:"genre" xml:"genre,attr"`
	CoverArt         string `json:"coverArt" xml:"coverArt,attr"`
	Size             int    `json:"size" xml:"size,attr"`
	ContentType      string `json:"contentType" xml:"contentType,attr"`
	Suffix           string `json:"suffix" xml:"suffix,attr"`
	Duration         int    `json:"duration" xml:"duration,attr"`
	BitRate          int    `json:"bitRate" xml:"bitRate,attr"`
	BitDepth         int    `json:"bitDepth" xml:"bitDepth,attr"`
	SamplingRate     int    `json:"samplingRate" xml:"samplingRate,attr"`
	ChannelCount     int    `json:"channelCount" xml:"channelCount,attr"`
	Path             string `json:"path" xml:"path,attr"`
	Created          string `json:"created" xml:"created,attr"`
	AlbumID          string `json:"albumId" xml:"albumId,attr"`
	ArtistID         string `json:"artistId" xml:"artistId,attr"`
	Type             string `json:"type" xml:"type,attr"`
	IsVideo          bool   `json:"isVideo" xml:"isVideo,attr"`
	BookmarkPosition int    `json:"bookmarkPosition" xml:"bookmarkPosition,attr"`
}

type Bookmark struct {
	Entry    BookmarkEntry `json:"entry" xml:"entry"`
	Position int           `json:"position" xml:"position,attr"`
	Username string        `json:"username" xml:"username,attr"`
	Comment  string        `json:"comment" xml:"comment,attr"`
	Created  string        `json:"created" xml:"created,attr"`
	Changed  string        `json:"changed" xml:"changed,attr"`
}

type Bookmarks struct {
	Bookmarks []Bookmark `json:"bookmark" xml:"bookmark"`
}
