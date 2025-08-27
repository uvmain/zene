package types

type PlaylistRow struct {
	Id          string          `json:"id" xml:"id,attr"`
	Name        string          `json:"name" xml:"name,attr"`
	Comment     string          `json:"comment,omitempty" xml:"comment,attr,omitempty"`
	Owner       string          `json:"owner" xml:"owner,attr"`
	Public      bool            `json:"public" xml:"public,attr"`
	SongCount   int             `json:"songCount" xml:"song_count,attr"`
	Duration    int             `json:"duration" xml:"duration,attr"`
	Created     string          `json:"created" xml:"created,attr"`
	Changed     string          `json:"changed" xml:"changed,attr"`
	CoverArt    string          `json:"coverArt" xml:"cover_art,attr"`
	AllowedUser []string        `json:"allowedUser" xml:"allowed_user,attr"`
	Entries     []SubsonicChild `json:"entry,omitempty" xml:"entry,omitempty"`
}

type Playlist struct {
	Playlist PlaylistRow `json:"playlist" xml:"playlist"`
}

type Playlists struct {
	Playlist []PlaylistRow `json:"playlist" xml:"playlist"`
}
