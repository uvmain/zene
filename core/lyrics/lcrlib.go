package lyrics

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
	"zene/core/types"
)

func GetLyricsForMusicBrainzTrackId(ctx context.Context, musicBrainzTrackId string) (types.Lyrics, error) {
	lyrics, err := database.GetLyricsForMusicBrainzTrackId(ctx, musicBrainzTrackId)
	if err != nil {
		logger.Printf("Error querying database in GetLyricsForMusicBrainzTrackId: %v", err)
		return types.Lyrics{}, fmt.Errorf("failed to query database for lyrics: %w", err)
	}
	if lyrics.PlainLyrics == "" && lyrics.SyncedLyrics == "" {
		logger.Printf("No lyrics found in database for track ID: %s", musicBrainzTrackId)
		return GetLyricsForMusicBrainzTrackIdFromLrclib(ctx, musicBrainzTrackId)
	}
	logger.Printf("Found lyrics in database for track ID: %s", musicBrainzTrackId)
	return lyrics, nil
}

func GetLyricsForMusicBrainzTrackIdFromLrclib(ctx context.Context, musicBrainzTrackId string) (types.Lyrics, error) {
	logger.Printf("Fetching lyrics from lrclib.net for track ID: %s", musicBrainzTrackId)

	trackMetadata, err := database.SelectTrack(ctx, musicBrainzTrackId)

	if err != nil {
		return types.Lyrics{}, err
	}

	trackName := url.QueryEscape(trackMetadata.Title)
	artistName := url.QueryEscape(trackMetadata.Artist)
	albumName := url.QueryEscape(trackMetadata.Album)

	durationFloat, err := strconv.ParseFloat(trackMetadata.Duration, 64)
	if err != nil {
		fmt.Println("Error converting string to float:", err)
		return types.Lyrics{}, err
	}

	durationInt := int(math.Round(durationFloat))

	url := fmt.Sprintf("https://lrclib.net/api/get?artist_name=%s&track_name=%s&album_name=%s&duration=%d", artistName, trackName, albumName, durationInt)

	logger.Printf("Constructed URL for lyrics: %s", url)

	if err := logic.CheckContext(ctx); err != nil {
		return types.Lyrics{}, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return types.Lyrics{}, fmt.Errorf("HTTP New Request failed: %v", err)
	}

	net.AddUserAgentHeaderToRequest(req)

	if err := logic.CheckContext(ctx); err != nil {
		return types.Lyrics{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return types.Lyrics{}, fmt.Errorf("HTTP error: %w", err)
	}
	defer res.Body.Close()

	if err := logic.CheckContext(ctx); err != nil {
		return types.Lyrics{}, err
	}

	if res.StatusCode != http.StatusOK {
		return types.Lyrics{}, fmt.Errorf("unexpected status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return types.Lyrics{}, err
	}

	var data types.Lyrics
	if err := json.Unmarshal(body, &data); err != nil {
		return types.Lyrics{}, err
	}

	if err := logic.CheckContext(ctx); err != nil {
		return types.Lyrics{}, err
	}

	err = database.UpsertTrackLyrics(ctx, musicBrainzTrackId, data)
	if err != nil {
		logger.Printf("Error upserting track lyrics for %s: %v", musicBrainzTrackId, err)
	}

	return data, nil
}
