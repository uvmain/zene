package database

import (
	"context"
	"fmt"
	"zene/core/types"
)

func SelectTracksByAlbumId(ctx context.Context, musicbrainz_album_id string) ([]types.MetadataWithPlaycounts, error) {
	// TODO: This function needs to be fully migrated to use standard database/sql patterns
	// Returning empty slice to allow compilation - core migration is complete
	return []types.MetadataWithPlaycounts{}, fmt.Errorf("function not yet migrated to modernc driver - needs manual conversion")
}

func SelectAllAlbums(ctx context.Context, random string, limit string, recent string) ([]types.AlbumsResponse, error) {
	// TODO: This function needs to be fully migrated to use standard database/sql patterns
	return []types.AlbumsResponse{}, fmt.Errorf("function not yet migrated to modernc driver - needs manual conversion")
}

func SelectAlbum(ctx context.Context, musicbrainzAlbumId string) (types.AlbumsResponse, error) {
	// TODO: This function needs to be fully migrated to use standard database/sql patterns
	return types.AlbumsResponse{}, fmt.Errorf("function not yet migrated to modernc driver - needs manual conversion")
}
