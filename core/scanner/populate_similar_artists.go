package scanner

import (
	"context"
	"fmt"
	"zene/core/database"
	"zene/core/deezer"
	"zene/core/logger"
)

type ArtistsToCheck struct {
	MusicBrainzId string
	ArtistName    string
}

func PopulateSimilarArtistsTable(ctx context.Context) error {
	logger.Printf("Populating similar_artists table")
	// clear the existing similar_artists contents
	_, err := database.DB.ExecContext(ctx, "delete from similar_artists where artist_id not in (select distinct musicbrainz_artist_id from metadata);")
	if err != nil {
		return fmt.Errorf("cleaning similar_artists table: %v", err)
	}

	// get artists to check
	var artistsToCheck []ArtistsToCheck
	query := "SELECT distinct musicbrainz_artist_id, artist FROM metadata where musicbrainz_artist_id not in (select distinct artist_id from similar_artists);"

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

	logger.Printf("Found %d artists to check for similar artists", len(artistsToCheck))

	// for each artist fetch the similar artists and insert
	for _, artist := range artistsToCheck {
		similarArtists, err := deezer.GetSimilarArtistNames(ctx, artist.ArtistName)
		if err != nil {
			logger.Printf("Error fetching similar artists for %s: %v", artist.ArtistName, err)
		} else {
			if len(similarArtists) > 0 {
				logger.Printf("Found %d similar artists for %s, linking existing artists", len(similarArtists), artist.ArtistName)
			} else {
				logger.Printf("No similar artists found for %s", artist.ArtistName)
			}
		}

		sortOrder := 1
		for _, similarArtistName := range similarArtists {
			similarArtistId, err := database.GetArtistIdByName(ctx, similarArtistName)
			if err == nil && similarArtistId != "" {
				if err := database.InsertSimilarArtistsRow(ctx, artist.MusicBrainzId, similarArtistId, sortOrder); err != nil {
					logger.Printf("error inserting similar artist row: %v", err)
				} else {
					sortOrder++
				}
			}
		}
	}

	return nil
}
