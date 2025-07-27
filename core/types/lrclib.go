package types

type Lyrics struct {
	Id           int    `json:"id"`
	PlainLyrics  string `json:"plainLyrics"`
	SyncedLyrics string `json:"syncedLyrics"`
}

type LyricsDatabaseRow struct {
	MusicBrainzTrackID string `json:"musicBrainzTrackId"`
	PlainLyrics        string `json:"plainLyrics"`
	SyncedLyrics       string `json:"syncedLyrics"`
}

type LrclibLyricsResponse struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	TrackName    string `json:"trackName"`
	ArtistName   string `json:"artistName"`
	AlbumName    string `json:"albumName"`
	Duration     int    `json:"duration"`
	Instrumental bool   `json:"instrumental"`
	PlainLyrics  string `json:"plainLyrics"`
	SyncedLyrics string `json:"syncedLyrics"`
}
