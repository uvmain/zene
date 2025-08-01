package deezer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"zene/core/logic"
	"zene/core/net"
)

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

func GetArtistArtUrlWithArtistName(ctx context.Context, artistName string) (string, error) {
	url := fmt.Sprintf("https://api.deezer.com/search/artist?q=%s", url.QueryEscape(artistName))

	if err := logic.CheckContext(ctx); err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("HTTP New Request failed: %v", err)
	}

	net.AddUserAgentHeaderToRequest(req)

	if err := logic.CheckContext(ctx); err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var data DeezerArtistResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	if len(data.Data) < 1 {
		return "", fmt.Errorf("no Deezer picture found for artist: %s", artistName)
	}

	dataArray := data.Data[0]

	imageUrl := ""
	if dataArray.PictureXl != "" {
		imageUrl = dataArray.PictureXl
	} else if dataArray.PictureBig != "" {
		imageUrl = dataArray.PictureBig
	} else if dataArray.PictureMedium != "" {
		imageUrl = dataArray.PictureMedium
	} else if dataArray.PictureSmall != "" {
		imageUrl = dataArray.PictureSmall
	} else if dataArray.Picture != "" {
		imageUrl = dataArray.Picture
	}

	if imageUrl == "" {
		return "", fmt.Errorf("no Deezer artist image found for: %s", artistName)
	}

	return imageUrl, nil
}

func GetAlbumArtUrlWithArtistNameAndAlbumName(ctx context.Context, artistName string, albumName string) (string, error) {
	if artistName == "" || albumName == "" {
		return "", fmt.Errorf("artist name and album name must not be empty")
	}
	queryParam := fmt.Sprintf("%s %s", artistName, albumName)
	url := fmt.Sprintf("https://api.deezer.com/search/album?q=%s", url.QueryEscape(queryParam))

	if err := logic.CheckContext(ctx); err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("HTTP New Request failed: %v", err)
	}

	net.AddUserAgentHeaderToRequest(req)

	if err := logic.CheckContext(ctx); err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var data DeezerAlbumResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	if len(data.Data) < 1 {
		return "", fmt.Errorf("no Deezer picture found for album: %s", albumName)
	}

	dataArray := data.Data[0]

	imageUrl := ""
	if dataArray.CoverXl != "" {
		imageUrl = dataArray.CoverXl
	} else if dataArray.CoverBig != "" {
		imageUrl = dataArray.CoverBig
	} else if dataArray.CoverMedium != "" {
		imageUrl = dataArray.CoverMedium
	} else if dataArray.CoverSmall != "" {
		imageUrl = dataArray.CoverSmall
	} else if dataArray.Cover != "" {
		imageUrl = dataArray.Cover
	}

	if imageUrl == "" {
		return "", fmt.Errorf("no Deezer album image found for: %s", artistName)
	}

	return imageUrl, nil
}
