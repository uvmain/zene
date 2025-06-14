package types

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
}

type File struct {
	FilePathAbs  string `json:"file_path_absolute"`
	DateModified string `json:"date_modified"`
}

type User struct {
	Id           int    `json:"id"` // Changed to int for consistency
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // Prevent password hash from being sent in JSON responses
	CreatedAt    string `json:"created_at"`
	IsAdmin      bool   `json:"is_admin"`
}

// HttpError represents an error with an HTTP status code and a message.
type HttpError struct {
	Code    int    `json:"-"` // The HTTP status code (e.g., http.StatusNotFound)
	Message string `json:"message"` // The error message to be sent to the client
}

// Error makes HttpError satisfy the error interface.
func (e *HttpError) Error() string {
	return e.Message
}

// NewHttpError creates a new HttpError.
func NewHttpError(code int, message string) *HttpError {
	return &HttpError{Code: code, Message: message}
}
