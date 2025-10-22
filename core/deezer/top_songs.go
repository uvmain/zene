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
	"zene/core/types"
)

func GetTopSongs(ctx context.Context, artistName string, limit int) ([]types.TopSongRow, error) {
	artistId, err := GetDeezerArtistId(ctx, artistName)
	if err != nil {
		return []types.TopSongRow{}, err
	}

	url := fmt.Sprintf("https://api.deezer.com/artist/%d/top?limit=%d", artistId, limit)

	if err := logic.CheckContext(ctx); err != nil {
		return []types.TopSongRow{}, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []types.TopSongRow{}, fmt.Errorf("HTTP New Request failed: %v", err)
	}

	net.AddUserAgentHeaderToRequest(req)

	if err := logic.CheckContext(ctx); err != nil {
		return []types.TopSongRow{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return []types.TopSongRow{}, fmt.Errorf("HTTP error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []types.TopSongRow{}, fmt.Errorf("unexpected status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []types.TopSongRow{}, err
	}

	var data DeezerTopTracksResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return []types.TopSongRow{}, err
	}

	if len(data.Data) < 1 {
		logger.Printf("no Deezer results found for artist: %s", artistName)
		return []types.TopSongRow{}, nil
	}

	var topSongs []types.TopSongRow

	for _, track := range data.Data {
		topSongs = append(topSongs, types.TopSongRow{
			ArtistName: track.Artist.Name,
			AlbumName:  track.Album.Title,
			TrackName:  track.Title,
			SortOrder:  track.Rank,
		})
	}

	return topSongs, nil
}
