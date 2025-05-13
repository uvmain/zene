package art

import (
	"log"
	"path/filepath"
	"slices"
	"zene/database"
)

func GetArtForAlbum(music_brainz_album_id string) {
	trackMetadataRows, err := database.SelectMetadataByAlbumID(music_brainz_album_id)
	if err != nil {
		log.Printf("Error getting track data from database: %v", err)
	}
	directories := []string{}

	for _, trackMetadata := range trackMetadataRows {
		directory := filepath.Dir(trackMetadata.Filename)
		if !slices.Contains(directories, directory) {
			directories = append(directories, directory)
		}
	}
	directories = slices.Compact(directories)
	log.Printf("Checking album art for %s in %v", music_brainz_album_id, directories)
}
