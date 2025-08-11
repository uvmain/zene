package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func SelectAllTracks(ctx context.Context, random string, limit string, offset string, recent string, chronological string) ([]types.MetadataWithPlaycounts, error) {
	userId, _ := logic.GetUserIdFromContext(ctx)
	query := getUnendedMetadataWithPlaycountsSql(userId)

	if recent == "true" {
		query += " ORDER BY m.date_added desc"
	} else if random != "" {
		randomInteger, err := strconv.Atoi(random)
		if err == nil {
			query += fmt.Sprintf(" ORDER BY ((m.rowid * %d) %% 1000000)", randomInteger)
		} else {
			logger.Printf("Error setting randomness: %v", err)
		}
	} else if chronological == "true" {
		query += " ORDER BY m.release_date desc"
	}

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return []types.MetadataWithPlaycounts{}, fmt.Errorf("invalid limit value: %v", err)
		}
		query += fmt.Sprintf(" limit %d", limitInt)
	}

	if offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return []types.MetadataWithPlaycounts{}, fmt.Errorf("invalid offset value: %v", err)
		}
		query += fmt.Sprintf(" offset %d", offsetInt)
	}

	query += ";"

	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.MetadataWithPlaycounts{}, err
	}
	defer rows.Close()

	var results []types.MetadataWithPlaycounts

	for rows.Next() {
		var result types.MetadataWithPlaycounts
		if err := rows.Scan(&result.FilePath, &result.DateAdded, &result.DateModified, &result.FileName, &result.Format, &result.Duration,
			&result.Size, &result.Bitrate, &result.Title, &result.Artist, &result.Album, &result.AlbumArtist, &result.Genre, &result.TrackNumber,
			&result.TotalTracks, &result.DiscNumber, &result.TotalDiscs, &result.ReleaseDate, &result.MusicBrainzArtistID, &result.MusicBrainzAlbumID,
			&result.MusicBrainzTrackID, &result.Label, &result.MusicFolderId, &result.UserPlayCount, &result.GlobalPlayCount); err != nil {
			logger.Printf("Failed to scan row in SelectAllTracks: %v", err)
			return []types.MetadataWithPlaycounts{}, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return results, err
	}

	return results, nil
}

func SelectTrack(ctx context.Context, musicBrainzTrackId string) (types.MetadataWithPlaycounts, error) {
	userId, _ := logic.GetUserIdFromContext(ctx)

	query := getUnendedMetadataWithPlaycountsSql(userId)
	query += " where m.musicbrainz_track_id = ? limit 1;"

	var result types.MetadataWithPlaycounts
	err := DB.QueryRowContext(ctx, query, musicBrainzTrackId).Scan(&result.FilePath, &result.DateAdded, &result.DateModified, &result.FileName, &result.Format, &result.Duration,
		&result.Size, &result.Bitrate, &result.Title, &result.Artist, &result.Album, &result.AlbumArtist, &result.Genre, &result.TrackNumber,
		&result.TotalTracks, &result.DiscNumber, &result.TotalDiscs, &result.ReleaseDate, &result.MusicBrainzArtistID, &result.MusicBrainzAlbumID,
		&result.MusicBrainzTrackID, &result.Label, &result.MusicFolderId, &result.UserPlayCount, &result.GlobalPlayCount)
	if err == sql.ErrNoRows {
		return types.MetadataWithPlaycounts{}, nil
	} else if err != nil {
		return types.MetadataWithPlaycounts{}, err
	}
	return result, nil
}

func SelectTrackFilesForScanner(ctx context.Context) ([]types.File, error) {
	query := "SELECT file_path, file_name, date_modified FROM metadata;"

	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.File{}, err
	}
	defer rows.Close()

	var results []types.File

	for rows.Next() {
		var result types.File
		if err := rows.Scan(&result.FilePathAbs, &result.FileName, &result.DateModified); err != nil {
			logger.Printf("Failed to scan row in SelectTrackFilesForScanner: %v", err)
			return []types.File{}, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return results, err
	}

	return results, nil
}

func GetTrackIdByArtistAndTitle(artist string, title string) (string, error) {
	query := "SELECT musicbrainz_track_id FROM metadata WHERE lower(artist) = lower(?) AND lower(title) = lower(?) LIMIT 1;"
	var musicBrainzTrackId string
	err := DB.QueryRow(query, artist, title).Scan(&musicBrainzTrackId)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("no track found for artist '%s' and title '%s'", artist, title)
	} else if err != nil {
		return "", fmt.Errorf("error querying track ID: %v", err)
	}
	return musicBrainzTrackId, nil
}
