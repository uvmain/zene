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

func SelectArtistByMusicBrainzArtistId(ctx context.Context, musicbrainzArtistId string) (types.ArtistResponse, error) {
	query := `SELECT DISTINCT artist, musicbrainz_artist_id FROM metadata where musicbrainz_artist_id = ? limit 1;`

	var result types.ArtistResponse

	err := DB.QueryRowContext(ctx, query, musicbrainzArtistId).Scan(&result.Artist, &result.MusicBrainzArtistID)

	if err == sql.ErrNoRows {
		return types.ArtistResponse{}, nil
	} else if err != nil {
		return types.ArtistResponse{}, err
	}

	return result, nil
}

func SelectAlbumsByArtistId(ctx context.Context, musicbrainz_artist_id string, random string, recent string, chronological string, limit string, offset string) ([]types.AlbumsResponse, error) {
	query := "SELECT DISTINCT musicbrainz_album_id, album, musicbrainz_artist_id, artist, genre, release_date FROM metadata WHERE musicbrainz_artist_id = $musicbrainz_artist_id GROUP BY musicbrainz_album_id"

	if recent == "true" {
		query += " ORDER BY date_added desc"
	} else if random != "" {
		randomInt, err := strconv.Atoi(random)
		if err == nil {
			query += fmt.Sprintf(" ORDER BY ((rowid * %d) %% 1000000)", randomInt)
		}
	} else if chronological == "true" {
		query += " ORDER BY release_date desc"
	}

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err == nil {
			query = fmt.Sprintf("%s limit %d", query, limitInt)
		}
	}
	if offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err == nil {
			query = fmt.Sprintf("%s offset %d", query, offsetInt)
		}
	}

	query = fmt.Sprintf("%s;", query)

	rows, err := DB.QueryContext(ctx, query, musicbrainz_artist_id)
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

func SelectTracksByArtistId(ctx context.Context, musicbrainz_artist_id string, random string, limit string, offset string, recent string) ([]types.MetadataWithPlaycounts, error) {
	userId, _ := logic.GetUserIdFromContext(ctx)
	query := getUnendedMetadataWithPlaycountsSql(userId)

	query = fmt.Sprintf("%s where musicbrainz_artist_id = $musicbrainz_artist_id", query)

	if recent == "true" {
		query = fmt.Sprintf("%s ORDER BY date_added desc", query)
	} else if random != "" {
		randomInteger, err := strconv.Atoi(random)
		if err == nil {
			query += fmt.Sprintf(" ORDER BY ((rowid * %d) %% 1000000)", randomInteger)
		}
	}
	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err == nil {
			query = fmt.Sprintf("%s limit %d", query, limitInt)
		}
	}
	if offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err == nil {
			query = fmt.Sprintf("%s offset %d", query, offsetInt)
		}
	}

	query = fmt.Sprintf("%s;", query)

	rows, err := DB.QueryContext(ctx, query, musicbrainz_artist_id)
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

func SelectAlbumArtists(ctx context.Context, searchParam string, random string, recent string, chronological string, limit string, offset string) ([]types.ArtistResponse, error) {
	query := "select distinct m.album_artist, m.musicbrainz_artist_id FROM metadata m"
	if searchParam != "" {
		query += " JOIN artists_fts f ON m.file_path = f.file_path"
	}
	query = " where m.album_artist = m.artist"

	if searchParam != "" {
		query += " and artists_fts MATCH ?"
	}

	if recent == "true" {
		query += " ORDER BY m.date_added desc"
	} else if random != "" {
		randomInt, err := strconv.Atoi(random)
		if err == nil {
			query += fmt.Sprintf(" ORDER BY ((m.rowid * %d) %% 1000000)", randomInt)
		}
	} else if chronological == "true" {
		query += " ORDER BY m.release_date desc"
	}

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err == nil {
			query += fmt.Sprintf(" limit %d", limitInt)
		}
	}
	if offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err == nil {
			query += fmt.Sprintf(" offset %d", offsetInt)
		}
	}

	query += ";"

	var rows *sql.Rows

	if searchParam != "" {
		rows, err = DB.QueryContext(ctx, query, searchParam)
	} else {
		rows, err = DB.QueryContext(ctx, query)
	}

	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.ArtistResponse{}, err
	}
	defer rows.Close()

	var results []types.ArtistResponse

	for rows.Next() {
		var result types.ArtistResponse
		if err := rows.Scan(&result.Artist, &result.MusicBrainzArtistID); err != nil {
			logger.Printf("Failed to scan row: %v", err)
			return []types.ArtistResponse{}, err
		}
		result.ImageURL = fmt.Sprintf("/api/artists/%s/art", result.MusicBrainzArtistID)
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return results, err
	}

	return results, nil
}
