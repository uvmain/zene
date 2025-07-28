package io

import (
	"zene/core/config"
	"zene/core/logger"
)

func CreateDirs() {

	if FileExists(config.DatabaseDirectory) {
		logger.Println("Database folder already exists")
	} else {
		CreateDir(config.DatabaseDirectory)
	}

	if FileExists(config.LibraryDirectory) {
		logger.Println("Library folder already exists")
	} else {
		CreateDir(config.LibraryDirectory)
	}

	if FileExists(config.ArtworkFolder) {
		logger.Println("Artwork folder already exists")
	} else {
		CreateDir(config.ArtworkFolder)
	}

	if FileExists(config.AlbumArtFolder) {
		logger.Println("Album artwork folder already exists")
	} else {
		CreateDir(config.AlbumArtFolder)
	}

	if FileExists(config.ArtistArtFolder) {
		logger.Println("Artist artwork folder already exists")
	} else {
		CreateDir(config.ArtistArtFolder)
	}

	if FileExists(config.AudioCacheFolder) {
		logger.Println("Database folder already exists")
	} else {
		CreateDir(config.AudioCacheFolder)
	}
}
