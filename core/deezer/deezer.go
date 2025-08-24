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
