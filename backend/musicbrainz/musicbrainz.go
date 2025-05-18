package musicbrainz

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"zene/net"
	"zene/types"
)

func GetAlbumArtUrl(musicBrainzAlbumId string) (string, error) {
	url := fmt.Sprintf("https://coverartarchive.org/release/%s", musicBrainzAlbumId)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("HTTP New Request failed: %v", err)
	}

	net.AddUserAgentHeaderToRequest(req)

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

	var data types.CoverArtResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	return data.Images[0].Image, nil
}
