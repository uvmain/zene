package deezer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
)

func GetSimilarArtistNames(ctx context.Context, artistName string) ([]string, error) {
	artistId, err := GetDeezerArtistId(ctx, artistName)
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
		logger.Printf("no Deezer results found for artist: %s", artistName)
		return []string{}, nil
	}

	var artistNames []string

	for _, artist := range data.Data {
		artistNames = append(artistNames, artist.Name)
	}

	return artistNames, nil
}
