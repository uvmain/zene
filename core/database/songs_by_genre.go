package database

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"zene/core/logic"
	"zene/core/types"
)

func GetSongsByGenre(ctx context.Context, genre string, count int, offset int) ([]types.SubsonicSong, error) {
	requestUser, err := GetUserByContext(ctx)
	if err != nil {
		return []types.SubsonicSong{}, err
	}

	query := `select musicbrainz_track_id as id, musicbrainz_album_id as album_id, title, album, artist, track_number as track,
		substr(release_date,1,4) as year, substr(m.genre,1,(instr(m.genre,';')-1)) as genre, musicbrainz_track_id as cover_art,
		size, duration, bitrate, m.file_path as path, date_added as created, disc_number, musicbrainz_artist_id as artist_id,
		m.genre as genres, album_artist, bit_depth, sample_rate, channels
		from metadata m
		join user_music_folders f on f.folder_id = m.music_folder_id
		join track_genres g on m.file_path = g.file_path
		where f.user_id = ? and lower(g.genre) = lower(?)
		limit ? offset ?`

	rows, err := DB.QueryContext(ctx, query, requestUser.Id, genre, count, offset)
	if err != nil {
		return nil, fmt.Errorf("getting scans: %v", err)
	}
	defer rows.Close()

	var songs []types.SubsonicSong
	for rows.Next() {
		var song types.SubsonicSong

		var genreString string
		var durationFloat float64
		var albumArtist string

		song.IsDir = false
		song.MediaType = "song"
		song.Type = "music"
		song.IsVideo = false
		song.Bpm = 0
		song.Comment = ""
		song.Contributors = []types.SongContributors{}
		song.Moods = []string{}

		if err := rows.Scan(&song.Id, &song.AlbumId, &song.Title, &song.Album, &song.Artist,
			&song.Track, &song.Year, &song.Genre, &song.CoverArt, &song.Size,
			&durationFloat, &song.BitRate, &song.Path, &song.Created, &song.DiscNumber,
			&song.ArtistId, &genreString, &albumArtist, &song.BitDepth, &song.SamplingRate, &song.ChannelCount,
		); err != nil {
			return nil, fmt.Errorf("scanning song row: %v", err)
		}

		song.ContentType = logic.InferContentTypeFromFileExtension(song.Path)
		song.Suffix = strings.Replace(filepath.Ext(song.Path), ".", "", 1)
		song.Duration = int(durationFloat)
		song.Parent = song.AlbumId
		song.SortName = strings.ToLower(song.Title)
		song.MusicBrainzId = song.Id

		song.Genres = []types.SongGenre{}
		for _, genre := range strings.Split(genreString, ";") {
			song.Genres = append(song.Genres, types.SongGenre{Name: genre})
		}

		song.Artists = []types.SongArtist{}
		song.Artists = append(song.Artists, types.SongArtist{Id: song.ArtistId, Name: song.Artist})

		song.DisplayArtist = song.Artist

		song.AlbumArtists = []types.SongArtist{}
		song.AlbumArtists = append(song.AlbumArtists, types.SongArtist{Id: song.ArtistId, Name: albumArtist})

		song.DisplayAlbumArtist = albumArtist

		songs = append(songs, song)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating song rows: %v", err)
	}

	return songs, nil
}
