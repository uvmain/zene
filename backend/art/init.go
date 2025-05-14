package art

import (
	"log"
	"zene/config"
	"zene/io"
)

func Initialise() {
	if io.FileExists(config.ArtworkFolder) {
		log.Println("Artwork folder already exists")
	} else {
		io.CreateDir(config.ArtworkFolder)
	}

	if io.FileExists(config.AlbumArtFolder) {
		log.Println("Album artwork folder already exists")
	} else {
		io.CreateDir(config.AlbumArtFolder)
	}

	if io.FileExists(config.ArtistArtFolder) {
		log.Println("Artist artwork folder already exists")
	} else {
		io.CreateDir(config.ArtistArtFolder)
	}
}
