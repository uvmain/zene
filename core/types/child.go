package types

type SubsonicChild struct {
	Id                 string              `xml:"id,attr" json:"id"`
	Parent             string              `xml:"parent,attr,omitempty" json:"parent,omitempty"`
	IsDir              bool                `xml:"isDir,attr" json:"isDir"`
	Title              string              `xml:"title,attr" json:"title"`
	Album              string              `xml:"album,attr,omitempty" json:"album,omitempty"`
	Artist             string              `xml:"artist,attr,omitempty" json:"artist,omitempty"`
	Track              int                 `xml:"track,attr,omitempty" json:"track,omitempty"`
	Year               int                 `xml:"year,attr,omitempty" json:"year,omitempty"`
	Genre              string              `xml:"genre,attr,omitempty" json:"genre,omitempty"`
	CoverArt           string              `xml:"coverArt,attr,omitempty" json:"coverArt,omitempty"`
	Size               int                 `xml:"size,attr,omitempty" json:"size,omitempty"`
	ContentType        string              `xml:"contentType,attr,omitempty" json:"contentType,omitempty"`
	Suffix             string              `xml:"suffix,attr,omitempty" json:"suffix,omitempty"`
	Duration           int                 `xml:"duration,attr,omitempty" json:"duration,omitempty"`
	BitRate            int                 `xml:"bitRate,attr,omitempty" json:"bitRate,omitempty"`
	BitDepth           int                 `xml:"bitDepth,attr,omitempty" json:"bitDepth,omitempty"`
	SamplingRate       int                 `xml:"samplingRate,attr,omitempty" json:"samplingRate,omitempty"`
	ChannelCount       int                 `xml:"channelCount,attr,omitempty" json:"channelCount,omitempty"`
	Path               string              `xml:"path,attr,omitempty" json:"path,omitempty"`
	IsVideo            bool                `xml:"isVideo,attr,omitempty" json:"isVideo,omitempty"`
	UserRating         int                 `xml:"userRating,attr,omitempty" json:"userRating,omitempty"`
	RecordLabels       []ChildRecordLabel  `xml:"recordLabels" json:"recordLabels"`
	AverageRating      float64             `xml:"averageRating,attr,omitempty" json:"averageRating,omitempty"`
	PlayCount          int                 `xml:"playCount,attr,omitempty" json:"playCount,omitempty"`
	SongCount          int                 `xml:"songCount,attr,omitempty" json:"songCount,omitempty"`
	DiscNumber         int                 `xml:"discNumber,attr,omitempty" json:"discNumber,omitempty"`
	Created            string              `xml:"created,attr,omitempty" json:"created,omitempty"`
	Starred            string              `xml:"starred,attr,omitempty" json:"starred,omitempty"`
	AlbumId            string              `xml:"albumId,attr,omitempty" json:"albumId,omitempty"`
	ArtistId           string              `xml:"artistId,attr,omitempty" json:"artistId,omitempty"`
	Type               string              `xml:"type,attr,omitempty" json:"type,omitempty"`
	MediaType          string              `xml:"mediaType,attr,omitempty" json:"mediaType,omitempty"`
	Played             string              `xml:"played,attr,omitempty" json:"played,omitempty"`
	Bpm                int                 `xml:"bpm,attr,omitempty" json:"bpm,omitempty"`
	Comment            string              `xml:"comment,attr,omitempty" json:"comment,omitempty"`
	SortName           string              `xml:"sortName,attr,omitempty" json:"sortName,omitempty"`
	MusicBrainzId      string              `xml:"musicBrainzId,attr,omitempty" json:"musicBrainzId,omitempty"`
	Genres             []ChildGenre        `xml:"genres" json:"genres,omitempty"`
	Artists            []ChildArtist       `xml:"artists" json:"artists,omitempty"`
	DisplayArtist      string              `xml:"displayArtist,attr,omitempty" json:"displayArtist,omitempty"`
	AlbumArtists       []ChildArtist       `xml:"albumArtists" json:"albumArtists,omitempty"`
	DisplayAlbumArtist string              `xml:"displayAlbumArtist,attr,omitempty" json:"displayAlbumArtist,omitempty"`
	Contributors       []ChildContributors `xml:"contributors" json:"contributors,omitempty"`
	Moods              []string            `xml:"moods" json:"moods,omitempty"`
	DisplayComposer    string              `xml:"displayComposer,attr,omitempty" json:"displayComposer,omitempty"`
	ExplicitStatus     string              `xml:"explicitStatus,attr,omitempty" json:"explicitStatus,omitempty"`
}

type ChildGenre struct {
	Name string `xml:"name,attr" json:"name"`
}

type ChildArtist struct {
	Id   string `xml:"id,attr" json:"id"`
	Name string `xml:"name,attr" json:"name"`
}

type ChildContributors struct {
	Role   string       `xml:"role,attr" json:"role"`
	Artist *ChildArtist `xml:"artist" json:"artist"`
}

type ChildRecordLabel struct {
	Name string `xml:"name,attr" json:"name"`
}

type SongsByGenre struct {
	Songs []SubsonicChild `xml:"song" json:"song"`
}

type RandomSongs struct {
	Songs []SubsonicChild `xml:"song" json:"song"`
}
