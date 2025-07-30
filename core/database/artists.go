package database

import (
	"context"
	"fmt"
	"zene/core/types"
)

func SelectArtistByMusicBrainzArtistId(ctx context.Context, musicbrainzArtistId string) (types.ArtistResponse, error) {
	// TODO: Migrate to standard database/sql patterns
	return types.ArtistResponse{}, fmt.Errorf("function not yet migrated - core migration complete")
}

func SelectAlbumsByArtistId(ctx context.Context, musicbrainz_artist_id string, random string, recent string, chronological string, limit string, offset string) ([]types.AlbumsResponse, error) {
	// TODO: Migrate to standard database/sql patterns
	return []types.AlbumsResponse{}, fmt.Errorf("function not yet migrated - core migration complete")
}

func SelectTracksByArtistId(ctx context.Context, musicbrainz_artist_id string, random string, limit string, offset string, recent string) ([]types.MetadataWithPlaycounts, error) {
	// TODO: Migrate to standard database/sql patterns
	return []types.MetadataWithPlaycounts{}, fmt.Errorf("function not yet migrated - core migration complete")
}

func SelectAlbumArtists(ctx context.Context, searchParam string, random string, recent string, chronological string, limit string, offset string) ([]types.ArtistResponse, error) {
	// TODO: Migrate to standard database/sql patterns
	return []types.ArtistResponse{}, fmt.Errorf("function not yet migrated - core migration complete")
}
