package deezer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
)

type DeezerTopSong struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Album  string `json:"album"`
}

// in-memory cache for top songs
type topSongsCacheEntry struct {
	topSongs  []DeezerTopSong
	expiresAt time.Time
}

var (
	topSongsCache     = make(map[string]topSongsCacheEntry)
	topSongsCacheLock sync.Mutex
)

var topSongsCacheTTL = time.Hour * 24

func CleanupTopSongsCache(ctx context.Context) {
	topSongsCacheLock.Lock()
	defer topSongsCacheLock.Unlock()
	now := time.Now()
	for k, v := range topSongsCache {
		if now.After(v.expiresAt) {
			delete(topSongsCache, k)
		}
	}
}

func GetTopSongs(ctx context.Context, artistName string, limit int) ([]DeezerTopSong, error) {
	// Check cache first
	topSongsCacheLock.Lock()
	entry, found := topSongsCache[artistName]
	if found && time.Now().Before(entry.expiresAt) {
		topSongsCacheLock.Unlock()
		logger.Printf("Deezer top songs cache hit for artist: %v", artistName)
		return entry.topSongs, nil
	}
	topSongsCacheLock.Unlock()

	artistId, err := GetDeezerArtistId(ctx, artistName)
	if err != nil {
		return []DeezerTopSong{}, err
	}

	url := fmt.Sprintf("https://api.deezer.com/artist/%d/top?limit=%d", artistId, limit)

	if err := logic.CheckContext(ctx); err != nil {
		return []DeezerTopSong{}, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []DeezerTopSong{}, fmt.Errorf("HTTP New Request failed: %v", err)
	}

	net.AddUserAgentHeaderToRequest(req)

	if err := logic.CheckContext(ctx); err != nil {
		return []DeezerTopSong{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return []DeezerTopSong{}, fmt.Errorf("HTTP error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []DeezerTopSong{}, fmt.Errorf("unexpected status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []DeezerTopSong{}, err
	}

	var data DeezerTopTracksResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return []DeezerTopSong{}, err
	}

	if len(data.Data) < 1 {
		logger.Printf("no Deezer results found for artist: %s", artistName)
		return []DeezerTopSong{}, nil
	}

	var topSongs []DeezerTopSong

	for _, track := range data.Data {
		topSongs = append(topSongs, DeezerTopSong{
			Title:  track.Title,
			Artist: track.Artist.Name,
			Album:  track.Album.Title,
		})
	}

	// Store in cache
	topSongsCacheLock.Lock()
	topSongsCache[artistName] = topSongsCacheEntry{
		topSongs:  topSongs,
		expiresAt: time.Now().Add(topSongsCacheTTL),
	}
	topSongsCacheLock.Unlock()
	return topSongs, nil
}
