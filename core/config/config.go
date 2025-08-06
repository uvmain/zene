package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"zene/core/logger"

	"github.com/joho/godotenv"
)

var MusicDirs []string
var DatabaseDirectory string
var LibraryDirectory string
var TempDirectory string
var FfmpegPath string
var FfmpegBinaryName string
var FfprobePath string
var FfprobeBinaryName string
var AudioFileTypes []string
var ArtworkFolder string
var AlbumArtFolder string
var ArtistArtFolder string
var AudioCacheFolder string
var AudioCacheMaxMB int
var AudioCacheMaxDays int
var AdminUsername string
var AdminPassword string

func LoadConfig() {

	godotenv.Load(".env")

	musicDirEnv := os.Getenv("MUSIC_DIRS")
	if musicDirEnv == "" {
		MusicDirs = []string{"./music"}
	} else {
		MusicDirs = strings.Split(musicDirEnv, ",")
		for i, dir := range MusicDirs {
			MusicDirs[i], _ = filepath.Abs(dir)
		}
	}
	logger.Printf("Using music directories: %v", MusicDirs)

	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		dataPath = "./data"
	}

	DatabaseDirectory = filepath.Join(dataPath, "database")
	AudioCacheFolder = filepath.Join(dataPath, "audio-cache")
	LibraryDirectory = filepath.Join(dataPath, "library")
	TempDirectory = filepath.Join(dataPath, "temp")
	ArtworkFolder = filepath.Join(dataPath, "artwork")
	AlbumArtFolder = filepath.Join(ArtworkFolder, "album")
	ArtistArtFolder = filepath.Join(ArtworkFolder, "artist")

	audioCacheMaxMB := os.Getenv("AUDIO_CACHE_MAX_MB")
	if audioCacheMaxMB == "" {
		AudioCacheMaxMB = 500
	} else {
		audioCacheMaxMbInt, err := strconv.Atoi(audioCacheMaxMB)
		if err != nil {
			AudioCacheMaxMB = 500
		} else {
			AudioCacheMaxMB = audioCacheMaxMbInt
		}
	}

	audioCacheMaxDays := os.Getenv("AUDIO_CACHE_MAX_DAYS")
	if audioCacheMaxDays == "" {
		AudioCacheMaxDays = 30
	} else {
		audioCacheMaxDaysInt, err := strconv.Atoi(audioCacheMaxMB)
		if err != nil {
			AudioCacheMaxDays = 30
		} else {
			AudioCacheMaxDays = audioCacheMaxDaysInt
		}
	}

	ffmpegPath := os.Getenv("FFMPEG_PATH")
	if ffmpegPath == "" {
		FfmpegBinaryName := "ffmpeg"
		if runtime.GOOS == "windows" {
			FfmpegBinaryName += ".exe"
		}
		FfmpegPath = filepath.Join(LibraryDirectory, FfmpegBinaryName)
	} else {
		FfmpegPath, _ = filepath.Abs(ffmpegPath)
	}

	ffprobePath := os.Getenv("FFPROBE_PATH")
	if ffprobePath == "" {
		FfprobeBinaryName := "ffprobe"
		if runtime.GOOS == "windows" {
			FfprobeBinaryName += ".exe"
		}
		FfprobePath = filepath.Join(LibraryDirectory, FfprobeBinaryName)
	} else {
		FfprobePath, _ = filepath.Abs(ffprobePath)
	}

	audioFileTypesEnv := os.Getenv("AUDIO_FILE_TYPES")
	if audioFileTypesEnv == "" {
		AudioFileTypes = []string{
			".aac", ".alac", ".flac", ".m4a", ".mp3", ".ogg", ".opus", ".wav", ".wma",
		}
	} else {
		AudioFileTypes = strings.Split(audioFileTypesEnv, ",")
		for i, ext := range AudioFileTypes {
			AudioFileTypes[i] = strings.TrimSpace(ext)
		}
	}
	logger.Printf("Audio file types: %v", AudioFileTypes)

	AdminUsername = os.Getenv("ADMIN_USERNAME")
	AdminPassword = os.Getenv("ADMIN_PASSWORD")
}

func IsLocalDevEnv() bool {
	localDev := os.Getenv("LOCAL_DEV_ENV")
	localDevBool, _ := strconv.ParseBool(localDev)
	return localDevBool
}
