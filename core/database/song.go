package database

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func GetSong(ctx context.Context, musicbrainzTrackId string) (types.SubsonicChild, error) {
	requestUser, err := GetUserByContext(ctx)
	if err != nil {
		return types.SubsonicChild{}, err
	}

	query := `select musicbrainz_track_id as id, musicbrainz_album_id as album_id, title, album, artist, track_number as track,
		substr(release_date,1,4) as year, substr(genre,1,(instr(genre,';')-1)) as genre, musicbrainz_track_id as cover_art,
		size, duration, bitrate, file_path as path, date_added as created, disc_number, musicbrainz_artist_id as artist_id,
		genre, album_artist, bit_depth, sample_rate, channels
		from metadata m
		join user_music_folders f on f.folder_id = m.music_folder_id
		where m.musicbrainz_track_id = ?
		and f.user_id = ?`

	var result types.SubsonicChild
	var genreString string
	var durationFloat float64
	var albumArtist string

	result.IsDir = false
	result.MediaType = "song"
	result.Type = "music"
	result.IsVideo = false
	result.Bpm = 0
	result.Comment = ""
	result.Contributors = []types.ChildContributors{}
	result.Moods = []string{}

	err = DB.QueryRowContext(ctx, query, musicbrainzTrackId, requestUser.Id).Scan(
		&result.Id, &result.AlbumId, &result.Title, &result.Album, &result.Artist,
		&result.Track, &result.Year, &result.Genre, &result.CoverArt, &result.Size,
		&durationFloat, &result.BitRate, &result.Path, &result.Created, &result.DiscNumber,
		&result.ArtistId, &genreString, &albumArtist, &result.BitDepth, &result.SamplingRate, &result.ChannelCount,
	)
	if err == sql.ErrNoRows {
		logger.Printf("No song found for %s", musicbrainzTrackId)
		return types.SubsonicChild{}, nil
	} else if err != nil {
		logger.Printf("Error querying song for %s: %v", musicbrainzTrackId, err)
		return types.SubsonicChild{}, err
	}

	result.ContentType = logic.InferContentTypeFromFileExtension(result.Path)
	result.Suffix = strings.Replace(filepath.Ext(result.Path), ".", "", 1)
	result.Duration = int(durationFloat)
	result.Parent = result.AlbumId
	result.SortName = strings.ToLower(result.Title)
	result.MusicBrainzId = result.Id

	result.Genres = []types.ChildGenre{}
	for _, genre := range strings.Split(genreString, ";") {
		result.Genres = append(result.Genres, types.ChildGenre{Name: genre})
	}

	result.Artists = []types.ChildArtist{}
	result.Artists = append(result.Artists, types.ChildArtist{Id: result.ArtistId, Name: result.Artist})

	result.DisplayArtist = result.Artist

	result.AlbumArtists = []types.ChildArtist{}
	result.AlbumArtists = append(result.AlbumArtists, types.ChildArtist{Id: result.ArtistId, Name: albumArtist})

	result.DisplayAlbumArtist = albumArtist

	return result, nil
}

func GetSongsForAlbum(ctx context.Context, musicbrainzAlbumId string) ([]types.SubsonicChild, error) {
	requestUser, err := GetUserByContext(ctx)
	if err != nil {
		return []types.SubsonicChild{}, err
	}

	query := `select musicbrainz_track_id as id, musicbrainz_album_id as album_id, title, album, artist, track_number as track,
		substr(release_date,1,4) as year, substr(genre,1,(instr(genre,';')-1)) as genre, musicbrainz_track_id as cover_art,
		size, duration, bitrate, file_path as path, date_added as created, disc_number, musicbrainz_artist_id as artist_id,
		genre, album_artist, bit_depth, sample_rate, channels
		from metadata m
		join user_music_folders f on f.folder_id = m.music_folder_id
		where m.musicbrainz_album_id = ?
		and f.user_id = ?`

	var results []types.SubsonicChild

	rows, err := DB.QueryContext(ctx, query, musicbrainzAlbumId, requestUser.Id)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []types.SubsonicChild{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var child types.SubsonicChild

		var genreString string
		var durationFloat float64
		var albumArtist string

		child.IsDir = false
		child.MediaType = "song"
		child.Type = "music"
		child.IsVideo = false
		child.Bpm = 0
		child.Comment = ""
		child.Contributors = []types.ChildContributors{}
		child.Moods = []string{}

		if err := rows.Scan(&child.Id, &child.AlbumId, &child.Title, &child.Album, &child.Artist,
			&child.Track, &child.Year, &child.Genre, &child.CoverArt, &child.Size,
			&durationFloat, &child.BitRate, &child.Path, &child.Created, &child.DiscNumber,
			&child.ArtistId, &genreString, &albumArtist, &child.BitDepth, &child.SamplingRate, &child.ChannelCount); err != nil {
			logger.Printf("Failed to scan row in SelectTracksByAlbumId: %v", err)
			return []types.SubsonicChild{}, err
		}
		child.ContentType = logic.InferContentTypeFromFileExtension(child.Path)
		child.Suffix = strings.Replace(filepath.Ext(child.Path), ".", "", 1)
		child.Duration = int(durationFloat)
		child.Parent = child.AlbumId
		child.SortName = strings.ToLower(child.Title)
		child.MusicBrainzId = child.Id

		child.Genres = []types.ChildGenre{}
		for _, genre := range strings.Split(genreString, ";") {
			child.Genres = append(child.Genres, types.ChildGenre{Name: genre})
		}

		child.Artists = []types.ChildArtist{}
		child.Artists = append(child.Artists, types.ChildArtist{Id: child.ArtistId, Name: child.Artist})

		child.DisplayArtist = child.Artist

		child.AlbumArtists = []types.ChildArtist{}
		child.AlbumArtists = append(child.AlbumArtists, types.ChildArtist{Id: child.ArtistId, Name: albumArtist})

		child.DisplayAlbumArtist = albumArtist

		results = append(results, child)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return results, err
	}

	return results, nil
}
