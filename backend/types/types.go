package types

type ScanRow = struct {
	Id           int    `json:"id"`
	ScanDate     string `json:"scan_date"`
	FileCount    string `json:"file_count"`
	DateModified string `json:"date_modified"`
}

type FilesRow = struct {
	Id           int    `json:"id"`
	DirPath      string `json:"directory"`
	Filename     string `json:"filename"`
	DateAdded    string `json:"date_added"`
	DateModified string `json:"date_modified"`
}

type TrackMetadata struct {
	Id                  int    `json:"id"`
	FileId              int    `json:"file_id"`
	Filename            string `json:"filename"`
	Format              string `json:"format"`
	Duration            string `json:"duration"`
	Size                string `json:"size"`
	Bitrate             string `json:"bitrate"`
	Title               string `json:"title"`
	Artist              string `json:"artist"`
	Album               string `json:"album"`
	AlbumArtist         string `json:"album_artist"`
	Genre               string `json:"genre"`
	TrackNumber         string `json:"track_number"`
	TotalTracks         string `json:"total_tracks"`
	DiscNumber          string `json:"disc_number"`
	TotalDiscs          string `json:"total_discs"`
	ReleaseDate         string `json:"release_date"`
	MusicBrainzArtistID string `json:"musicbrainz_artist_id"`
	MusicBrainzAlbumID  string `json:"musicbrainz_album_id"`
	MusicBrainzTrackID  string `json:"musicbrainz_track_id"`
	Label               string `json:"label"`
}

type AlbumArtRow struct {
	MusicbrainzAlbumId string `json:"musicbrainz_album_id"`
	DateAdded          string `json:"date_added"`
	DateModified       string `json:"date_modified"`
}
