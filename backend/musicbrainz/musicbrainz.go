package musicbrainz

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"zene/net"
	"zene/types"
)

func GetMetadataForMusicBrainzAlbumId(musicBrainzAlbumId string) (types.MbRelease, error) {
	log.Printf("Fetching metadata from MB for album ID: %s", musicBrainzAlbumId)
	url := fmt.Sprintf("http://musicbrainz.org/ws/2/release/%s?fmt=json", musicBrainzAlbumId)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return types.MbRelease{}, fmt.Errorf("HTTP New Request failed: %v", err)
	}

	net.AddUserAgentHeaderToRequest(req)

	res, err := client.Do(req)
	if err != nil {
		return types.MbRelease{}, fmt.Errorf("HTTP error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return types.MbRelease{}, fmt.Errorf("unexpected status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return types.MbRelease{}, err
	}

	var data types.MbRelease
	if err := json.Unmarshal(body, &data); err != nil {
		return types.MbRelease{}, err
	}

	return data, nil
}

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

	var data types.MbCoverArtResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	return data.Images[0].Image, nil
}

func GetArtistArtUrl(musicBrainzArtistId string) (string, error) {
	url := fmt.Sprintf("http://musicbrainz.org/ws/2/artist/%s?inc=url-rels&fmt=json", musicBrainzArtistId)

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

	var data types.MbArtist
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	var entityId string

	for _, relation := range data.Relations {
		if strings.Contains(relation.URL.Resource, "wikidata") {
			var splitStrings = strings.Split(relation.URL.Resource, "/")
			entityId = splitStrings[len(splitStrings)-1]
			break
		}
	}

	if entityId == "" {
		return "", fmt.Errorf("no wikidata entity found for artist %s", musicBrainzArtistId)
	}

	url = fmt.Sprintf("https://www.wikidata.org/wiki/Special:EntityData/%s.json", entityId)

	client = &http.Client{}
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("HTTP New Request failed: %v", err)
	}

	net.AddUserAgentHeaderToRequest(req)

	res, err = client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", res.Status)
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var wikidataResponse map[string]any
	if err := json.Unmarshal(body, &wikidataResponse); err != nil {
		return "", err
	}
	entities, ok := wikidataResponse["entities"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid entities format in wikidata response")
	}

	entity, ok := entities[entityId].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid entity format for ID %s in wikidata response", entityId)
	}

	claims, ok := entity["claims"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid claims format in wikidata response")
	}

	p18, ok := claims["P18"].([]any)
	if !ok || len(p18) == 0 {
		return "", fmt.Errorf("no P18 claim found in wikidata response")
	}

	mainSnak, ok := p18[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid P18 claim format in wikidata response")
	}

	dataValue, ok := mainSnak["mainsnak"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid mainsnak format in wikidata response")
	}

	valueField, ok := dataValue["datavalue"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid datavalue format in wikidata response")
	}

	value, ok := valueField["value"].(string)
	if !ok {
		return "", fmt.Errorf("invalid value format in wikidata response")
	}

	url = fmt.Sprintf("https://commons.wikimedia.org/wiki/Special:FilePath/%s", value)

	return url, nil
}
