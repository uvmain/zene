package io

import (
	"log"
	"zene/core/config"
)

func CreateDirs() {

	if FileExists(config.DatabaseDirectory) {
		log.Println("Database folder already exists")
	} else {
		CreateDir(config.DatabaseDirectory)
	}

	if FileExists(config.ArtworkFolder) {
		log.Println("Artwork folder already exists")
	} else {
		CreateDir(config.ArtworkFolder)
	}

	if FileExists(config.AlbumArtFolder) {
		log.Println("Album artwork folder already exists")
	} else {
		CreateDir(config.AlbumArtFolder)
	}

	if FileExists(config.ArtistArtFolder) {
		log.Println("Artist artwork folder already exists")
	} else {
		CreateDir(config.ArtistArtFolder)
	}

	if FileExists(config.AudioCacheFolder) {
		log.Println("Database folder already exists")
	} else {
		CreateDir(config.AudioCacheFolder)
	}
}
