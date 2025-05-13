package art

import (
	"log"
	"path/filepath"
	"zene/config"
	"zene/io"
)

func Initialise() {
	artworkFolder := filepath.Join(config.DatabaseDirectory, "artwork")

	if io.FileExists(artworkFolder) {
		log.Println("Artwork folder already exists")
	} else {
		io.CreateDir(artworkFolder)
	}
}
