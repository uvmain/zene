package database

import (
	"context"
	"database/sql"
	"strings"
	"zene/core/logger"
	"zene/core/types"
)

func GetSong(ctx context.Context, musicbrainzTrackId string) (types.SubsonicChild, error) {
	requestUser, err := GetUserByContext(ctx)
	if err != nil {
		return types.SubsonicChild{}, err
	}

	query := `select m.musicbrainz_track_id as id, m.musicbrainz_album_id as album_id, m.title, m.album, m.artist, COALESCE(m.track_number, 0) as track,
		REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year, substr(m.genre,1,(instr(m.genre,';')-1)) as genre, m.musicbrainz_track_id as cover_art,
		m.size, m.duration, m.bitrate, m.file_path as path, m.date_added as created, m.disc_number, m.musicbrainz_artist_id as artist_id,
		m.genre, m.album_artist, m.bit_depth, m.sample_rate, m.channels,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(AVG(gr.rating), 0.0) AS average_rating,
		COALESCE(SUM(pc.play_count), 0) AS play_count,
		max(pc.last_played) as played,
		us.created_at AS starred,
		maa.musicbrainz_artist_id as album_artist_id
	from metadata m
	join user_music_folders f on f.folder_id = m.music_folder_id
	LEFT JOIN user_stars s ON m.musicbrainz_album_id = s.metadata_id AND s.user_id = f.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_album_id = ur.metadata_id AND ur.user_id = f.user_id
	LEFT JOIN user_ratings gr ON m.musicbrainz_album_id = gr.metadata_id
	LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = f.user_id
	LEFT JOIN user_stars us ON m.musicbrainz_track_id = us.metadata_id AND us.user_id = f.user_id
	left join metadata maa on maa.artist = m.album_artist
	where m.musicbrainz_track_id = ?
	and f.user_id = ?
	group by m.musicbrainz_track_id;`

	var result types.SubsonicChild

	var albumArtistName sql.NullString
	var albumArtistId sql.NullString
	var genreString string
	var durationFloat float64
	var played sql.NullString
	var starred sql.NullString

	err = DB.QueryRowContext(ctx, query, musicbrainzTrackId, requestUser.Id).Scan(
		&result.Id, &result.Parent, &result.Title, &result.Album, &result.Artist, &result.Track,
		&result.Year, &result.Genre, &result.CoverArt,
		&result.Size, &durationFloat, &result.BitRate, &result.Path, &result.Created, &result.DiscNumber, &result.ArtistId,
		&genreString, &albumArtistName, &result.BitDepth, &result.SamplingRate, &result.ChannelCount,
		&result.UserRating, &result.AverageRating, &result.PlayCount, &played, &starred, &albumArtistId,
	)
	if err == sql.ErrNoRows {
		logger.Printf("No song found for %s", musicbrainzTrackId)
		return types.SubsonicChild{}, nil
	} else if err != nil {
		logger.Printf("Error querying song for %s: %v", musicbrainzTrackId, err)
		return types.SubsonicChild{}, err
	}
	result.Genres = []types.ChildGenre{}
	for _, genre := range strings.Split(genreString, ";") {
		result.Genres = append(result.Genres, types.ChildGenre{Name: genre})
	}

	if played.Valid {
		result.Played = played.String
	}
	if starred.Valid {
		result.Starred = starred.String
	}

	result.Duration = int(durationFloat)
	result.IsDir = false
	result.MusicBrainzId = result.Id
	result.AlbumId = result.Parent

	result.Artists = []types.ChildArtist{}
	result.Artists = append(result.Artists, types.ChildArtist{Id: result.ArtistId, Name: result.Artist})

	result.DisplayArtist = result.Artist

	result.AlbumArtists = []types.ChildArtist{}
	if albumArtistId.Valid && albumArtistName.Valid {
		result.AlbumArtists = append(result.AlbumArtists, types.ChildArtist{Id: albumArtistId.String, Name: albumArtistName.String})
	}

	result.DisplayAlbumArtist = albumArtistName.String

	return result, nil
}

func GetSongsForAlbum(ctx context.Context, musicbrainzAlbumId string) ([]types.SubsonicChild, error) {
	requestUser, err := GetUserByContext(ctx)
	if err != nil {
		return []types.SubsonicChild{}, err
	}

	query := `select m.musicbrainz_track_id as id, m.musicbrainz_album_id as parent, m.title, m.album, m.artist, COALESCE(m.track_number, 0) as track,
		REPLACE(PRINTF('%4s', substr(m.release_date,1,4)), ' ', '0') as year, substr(m.genre,1,(instr(m.genre,';')-1)) as genre, m.musicbrainz_track_id as cover_art,
		m.size, m.duration, m.bitrate, m.file_path as path, m.date_added as created, m.disc_number, m.musicbrainz_artist_id as artist_id,
		m.genre, m.album_artist, maa.musicbrainz_artist_id as album_artist_id, m.bit_depth, m.sample_rate, m.channels,
		COALESCE(ur.rating, 0) AS user_rating,
		COALESCE(AVG(gr.rating), 0.0) AS average_rating,
		COALESCE(SUM(pc.play_count), 0) AS play_count,
		max(pc.last_played) as played,
		us.created_at AS starred
	from user_music_folders u
	join metadata m on m.music_folder_id = u.folder_id
	LEFT JOIN user_stars us ON m.musicbrainz_album_id = us.metadata_id AND us.user_id = u.user_id
	LEFT JOIN user_ratings ur ON m.musicbrainz_album_id = ur.metadata_id AND ur.user_id = u.user_id
	LEFT JOIN user_ratings gr ON m.musicbrainz_album_id = gr.metadata_id
	LEFT JOIN play_counts pc ON m.musicbrainz_track_id = pc.musicbrainz_track_id AND pc.user_id = u.user_id
	left join metadata maa on maa.artist = m.album_artist
	where m.musicbrainz_album_id = ?
	and u.user_id = ?
	group by m.musicbrainz_track_id;`

	var results []types.SubsonicChild

	rows, err := DB.QueryContext(ctx, query, musicbrainzAlbumId, requestUser.Id)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.SubsonicChild{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var result types.SubsonicChild

		var albumArtistName sql.NullString
		var albumArtistId sql.NullString
		var genreString string
		var durationFloat float64
		var played sql.NullString
		var starred sql.NullString

		if err := rows.Scan(&result.Id, &result.Parent, &result.Title, &result.Album, &result.Artist, &result.Track,
			&result.Year, &result.Genre, &result.CoverArt,
			&result.Size, &durationFloat, &result.BitRate, &result.Path, &result.Created, &result.DiscNumber, &result.ArtistId,
			&genreString, &albumArtistName, &albumArtistId, &result.BitDepth, &result.SamplingRate, &result.ChannelCount,
			&result.UserRating, &result.AverageRating, &result.PlayCount, &played, &starred); err != nil {
			return nil, err
		}
		result.Genres = []types.ChildGenre{}
		for _, genre := range strings.Split(genreString, ";") {
			result.Genres = append(result.Genres, types.ChildGenre{Name: genre})
		}

		if played.Valid {
			result.Played = played.String
		}
		if starred.Valid {
			result.Starred = starred.String
		}

		result.Duration = int(durationFloat)
		result.IsDir = false
		result.MusicBrainzId = result.Id
		result.AlbumId = result.Parent

		result.Artists = []types.ChildArtist{}
		result.Artists = append(result.Artists, types.ChildArtist{Id: result.ArtistId, Name: result.Artist})

		result.DisplayArtist = result.Artist

		result.AlbumArtists = []types.ChildArtist{}
		if albumArtistId.Valid && albumArtistName.Valid {
			result.AlbumArtists = append(result.AlbumArtists, types.ChildArtist{Id: albumArtistId.String, Name: albumArtistName.String})
		}

		result.DisplayAlbumArtist = albumArtistName.String

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return results, err
	}

	return results, nil
}
