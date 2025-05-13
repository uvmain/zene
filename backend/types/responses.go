package types

type ArtistResponse struct {
	MusicBrainzArtistID string `json:"musicbrainz_artist_id"`
	Artist              string `json:"artist"`
}

type AlbumsResponse struct {
	MusicBrainzAlbumID  string `json:"musicbrainz_album_id"`
	Album               string `json:"album"`
	MusicBrainzArtistID string `json:"musicbrainz_artist_id"`
	Artist              string `json:"artist"`
}
