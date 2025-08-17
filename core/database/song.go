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

func GetSong(ctx context.Context, musicbrainzTrackId string) (types.SubsonicSong, error) {
	requestUser, err := GetUserByContext(ctx)
	if err != nil {
		return types.SubsonicSong{}, err
	}

	query := `select musicbrainz_track_id as id, musicbrainz_album_id as album_id, title, album, artist, track_number as track,
		substr(release_date,1,4) as year, substr(genre,1,(instr(genre,';')-1)) as genre, musicbrainz_track_id as cover_art,
		size, duration, bitrate, file_path as path, date_added as created, disc_number, musicbrainz_artist_id as artist_id,
		genre, album_artist, bit_depth, sample_rate, channels
		from metadata m
		join user_music_folders f on f.folder_id = m.music_folder_id
		where m.musicbrainz_track_id = ?
		and f.user_id = ?`

	var result types.SubsonicSong
	var genreString string
	var durationFloat float64
	var albumArtist string

	result.IsDir = false
	result.MediaType = "song"
	result.Type = "music"
	result.IsVideo = false
	result.Bpm = 0
	result.Comment = ""
	result.Contributors = []types.SongContributors{}
	result.Moods = []string{}

	err = DB.QueryRowContext(ctx, query, musicbrainzTrackId, requestUser.Id).Scan(
		&result.Id, &result.AlbumId, &result.Title, &result.Album, &result.Artist,
		&result.Track, &result.Year, &result.Genre, &result.CoverArt, &result.Size,
		&durationFloat, &result.BitRate, &result.Path, &result.Created, &result.DiscNumber,
		&result.ArtistId, &genreString, &albumArtist, &result.BitDepth, &result.SamplingRate, &result.ChannelCount,
	)
	if err == sql.ErrNoRows {
		logger.Printf("No song found for %s", musicbrainzTrackId)
		return types.SubsonicSong{}, nil
	} else if err != nil {
		logger.Printf("Error querying song for %s: %v", musicbrainzTrackId, err)
		return types.SubsonicSong{}, err
	}

	result.ContentType = logic.InferContentTypeFromFileExtension(result.Path)
	result.Suffix = strings.Replace(filepath.Ext(result.Path), ".", "", 1)
	result.Duration = int(durationFloat)
	result.Parent = result.AlbumId
	result.SortName = strings.ToLower(result.Title)
	result.MusicBrainzId = result.Id

	result.Genres = []types.SongGenre{}
	for _, genre := range strings.Split(genreString, ";") {
		result.Genres = append(result.Genres, types.SongGenre{Name: genre})
	}

	result.Artists = []types.SongArtist{}
	result.Artists = append(result.Artists, types.SongArtist{Id: result.ArtistId, Name: result.Artist})

	result.DisplayArtist = result.Artist

	result.AlbumArtists = []types.SongArtist{}
	result.AlbumArtists = append(result.AlbumArtists, types.SongArtist{Id: result.ArtistId, Name: albumArtist})

	result.DisplayAlbumArtist = albumArtist

	return result, nil
}
