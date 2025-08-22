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

func GetSimilarArtistNames(ctx context.Context, artistName string) ([]string, error) {
	artistId, err := GetArtistId(ctx, artistName)
	if err != nil {
		return []string{}, err
	}

	url := fmt.Sprintf("https://api.deezer.com/artist/%d/related", artistId)

	if err := logic.CheckContext(ctx); err != nil {
		return []string{}, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []string{}, fmt.Errorf("HTTP New Request failed: %v", err)
	}

	net.AddUserAgentHeaderToRequest(req)

	if err := logic.CheckContext(ctx); err != nil {
		return []string{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return []string{}, fmt.Errorf("HTTP error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []string{}, fmt.Errorf("unexpected status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []string{}, err
	}

	var data DeezerArtistResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return []string{}, err
	}

	if len(data.Data) < 1 {
		return []string{}, fmt.Errorf("no Deezer results found for artist: %s", artistName)
	}

	var artistNames []string

	for _, artist := range data.Data {
		artistNames = append(artistNames, artist.Name)
	}

	return artistNames, nil
}

func GetArtistId(ctx context.Context, artistName string) (int, error) {
	if artistName == "" {
		return 0, fmt.Errorf("artist name must not be empty")
	}

	url := fmt.Sprintf("https://api.deezer.com/search/artist?q=%s", url.QueryEscape(artistName))

	if err := logic.CheckContext(ctx); err != nil {
		return 0, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("HTTP New Request failed: %v", err)
	}

	net.AddUserAgentHeaderToRequest(req)

	if err := logic.CheckContext(ctx); err != nil {
		return 0, err
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("HTTP error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var data DeezerArtistResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	if len(data.Data) < 1 {
		return 0, fmt.Errorf("no Deezer results found for artist: %s", artistName)
	}

	artistId := data.Data[0].ID

	return artistId, nil
}
