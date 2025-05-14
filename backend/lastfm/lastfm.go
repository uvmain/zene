package lastfm

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"zene/config"
	"zene/downloader"
	"zene/types"
)

type AlbumDownloadLinks = struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
	XL     string `json:"extra_large"`
}

func GetAlbumArt(musicBrainzAlbumId string) (string, error) {
	if config.LastFmApiKey == "" {
		return "", fmt.Errorf("error fetching album art for %s: LASTFM_API_KEY environment variable is not defined", musicBrainzAlbumId)
	}

	url := fmt.Sprintf(
		"https://ws.audioscrobbler.com/2.0/?method=album.getinfo&api_key=%s&mbid=%s&format=json",
		config.LastFmApiKey,
		musicBrainzAlbumId,
	)

	res, err := http.Get(url)
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

	var data types.LastFmAlbumInfoResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	albumDownloadLinks := AlbumDownloadLinks{}

	for _, image := range data.Album.Image {
		switch image.Size {
		case "small":
			albumDownloadLinks.Small = image.URL
		case "medium":
			albumDownloadLinks.Medium = image.URL
		case "large":
			albumDownloadLinks.Large = image.URL
		case "extralarge":
			albumDownloadLinks.XL = image.URL
		}
	}

	go downloader.DownloadAndSaveAsJPG(albumDownloadLinks.Small, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "sm"}, "_")))
	go downloader.DownloadAndSaveAsJPG(albumDownloadLinks.Medium, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "md"}, "_")))
	go downloader.DownloadAndSaveAsJPG(albumDownloadLinks.Large, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "lg"}, "_")))
	go downloader.DownloadAndSaveAsJPG(albumDownloadLinks.XL, filepath.Join(config.AlbumArtFolder, strings.Join([]string{musicBrainzAlbumId, "xl"}, "_")))

	log.Printf("%v", albumDownloadLinks.Large)

	return "", nil
}
