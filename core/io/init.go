package io

import (
	"zene/core/config"
	"zene/core/logger"
)

func CreateDirs() {

	dirs := []struct {
		path string
		msg  string
	}{
		{config.DatabaseDirectory, "Database folder already exists"},
		{config.LibraryDirectory, "Library folder already exists"},
		{config.TempDirectory, "Temp folder already exists"},
		{config.UserAvatarFolder, "User avatar folder already exists"},
		{config.ArtworkFolder, "Artwork folder already exists"},
		{config.AlbumArtFolder, "Album artwork folder already exists"},
		{config.ArtistArtFolder, "Artist artwork folder already exists"},
		{config.PodcastArtFolder, "Podcast artwork folder already exists"},
		{config.AudioCacheFolder, "Database folder already exists"},
		{config.PodcastDirectory, "Podcast folder already exists"},
	}

	for _, dir := range dirs {
		if FileExists(dir.path) {
			logger.Printf("[Bootstrap] %s", dir.msg)
		} else {
			CreateDir(dir.path)
		}
	}
}
