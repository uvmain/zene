package types

import "time"

type Tags struct {
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

type Metadata struct {
	FilePath            string `json:"file_path"`
	FileName            string `json:"file_name"`
	DateAdded           string `json:"date_added"`
	DateModified        string `json:"date_modified"`
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

type MetadataWithPlaycounts struct {
	FilePath            string `json:"file_path"`
	FileName            string `json:"file_name"`
	DateAdded           string `json:"date_added"`
	DateModified        string `json:"date_modified"`
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
	UserPlayCount       int64  `json:"user_play_count"`
	GlobalPlayCount     int64  `json:"global_play_count"`
}

type AlbumArtRow struct {
	MusicbrainzAlbumId string `json:"musicbrainz_album_id"`
	DateAdded          string `json:"date_added"`
	DateModified       string `json:"date_modified"`
}

type ArtistArtRow struct {
	MusicbrainzArtistId string `json:"musicbrainz_artist_id"`
	DateAdded           string `json:"date_added"`
	DateModified        string `json:"date_modified"`
}

type SessionCheck struct {
	LoggedIn bool `json:"loggedIn"`
	IsAdmin  bool `json:"isAdmin"`
}

type File struct {
	FileName     string `json:"file_name"`
	FilePathAbs  string `json:"file_path_absolute"`
	DateModified string `json:"date_modified"`
}

type Playcount struct {
	Id                 int64  `json:"id"`
	UserId             int64  `json:"user_id"`
	MusicBrainzTrackID string `json:"musicbrainz_track_id"`
	PlayCount          int64  `json:"play_count"`
	LastPlayed         string `json:"last_played"`
}

type AudioCacheEntry struct {
	CacheKey     string    `json:"cache_key"`
	LastAccessed time.Time `json:"last_accessed"`
}
