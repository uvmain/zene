package scanner

import (
	"context"
	"fmt"
	"zene/core/database"
	"zene/core/deezer"
	"zene/core/logger"
)

func PopulateTopSongsTable(ctx context.Context) error {
	logger.Printf("Populating top_songs table")
	_, err := database.DB.ExecContext(ctx, "delete from top_songs where musicbrainz_artist_id not in (select distinct musicbrainz_artist_id from metadata);")
	if err != nil {
		return fmt.Errorf("cleaning top_songs table: %v", err)
	}

	var artistsToCheck []ArtistsToCheck
	query := "SELECT distinct musicbrainz_artist_id, artist FROM metadata where musicbrainz_artist_id not in (select distinct musicbrainz_artist_id from top_songs);"

	rows, err := database.DB.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("querying metadata: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var artist ArtistsToCheck
		if err := rows.Scan(&artist.MusicBrainzId, &artist.ArtistName); err != nil {
			return fmt.Errorf("scanning row: %v", err)
		}
		artistsToCheck = append(artistsToCheck, artist)
	}

	logger.Printf("Found %d artists to check for top songs", len(artistsToCheck))

	for _, artist := range artistsToCheck {
		topSongs, err := deezer.GetTopSongs(ctx, artist.ArtistName, 100)
		if err != nil {
			logger.Printf("Error fetching top songs for %s: %v", artist.ArtistName, err)
		} else {
			if len(topSongs) > 0 {
				logger.Printf("Found %d top songs for %s, linking existing songs", len(topSongs), artist.ArtistName)
			} else {
				logger.Printf("No top songs found for %s", artist.ArtistName)
			}
		}

		err = database.InsertTopSongs(ctx, topSongs)
		if err != nil {
			logger.Printf("Error inserting top songs for %s: %v", artist.ArtistName, err)
		}
	}

	return nil
}
