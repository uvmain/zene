package types

type ArtistResponse struct {
	MusicBrainzArtistID string `json:"musicbrainz_artist_id"`
	Artist              string `json:"artist"`
	ImageURL            string `json:"image_url"`
}

type AlbumsResponse struct {
	MusicBrainzAlbumID  string `json:"musicbrainz_album_id"`
	Album               string `json:"album"`
	MusicBrainzArtistID string `json:"musicbrainz_artist_id"`
	Artist              string `json:"artist"`
	Genres              string `json:"genres"`
	ReleaseDate         string `json:"release_date"`
}
