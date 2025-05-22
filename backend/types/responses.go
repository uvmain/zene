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
	Genres              string `json:"genres"`
	ReleaseDate         string `json:"release_date"`
}

type ScanResponse struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
}

type GenreResponse struct {
	Genre string `json:"genre"`
	Count int    `json:"count"`
}

type LastFmAlbumInfoResponse struct {
	Album struct {
		Image []struct {
			Size string `json:"size"`
			URL  string `json:"#text"`
		} `json:"image"`
	} `json:"album"`
}
