package database

import (
	"context"
	"database/sql"
	"fmt"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func SelectTrack(ctx context.Context, musicBrainzTrackId string) (types.MetadataWithPlaycounts, error) {
	userId, _ := logic.GetUserIdFromContext(ctx)

	query := getUnendedMetadataWithPlaycountsSql(userId)
	query += " where m.musicbrainz_track_id = ? limit 1;"

	var result types.MetadataWithPlaycounts
	err := DbRead.QueryRowContext(ctx, query, musicBrainzTrackId).Scan(&result.FilePath, &result.DateAdded, &result.DateModified, &result.FileName, &result.Format, &result.Duration,
		&result.Size, &result.Bitrate, &result.Title, &result.Artist, &result.Album, &result.AlbumArtist, &result.Genre, &result.TrackNumber,
		&result.TotalTracks, &result.DiscNumber, &result.TotalDiscs, &result.ReleaseDate, &result.MusicBrainzArtistID, &result.MusicBrainzAlbumID,
		&result.MusicBrainzTrackID, &result.Label, &result.MusicFolderId, &result.Codec, &result.BitDepth, &result.SampleRate, &result.Channels,
		&result.UserPlayCount, &result.GlobalPlayCount)
	if err == sql.ErrNoRows {
		return types.MetadataWithPlaycounts{}, nil
	} else if err != nil {
		return types.MetadataWithPlaycounts{}, err
	}
	return result, nil
}

func SelectTrackFilesForScanner(ctx context.Context, musicDir string) ([]types.File, error) {
	query := "SELECT m.file_path, m.file_name, m.date_modified FROM metadata m join music_folders f on m.music_folder_id = f.id where f.name = ?;"

	rows, err := DbRead.QueryContext(ctx, query, musicDir)
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
	err := DbRead.QueryRow(query, artist, title).Scan(&musicBrainzTrackId)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("no track found for artist '%s' and title '%s'", artist, title)
	} else if err != nil {
		return "", fmt.Errorf("error querying track ID: %v", err)
	}
	return musicBrainzTrackId, nil
}
