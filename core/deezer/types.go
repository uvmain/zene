package deezer

type DeezerArtistResponse struct {
	Data []struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Link          string `json:"link"`
		Picture       string `json:"picture"`
		PictureSmall  string `json:"picture_small"`
		PictureMedium string `json:"picture_medium"`
		PictureBig    string `json:"picture_big"`
		PictureXl     string `json:"picture_xl"`
		NbAlbum       int    `json:"nb_album"`
		NbFan         int    `json:"nb_fan"`
		Radio         bool   `json:"radio"`
		Tracklist     string `json:"tracklist"`
		Type          string `json:"type"`
	} `json:"data"`
	Total int `json:"total"`
}

type DeezerAlbumResponse struct {
	Data []struct {
		ID             int    `json:"id"`
		Title          string `json:"title"`
		Link           string `json:"link"`
		Cover          string `json:"cover"`
		CoverSmall     string `json:"cover_small"`
		CoverMedium    string `json:"cover_medium"`
		CoverBig       string `json:"cover_big"`
		CoverXl        string `json:"cover_xl"`
		Md5Image       string `json:"md5_image"`
		GenreID        int    `json:"genre_id"`
		NbTracks       int    `json:"nb_tracks"`
		RecordType     string `json:"record_type"`
		Tracklist      string `json:"tracklist"`
		ExplicitLyrics bool   `json:"explicit_lyrics"`
		Artist         struct {
			ID            int    `json:"id"`
			Name          string `json:"name"`
			Link          string `json:"link"`
			Picture       string `json:"picture"`
			PictureSmall  string `json:"picture_small"`
			PictureMedium string `json:"picture_medium"`
			PictureBig    string `json:"picture_big"`
			PictureXl     string `json:"picture_xl"`
			Tracklist     string `json:"tracklist"`
			Type          string `json:"type"`
		} `json:"artist"`
		Type string `json:"type"`
	} `json:"data"`
	Total int `json:"total"`
}

type DeezerTopTracksResponse struct {
	Data []struct {
		ID                    int    `json:"id"`
		Readable              bool   `json:"readable"`
		Title                 string `json:"title"`
		TitleShort            string `json:"title_short"`
		TitleVersion          string `json:"title_version"`
		Link                  string `json:"link"`
		Duration              int    `json:"duration"`
		Rank                  int    `json:"rank"`
		ExplicitLyrics        bool   `json:"explicit_lyrics"`
		ExplicitContentLyrics int    `json:"explicit_content_lyrics"`
		ExplicitContentCover  int    `json:"explicit_content_cover"`
		Preview               string `json:"preview"`
		Contributors          []struct {
			ID            int    `json:"id"`
			Name          string `json:"name"`
			Link          string `json:"link"`
			Share         string `json:"share"`
			Picture       string `json:"picture"`
			PictureSmall  string `json:"picture_small"`
			PictureMedium string `json:"picture_medium"`
			PictureBig    string `json:"picture_big"`
			PictureXl     string `json:"picture_xl"`
			Radio         bool   `json:"radio"`
			Tracklist     string `json:"tracklist"`
			Type          string `json:"type"`
			Role          string `json:"role"`
		} `json:"contributors"`
		Md5Image string `json:"md5_image"`
		Artist   struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Tracklist string `json:"tracklist"`
			Type      string `json:"type"`
		} `json:"artist"`
		Album struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			Cover       string `json:"cover"`
			CoverSmall  string `json:"cover_small"`
			CoverMedium string `json:"cover_medium"`
			CoverBig    string `json:"cover_big"`
			CoverXl     string `json:"cover_xl"`
			Md5Image    string `json:"md5_image"`
			Tracklist   string `json:"tracklist"`
			Type        string `json:"type"`
		} `json:"album"`
		Type string `json:"type"`
	} `json:"data"`
	Total int    `json:"total"`
	Next  string `json:"next"`
}
