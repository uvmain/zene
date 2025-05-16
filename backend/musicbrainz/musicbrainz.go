package musicbrainz

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"zene/types"
)

func GetAlbumArtUrl(musicBrainzAlbumId string) (string, error) {
	url := fmt.Sprintf("https://coverartarchive.org/release/%s", musicBrainzAlbumId)

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

	var data types.CoverArtResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	return data.Images[0].Image, nil
}
