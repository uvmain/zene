package types

import "time"

type ContextKey string

type FileMetadata struct {
	Format              string `json:"format"`
	Duration            string `json:"duration"`
	Size                string `json:"size"`
	Bitrate             string `json:"bitrate"`
	Title               string `json:"title"`
	Artist              string `json:"artist"`
	Album               string `json:"album"`
	AlbumArtist         string `json:"album_artist"`
	Genre               string `json:"genre"`
	TrackNumber         int    `json:"track_number"`
	TotalTracks         int    `json:"total_tracks"`
	DiscNumber          int    `json:"disc_number"`
	TotalDiscs          int    `json:"total_discs"`
	ReleaseDate         string `json:"release_date"`
	MusicBrainzArtistID string `json:"musicbrainz_artist_id"`
	MusicBrainzAlbumID  string `json:"musicbrainz_album_id"`
	MusicBrainzTrackID  string `json:"musicbrainz_track_id"`
	Label               string `json:"label"`
	Codec               string `json:"codec_name"`
	BitDepth            int    `json:"bits_per_raw_sample"`
	SampleRate          int    `json:"sample_rate"`
	Channels            int    `json:"channels"`
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
	TrackNumber         int    `json:"track_number"`
	TotalTracks         int    `json:"total_tracks"`
	DiscNumber          int    `json:"disc_number"`
	TotalDiscs          int    `json:"total_discs"`
	ReleaseDate         string `json:"release_date"`
	MusicBrainzArtistID string `json:"musicbrainz_artist_id"`
	MusicBrainzAlbumID  string `json:"musicbrainz_album_id"`
	MusicBrainzTrackID  string `json:"musicbrainz_track_id"`
	Label               string `json:"label"`
	MusicFolderId       int    `json:"music_folder_id"`
	Codec               string `json:"codec_name"`
	BitDepth            int    `json:"bits_per_raw_sample"`
	SampleRate          int    `json:"sample_rate"`
	Channels            int    `json:"channels"`
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
	TrackNumber         int    `json:"track_number"`
	TotalTracks         int    `json:"total_tracks"`
	DiscNumber          int    `json:"disc_number"`
	TotalDiscs          int    `json:"total_discs"`
	ReleaseDate         string `json:"release_date"`
	MusicBrainzArtistID string `json:"musicbrainz_artist_id"`
	MusicBrainzAlbumID  string `json:"musicbrainz_album_id"`
	MusicBrainzTrackID  string `json:"musicbrainz_track_id"`
	Label               string `json:"label"`
	MusicFolderId       int    `json:"music_folder_id"`
	Codec               string `json:"codec_name"`
	BitDepth            int    `json:"bits_per_raw_sample"`
	SampleRate          int    `json:"sample_rate"`
	Channels            int    `json:"channels"`
	UserPlayCount       int    `json:"user_play_count"`
	GlobalPlayCount     int    `json:"global_play_count"`
}

type AlbumArtRow struct {
	DateModified string `json:"date_modified"`
	FilePath     string `json:"file_path"`
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
	FilePath     string `json:"file_path"`
	FilePathAbs  string `json:"file_path_absolute"`
	DateModified string `json:"date_modified"`
}

type Playcount struct {
	Id                 int    `json:"id"`
	UserId             int    `json:"user_id"`
	MusicBrainzTrackID string `json:"musicbrainz_track_id"`
	PlayCount          int    `json:"play_count"`
	LastPlayed         string `json:"last_played"`
}

type AudioCacheEntry struct {
	CacheKey     string    `json:"cache_key"`
	LastAccessed time.Time `json:"last_accessed"`
}

type TopSongRow struct {
	ArtistName string
	AlbumName  string
	TrackName  string
	SortOrder  int
}
