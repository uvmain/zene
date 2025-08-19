package types

type SubsonicChild struct {
	Id                 string              `xml:"id,attr" json:"id"`
	Parent             string              `xml:"parent,attr" json:"parent"`
	IsDir              bool                `xml:"isDir,attr" json:"isDir"`
	Title              string              `xml:"title,attr" json:"title"`
	Album              string              `xml:"album,attr" json:"album"`
	Artist             string              `xml:"artist,attr" json:"artist"`
	Track              int                 `xml:"track,attr" json:"track"`
	Year               int                 `xml:"year,attr" json:"year"`
	Genre              string              `xml:"genre,attr" json:"genre"`
	CoverArt           string              `xml:"coverArt,attr" json:"coverArt"`
	Size               int                 `xml:"size,attr" json:"size"`
	ContentType        string              `xml:"contentType,attr" json:"contentType"`
	Suffix             string              `xml:"suffix,attr" json:"suffix"`
	Duration           int                 `xml:"duration,attr" json:"duration"`
	BitRate            int                 `xml:"bitRate,attr" json:"bitRate"`
	Path               string              `xml:"path,attr" json:"path"`
	DiscNumber         int                 `xml:"discNumber,attr" json:"discNumber"`
	Created            string              `xml:"created,attr" json:"created"`
	AlbumId            string              `xml:"albumId,attr" json:"albumId"`
	ArtistId           string              `xml:"artistId,attr" json:"artistId"`
	Type               string              `xml:"type,attr" json:"type"`
	IsVideo            bool                `xml:"isVideo,attr" json:"isVideo"`
	Bpm                int                 `xml:"bpm,attr" json:"bpm"`
	Comment            string              `xml:"comment,attr" json:"comment"`
	SortName           string              `xml:"sortName,attr" json:"sortName"`
	MediaType          string              `xml:"mediaType,attr" json:"mediaType"`
	MusicBrainzId      string              `xml:"musicBrainzId,attr" json:"musicBrainzId"`
	Genres             []ChildGenre        `xml:"genres>genre" json:"genres"`
	ChannelCount       int                 `xml:"channelCount,attr" json:"channelCount"`
	SamplingRate       int                 `xml:"samplingRate,attr" json:"samplingRate"`
	BitDepth           int                 `xml:"bitDepth,attr" json:"bitDepth"`
	Moods              []string            `xml:"moods>mood" json:"moods"`
	Artists            []ChildArtist       `xml:"artists>artist" json:"artists"`
	DisplayArtist      string              `xml:"displayArtist,attr" json:"displayArtist"`
	AlbumArtists       []ChildArtist       `xml:"albumArtists>artist" json:"albumArtists"`
	DisplayAlbumArtist string              `xml:"displayAlbumArtist,attr" json:"displayAlbumArtist"`
	Contributors       []ChildContributors `xml:"contributors>contributor" json:"contributors"`
	DisplayComposer    string              `xml:"displayComposer,attr" json:"displayComposer"`
	ExplicitStatus     string              `xml:"explicitStatus,attr" json:"explicitStatus"`
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

type SongsByGenre struct {
	Songs []SubsonicChild `xml:"song" json:"song"`
}
