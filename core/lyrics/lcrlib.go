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
	"strings"

	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
	"zene/core/types"
)

func GetLyricsForMusicBrainzTrackId(ctx context.Context, musicBrainzTrackId string) (types.LyricsDatabaseRow, error) {
	lyrics, err := database.GetLyricsForMusicBrainzTrackId(ctx, musicBrainzTrackId)
	if err != nil {
		logger.Printf("Error querying database in GetLyricsForMusicBrainzTrackId: %v", err)
		return types.LyricsDatabaseRow{}, fmt.Errorf("failed to query database for lyrics: %w", err)
	}
	if lyrics.PlainLyrics == "" && lyrics.SyncedLyrics == "" {
		logger.Printf("No lyrics found in database for track ID: %s", musicBrainzTrackId)
		return GetLyricsForMusicBrainzTrackIdFromLrclib(ctx, musicBrainzTrackId)
	}
	logger.Printf("Found lyrics in database for track ID: %s", musicBrainzTrackId)
	return lyrics, nil
}

func GetLyricsForMusicBrainzTrackIdFromLrclib(ctx context.Context, musicBrainzTrackId string) (types.LyricsDatabaseRow, error) {
	logger.Printf("Fetching lyrics from lrclib.net for track ID: %s", musicBrainzTrackId)

	trackMetadata, err := database.SelectTrack(ctx, musicBrainzTrackId)

	if err != nil {
		return types.LyricsDatabaseRow{}, err
	}

	trackName := url.QueryEscape(trackMetadata.Title)
	artistName := url.QueryEscape(trackMetadata.Artist)
	albumName := url.QueryEscape(trackMetadata.Album)

	durationFloat, err := strconv.ParseFloat(trackMetadata.Duration, 64)
	if err != nil {
		return types.LyricsDatabaseRow{}, fmt.Errorf("converting string to float: %v", err)
	}

	durationInt := int(math.Round(durationFloat))

	url := fmt.Sprintf("https://lrclib.net/api/get?artist_name=%s&track_name=%s&album_name=%s&duration=%d", artistName, trackName, albumName, durationInt)

	logger.Printf("Constructed URL for lyrics: %s", url)

	if err := logic.CheckContext(ctx); err != nil {
		return types.LyricsDatabaseRow{}, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return types.LyricsDatabaseRow{}, fmt.Errorf("HTTP New Request failed: %v", err)
	}

	net.AddUserAgentHeaderToRequest(req)

	if err := logic.CheckContext(ctx); err != nil {
		return types.LyricsDatabaseRow{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return types.LyricsDatabaseRow{}, fmt.Errorf("HTTP error: %w", err)
	}
	defer res.Body.Close()

	if err := logic.CheckContext(ctx); err != nil {
		return types.LyricsDatabaseRow{}, err
	}

	if res.StatusCode != http.StatusOK {
		return types.LyricsDatabaseRow{}, fmt.Errorf("unexpected status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return types.LyricsDatabaseRow{}, err
	}

	var data types.LrclibLyricsResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return types.LyricsDatabaseRow{}, err
	}

	if err := logic.CheckContext(ctx); err != nil {
		return types.LyricsDatabaseRow{}, err
	}

	var upsertData types.LyricsDatabaseRow
	upsertData.MusicBrainzTrackID = musicBrainzTrackId
	upsertData.PlainLyrics = data.PlainLyrics
	upsertData.SyncedLyrics = data.SyncedLyrics

	err = database.UpsertTrackLyrics(ctx, musicBrainzTrackId, upsertData)
	if err != nil {
		logger.Printf("Error upserting track lyrics for %s: %v", musicBrainzTrackId, err)
	}

	return upsertData, nil
}

func ParseSyncLyricTimeToMilliseconds(timeStr string) (int, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid format: %s", timeStr)
	}

	// Parse minutes
	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes: %v", err)
	}

	// parse seconds and milliseconds
	secParts := strings.Split(parts[1], ".")
	if len(secParts) != 2 {
		return 0, fmt.Errorf("invalid seconds format: %s", parts[1])
	}

	seconds, err := strconv.Atoi(secParts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid seconds: %v", err)
	}

	millis, err := strconv.Atoi(secParts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid milliseconds: %v", err)
	}
	// scale millis to be actual milliseconds (e.g., ".35" means 350 ms)
	if len(secParts[1]) == 1 {
		millis *= 100
	} else if len(secParts[1]) == 2 {
		millis *= 10
	}

	totalMillis := (minutes*60+seconds)*1000 + millis
	return totalMillis, nil
}
