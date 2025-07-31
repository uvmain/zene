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

func SelectTracksByAlbumId(ctx context.Context, musicbrainz_album_id string) ([]types.MetadataWithPlaycounts, error) {
	userId, _ := logic.GetUserIdFromContext(ctx)
	query := getUnendedMetadataWithPlaycountsSql(userId)

	query += " where musicbrainz_album_id = ? order by cast(disc_number AS INTEGER), cast(track_number AS INTEGER);"

	rows, err := DB.QueryContext(ctx, query, musicbrainz_album_id)
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
			&result.MusicBrainzTrackID, &result.Label, &result.UserPlayCount, &result.GlobalPlayCount); err != nil {
			logger.Printf("Failed to scan row: %v", err)
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

func SelectAllAlbums(ctx context.Context, random string, limit string, recent string) ([]types.AlbumsResponse, error) {
	query := "SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date, date_added FROM metadata group by album"

	if recent == "true" {
		query += " ORDER BY date_added desc"
	} else if random != "" {
		randomInt, err := strconv.Atoi(random)
		if err == nil {
			query += fmt.Sprintf(" ORDER BY ((rowid * %d) %% 1000000)", randomInt)
		}
	} else {
		query += " ORDER BY album"
	}

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err == nil {
			query += fmt.Sprintf(" limit %d", limitInt)
		}
	}

	query += ";"

	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.AlbumsResponse{}, err
	}
	defer rows.Close()

	var results []types.AlbumsResponse

	for rows.Next() {
		var result types.AlbumsResponse
		if err := rows.Scan(&result.Album, &result.MusicBrainzAlbumID, &result.Artist, &result.MusicBrainzArtistID, &result.Genres, &result.ReleaseDate); err != nil {
			logger.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return results, err
	}

	return results, nil
}

func SelectAlbum(ctx context.Context, musicbrainzAlbumId string) (types.AlbumsResponse, error) {
	query := `SELECT album, album_artist, musicbrainz_album_id, musicbrainz_artist_id, genre, release_date FROM metadata where musicbrainz_album_id = ? limit 1;`

	var result types.AlbumsResponse

	err := DB.QueryRowContext(ctx, query, musicbrainzAlbumId).Scan(&result.Album, &result.Artist, &result.MusicBrainzAlbumID, &result.MusicBrainzArtistID, &result.Genres, &result.ReleaseDate)
	if err == sql.ErrNoRows {
		return types.AlbumsResponse{}, nil
	} else if err != nil {
		return types.AlbumsResponse{}, err
	}
	return result, nil
}
