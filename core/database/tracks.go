package database

import (
	"context"
	"fmt"
	"zene/core/types"
)

func SelectTracksByRandomOffset(ctx context.Context, random string, limit string, offset string) ([]types.MetadataWithPlaycounts, error) {
	// TODO: Migrate to standard database/sql patterns
	return []types.MetadataWithPlaycounts{}, fmt.Errorf("function not yet migrated - core migration complete")
}

func SelectTracksByRecentOffset(ctx context.Context, recent string, limit string, offset string) ([]types.MetadataWithPlaycounts, error) {
	// TODO: Migrate to standard database/sql patterns
	return []types.MetadataWithPlaycounts{}, fmt.Errorf("function not yet migrated - core migration complete")
}

func SelectTrackByMusicBrainzTrackId(ctx context.Context, musicbrainzTrackId string) (types.MetadataWithPlaycounts, error) {
	// TODO: Migrate to standard database/sql patterns
	return types.MetadataWithPlaycounts{}, fmt.Errorf("function not yet migrated - core migration complete")
}
