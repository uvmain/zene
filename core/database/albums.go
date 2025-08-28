package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"

	"github.com/timematic/anytime"
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
		if err := rows.Scan(&result.FilePath, &result.FileName, &result.DateAdded, &result.DateModified, &result.Format, &result.Duration,
			&result.Size, &result.Bitrate, &result.Title, &result.Artist, &result.Album, &result.AlbumArtist, &result.Genre, &result.TrackNumber,
			&result.TotalTracks, &result.DiscNumber, &result.TotalDiscs, &result.ReleaseDate, &result.MusicBrainzArtistID, &result.MusicBrainzAlbumID,
			&result.MusicBrainzTrackID, &result.Label, &result.MusicFolderId, &result.Codec, &result.BitDepth, &result.SampleRate, &result.Channels,
			&result.UserPlayCount, &result.GlobalPlayCount); err != nil {
			logger.Printf("Failed to scan row in SelectTracksByAlbumId: %v", err)
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

func SelectAlbumIdByTrackId(ctx context.Context, musicbrainz_track_id string) (string, error) {
	var albumId string
	query := "SELECT musicbrainz_album_id FROM metadata WHERE musicbrainz_track_id = ? limit 1"
	err := DB.QueryRowContext(ctx, query, musicbrainz_track_id).Scan(&albumId)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return "", err
	}
	return albumId, nil
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
		var dateAdded string
		if err := rows.Scan(&result.Album, &result.MusicBrainzAlbumID, &result.Artist, &result.MusicBrainzArtistID, &result.Genres, &result.ReleaseDate, &dateAdded); err != nil {
			logger.Printf("Failed to scan row in SelectAllAlbums: %v", err)
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

func SelectAllAlbumsForMusicDir(ctx context.Context, musicDir string, random string, limit string, recent string) ([]types.AlbumsResponse, error) {
	query := "SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date, date_added FROM metadata m JOIN music_folders f ON m.music_folder_id = f.id WHERE f.name = ? group by album"

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

	rows, err := DB.QueryContext(ctx, query, musicDir)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.AlbumsResponse{}, err
	}
	defer rows.Close()

	var results []types.AlbumsResponse

	for rows.Next() {
		var result types.AlbumsResponse
		var dateAdded string
		if err := rows.Scan(&result.Album, &result.MusicBrainzAlbumID, &result.Artist, &result.MusicBrainzArtistID, &result.Genres, &result.ReleaseDate, &dateAdded); err != nil {
			logger.Printf("Failed to scan row in SelectAllAlbums: %v", err)
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

func GetAlbum(ctx context.Context, musicbrainzAlbumId string) (types.AlbumId3, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return types.AlbumId3{}, err
	}

	var album types.AlbumId3

	query := `select m.musicbrainz_album_id as id,
		m.album as name,
		m.artist as artist,
		m.musicbrainz_album_id as cover_art,
		count(m.musicbrainz_track_id) as song_count,
		cast(sum(m.duration) as integer) as duration,
		COALESCE(SUM(pc.play_count), 0) as play_count,
		min(m.date_added) as created,
		m.musicbrainz_artist_id as artist_id,
		s.created_at as starred,
		REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year,
		substr(m.genre,1,(instr(m.genre,';')-1)) as genre,
		max(pc.last_played) as played,
		COALESCE(ur.rating, 0) AS user_rating,
		m.label as label_string,
		m.musicbrainz_album_id as musicbrainz_id,
		m.genre as genre_string,
		m.artist as display_artist,
		lower(m.album) as sort_name,
		m.release_date as release_date_string
	from user_music_folders f
	join metadata m on m.music_folder_id = f.folder_id
	LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = f.user_id
	LEFT JOIN user_stars s ON m.musicbrainz_album_id = s.metadata_id AND s.user_id = f.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_artist_id = ur.metadata_id AND ur.user_id = f.user_id
	where m.musicbrainz_album_id = ?
	and f.user_id = ?`

	var starred sql.NullString
	var labelString sql.NullString
	var genresString sql.NullString
	var releaseDateString sql.NullString
	var played sql.NullString

	err = DB.QueryRowContext(ctx, query, musicbrainzAlbumId, user.Id).Scan(
		&album.Id, &album.Name, &album.Artist, &album.CoverArt, &album.SongCount,
		&album.Duration, &album.PlayCount, &album.Created, &album.ArtistId, &starred,
		&album.Year, &album.Genre, &played, &album.UserRating,
		&labelString, &album.MusicBrainzId, &genresString,
		&album.DisplayArtist, &album.SortName, &releaseDateString,
	)

	if err == sql.ErrNoRows {
		return types.AlbumId3{}, nil
	} else if err != nil {
		return types.AlbumId3{}, err
	}

	if starred.Valid {
		album.Starred = starred.String
	}

	if played.Valid {
		album.Played = played.String
	}

	album.RecordLabels = []types.ChildRecordLabel{}
	album.RecordLabels = append(album.RecordLabels, types.ChildRecordLabel{Name: labelString.String})

	album.Genres = []types.ItemGenre{}
	for _, genre := range strings.Split(genresString.String, ";") {
		album.Genres = append(album.Genres, types.ItemGenre{Name: genre})
	}

	releaseDateTime, err := anytime.Parse(releaseDateString.String)
	if err == nil {
		album.ReleaseDate = types.ItemDate{
			Year:  releaseDateTime.Year(),
			Month: int(releaseDateTime.Month()),
			Day:   releaseDateTime.Day(),
		}
	}

	album.Songs = []types.SubsonicChild{}

	return album, nil
}
