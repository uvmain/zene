package database

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"zene/core/logic"
	"zene/core/types"
)

func SelectDistinctGenres(ctx context.Context, limitParam string, searchParam string) ([]types.GenreResponse, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.GenreResponse{}, fmt.Errorf("taking a db conn from the pool in SelectDistinctGenres: %v", err)
	}
	defer DbPool.Put(conn)

	stmtText := "SELECT DISTINCT genre FROM metadata"

	if limitParam != "" {
		limitInt, err := strconv.Atoi(limitParam)
		if err != nil {
			return []types.GenreResponse{}, fmt.Errorf("invalid limit value: %v", err)
		}
		stmtText = fmt.Sprintf("%s limit %d", stmtText, limitInt)
	}

	stmtText = fmt.Sprintf("%s;", stmtText)

	stmt := conn.Prep(stmtText)
	defer stmt.Finalize()

	var genres []string

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return []types.GenreResponse{}, err
		} else if !hasRow {
			break
		} else {
			row := stmt.GetText("genre")
			splits := strings.Split(row, ";")
			for _, split := range splits {
				trimmed := strings.TrimSpace(split)
				if trimmed != "" {
					if searchParam != "" {
						if strings.Contains(strings.ToLower(trimmed), strings.ToLower(searchParam)) {
							genres = append(genres, trimmed)
						}
					} else {
						genres = append(genres, trimmed)
					}
				}
			}
		}
	}

	dict := map[string]int{}
	for _, num := range genres {
		dict[num]++
	}

	var ss []types.GenreResponse
	for k, v := range dict {
		ss = append(ss, types.GenreResponse{
			Genre: k,
			Count: v,
		})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Count > ss[j].Count
	})

	return ss, nil
}

func SelectTracksByGenres(ctx context.Context, genres []string, andOr string, limit int64, random string) ([]types.MetadataWithPlaycounts, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.MetadataWithPlaycounts{}, fmt.Errorf("taking a db conn from the pool in SelectTracksByAlbumId: %v", err)
	}
	defer DbPool.Put(conn)

	userId, _ := logic.GetUserIdFromContext(ctx)
	stmtText := getMetadataWithGenresSql(userId, genres, andOr, limit, random)

	stmt := conn.Prep(stmtText)
	defer stmt.Finalize()

	var rows []types.MetadataWithPlaycounts

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return []types.MetadataWithPlaycounts{}, err
		} else if !hasRow {
			break
		} else {
			row := types.MetadataWithPlaycounts{
				FilePath:            stmt.GetText("file_path"),
				DateAdded:           stmt.GetText("date_added"),
				DateModified:        stmt.GetText("date_modified"),
				FileName:            stmt.GetText("file_name"),
				Format:              stmt.GetText("format"),
				Duration:            stmt.GetText("duration"),
				Size:                stmt.GetText("size"),
				Bitrate:             stmt.GetText("bitrate"),
				Title:               stmt.GetText("title"),
				Artist:              stmt.GetText("artist"),
				Album:               stmt.GetText("album"),
				AlbumArtist:         stmt.GetText("album_artist"),
				Genre:               stmt.GetText("genre"),
				TrackNumber:         stmt.GetText("track_number"),
				TotalTracks:         stmt.GetText("total_tracks"),
				DiscNumber:          stmt.GetText("disc_number"),
				TotalDiscs:          stmt.GetText("total_discs"),
				ReleaseDate:         stmt.GetText("release_date"),
				MusicBrainzArtistID: stmt.GetText("musicbrainz_artist_id"),
				MusicBrainzAlbumID:  stmt.GetText("musicbrainz_album_id"),
				MusicBrainzTrackID:  stmt.GetText("musicbrainz_track_id"),
				Label:               stmt.GetText("label"),
				UserPlayCount:       stmt.GetInt64("user_play_count"),
				GlobalPlayCount:     stmt.GetInt64("global_play_count"),
			}
			rows = append(rows, row)
		}
	}
	if rows == nil {
		rows = []types.MetadataWithPlaycounts{}
	}
	return rows, nil
}
