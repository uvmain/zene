package database

import (
	"context"
	"zene/core/logger"
)

type getArtistsLine struct {
	Artist              string
	MusicBrainzArtistID string
	IsAlbumArtist       bool
}

func SelectArtistsForMusicDir(ctx context.Context, musicDir string) ([]getArtistsLine, error) {
	query := `select distinct m.artist,
		m.musicbrainz_artist_id,
		CASE WHEN m.artist = m.album_artist THEN true ELSE false end as is_album_artist
	FROM metadata m
	join music_folders mf on m.music_folder_id = mf.id and mf.name = ?`

	var results []getArtistsLine

	rows, err := DB.QueryContext(ctx, query, musicDir)
	if err != nil {
		logger.Printf("Query failed: %v", err)
		return []getArtistsLine{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var result getArtistsLine
		if err := rows.Scan(&result.Artist, &result.MusicBrainzArtistID, &result.IsAlbumArtist); err != nil {
			logger.Printf("Failed to scan row in SelectArtistsForMusicDir: %v", err)
			return []getArtistsLine{}, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Rows iteration error: %v", err)
		return results, err
	}

	return results, nil
}
