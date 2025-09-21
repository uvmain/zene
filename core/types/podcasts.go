package types

type PodcastStatus string

const (
	PodcastStatusNew         PodcastStatus = "new"
	PodcastStatusDownloading PodcastStatus = "downloading"
	PodcastStatusCompleted   PodcastStatus = "completed"
	PodcastStatusError       PodcastStatus = "error"
	PodcastStatusDeleted     PodcastStatus = "deleted"
	PodcastStatusSkipped     PodcastStatus = "skipped"
)

type PodcastChannel struct {
	Id               string           `json:"id" xml:"id,attr"`
	ParentId         string           `json:"parent" xml:"parent,attr"`
	IsDir            string           `json:"isDir" xml:"isDir,attr" default:"false"`
	Title            string           `json:"title" xml:"title,attr"`
	Url              string           `json:"url" xml:"url,attr"`
	Description      string           `json:"description" xml:"description,attr"`
	CoverArt         string           `json:"coverArt" xml:"coverArt,attr"`
	OriginalImageUrl string           `json:"originalImageUrl" xml:"originalImageUrl,attr"`
	Status           PodcastStatus    `json:"status" xml:"status,attr"`
	Type             string           `json:"type" xml:"type,attr" default:"podcast"`
	IsVideo          string           `json:"isVideo" xml:"isVideo,attr" default:"false"`
	StreamId         string           `json:"streamId" xml:"streamId,attr"`
	ChannelId        string           `json:"channelId" xml:"channelId,attr"`
	LastRefresh      string           `json:"lastRefresh" xml:"lastRefresh,attr"`
	CreatedAt        string           `json:"created" xml:"created,attr"`
	ErrorMessage     string           `json:"errorMessage,omitempty" xml:"errorMessage,attr,omitempty"`
	Episodes         []PodcastEpisode `json:"episode" xml:"episode"`
}

type PodcastEpisode struct {
	Id          string        `json:"id" xml:"id,attr"`
	StreamId    string        `json:"streamId" xml:"streamId,attr"`
	ChannelId   string        `json:"channelId" xml:"channelId,attr"`
	Title       string        `json:"title" xml:"title,attr"`
	Description string        `json:"description" xml:"description,attr"`
	PublishDate string        `json:"publishDate" xml:"publishDate,attr"`
	Status      PodcastStatus `json:"status" xml:"status,attr"`
	Parent      string        `json:"parent" xml:"parent,attr"`
	IsDir       string        `json:"isDir" xml:"isDir,attr"`
	Year        string        `json:"year" xml:"year,attr"`
	Genre       string        `json:"genre" xml:"genre,attr"`
	Genres      []ChildGenre  `xml:"genres" json:"genres"`
	CoverArt    string        `json:"coverArt" xml:"coverArt,attr"`
	Size        string        `json:"size" xml:"size,attr"`
	ContentType string        `json:"contentType" xml:"contentType,attr"`
	Suffix      string        `json:"suffix" xml:"suffix,attr"`
	Duration    string        `json:"duration" xml:"duration,attr"`
	BitRate     string        `json:"bitRate" xml:"bitRate,attr"`
	Path        string        `json:"path" xml:"path,attr"`
	SourceUrl   string        `json:"sourceUrl" xml:"sourceUrl,attr"`
}

type PodcastEpisodeRow struct {
	ChannelId   string
	Guid        string
	Title       string
	Album       string
	Artist      string
	Year        string
	CoverArt    string
	Size        string
	ContentType string
	Suffix      string
	Duration    int
	BitRate     int
	Description string
	PublishDate string
	Status      string
	FilePath    string
	CreatedAt   string
	SourceUrl   string
}

type PodcastChannels struct {
	PodcastChannels []PodcastChannel `json:"channel" xml:"channel"`
}

type NewestPodcasts struct {
	Episodes []PodcastEpisode `json:"episode,omitempty" xml:"episode,omitempty"`
}
