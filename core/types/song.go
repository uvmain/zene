package types

type SubsonicSong struct {
	Id                 string             `xml:"id,attr" json:"id"`
	Parent             string             `xml:"parent,attr" json:"parent"`
	IsDir              bool               `xml:"isDir,attr" json:"isDir"`
	Title              string             `xml:"title,attr" json:"title"`
	Album              string             `xml:"album,attr" json:"album"`
	Artist             string             `xml:"artist,attr" json:"artist"`
	Track              int64              `xml:"track,attr" json:"track"`
	Year               int64              `xml:"year,attr" json:"year"`
	Genre              string             `xml:"genre,attr" json:"genre"`
	CoverArt           string             `xml:"coverArt,attr" json:"coverArt"`
	Size               int64              `xml:"size,attr" json:"size"`
	ContentType        string             `xml:"contentType,attr" json:"contentType"`
	Suffix             string             `xml:"suffix,attr" json:"suffix"`
	Duration           int64              `xml:"duration,attr" json:"duration"`
	BitRate            int64              `xml:"bitRate,attr" json:"bitRate"`
	Path               string             `xml:"path,attr" json:"path"`
	DiscNumber         int64              `xml:"discNumber,attr" json:"discNumber"`
	Created            string             `xml:"created,attr" json:"created"`
	AlbumId            string             `xml:"albumId,attr" json:"albumId"`
	ArtistId           string             `xml:"artistId,attr" json:"artistId"`
	Type               string             `xml:"type,attr" json:"type"`
	IsVideo            bool               `xml:"isVideo,attr" json:"isVideo"`
	Bpm                int64              `xml:"bpm,attr" json:"bpm"`
	Comment            string             `xml:"comment,attr" json:"comment"`
	SortName           string             `xml:"sortName,attr" json:"sortName"`
	MediaType          string             `xml:"mediaType,attr" json:"mediaType"`
	MusicBrainzId      string             `xml:"musicBrainzId,attr" json:"musicBrainzId"`
	Genres             []SongGenre        `xml:"genres>genre" json:"genres"`
	ChannelCount       int64              `xml:"channelCount,attr" json:"channelCount"`
	SamplingRate       int64              `xml:"samplingRate,attr" json:"samplingRate"`
	BitDepth           int64              `xml:"bitDepth,attr" json:"bitDepth"`
	Moods              []string           `xml:"moods>mood" json:"moods"`
	Artists            []SongArtist       `xml:"artists>artist" json:"artists"`
	DisplayArtist      string             `xml:"displayArtist,attr" json:"displayArtist"`
	AlbumArtists       []SongArtist       `xml:"albumArtists>artist" json:"albumArtists"`
	DisplayAlbumArtist string             `xml:"displayAlbumArtist,attr" json:"displayAlbumArtist"`
	Contributors       []SongContributors `xml:"contributors>contributor" json:"contributors"`
	DisplayComposer    string             `xml:"displayComposer,attr" json:"displayComposer"`
	ExplicitStatus     string             `xml:"explicitStatus,attr" json:"explicitStatus"`
}

type SongGenre struct {
	Name string `xml:"name,attr" json:"name"`
}

type SongArtist struct {
	Id   string `xml:"id,attr" json:"id"`
	Name string `xml:"name,attr" json:"name"`
}

type SongContributors struct {
	Role   string      `xml:"role,attr" json:"role"`
	Artist *SongArtist `xml:"artist" json:"artist"`
}
