package deezer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
)

// in-memory cache for similar artist names
type similarArtistCacheEntry struct {
	names     []string
	expiresAt time.Time
}

var (
	similarArtistCache     = make(map[string]similarArtistCacheEntry)
	similarArtistCacheLock sync.Mutex
)

var similarArtistCacheTTL = time.Hour * 24

func CleanupSimilarArtistsCache(ctx context.Context) {
	similarArtistCacheLock.Lock()
	defer similarArtistCacheLock.Unlock()
	now := time.Now()
	for k, v := range similarArtistCache {
		if now.After(v.expiresAt) {
			delete(similarArtistCache, k)
		}
	}
}

func GetSimilarArtistNames(ctx context.Context, artistName string) ([]string, error) {
	// Check cache first
	similarArtistCacheLock.Lock()
	entry, found := similarArtistCache[artistName]
	if found && time.Now().Before(entry.expiresAt) {
		similarArtistCacheLock.Unlock()
		logger.Printf("Deezer similar artists cache hit for artist: %s", artistName)
		return entry.names, nil
	}
	similarArtistCacheLock.Unlock()

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

	// Store in cache
	similarArtistCacheLock.Lock()
	similarArtistCache[artistName] = similarArtistCacheEntry{
		names:     artistNames,
		expiresAt: time.Now().Add(similarArtistCacheTTL),
	}
	similarArtistCacheLock.Unlock()
	return artistNames, nil
}

func GetDeezerArtistId(ctx context.Context, artistName string) (int, error) {
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
