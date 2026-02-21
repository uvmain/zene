package config

import (
	"cmp"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"zene/core/logger"

	"github.com/joho/godotenv"
)

var BaseUrl string
var MusicDirs []string
var PodcastDirectory string
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
var PodcastArtFolder string
var AudioCacheFolder string
var AudioCacheMaxMB int
var AudioCacheMaxDays int
var AdminUsername string
var AdminPassword string
var AdminEmail string
var UserAvatarFolder string
var DefaultBitRate int
var FfprobeConcurrentProcesses int

func LoadConfig() {

	err := godotenv.Load(".env")
	if err != nil {
		logger.Printf("no .env file found, using only environment variables")
	}

	BaseUrl = cmp.Or(os.Getenv("BASE_URL"), "http://localhost:8080")

	musicDirs := cmp.Or(os.Getenv("MUSIC_DIRS"), "./music")

	MusicDirs = strings.Split(musicDirs, ",")
	for i, dir := range MusicDirs {
		MusicDirs[i], _ = filepath.Abs(dir)
	}
	logger.Printf("Using music directories: %v", MusicDirs)

	dataPath := cmp.Or(os.Getenv("DATA_PATH"), "./data")

	DatabaseDirectory = filepath.Join(dataPath, "database")
	AudioCacheFolder = filepath.Join(dataPath, "audio-cache")
	LibraryDirectory = filepath.Join(dataPath, "library")
	TempDirectory = filepath.Join(dataPath, "temp")
	UserAvatarFolder = filepath.Join(dataPath, "avatars")
	ArtworkFolder = filepath.Join(dataPath, "artwork")
	AlbumArtFolder = filepath.Join(ArtworkFolder, "album")
	ArtistArtFolder = filepath.Join(ArtworkFolder, "artist")
	PodcastArtFolder = filepath.Join(ArtworkFolder, "podcasts")

	podcastDirectory := os.Getenv("PODCAST_DIRECTORY")
	if podcastDirectory == "" {
		PodcastDirectory = filepath.Join(dataPath, "podcasts")
	} else {
		PodcastDirectory, _ = filepath.Abs(podcastDirectory)
	}

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

	defaultBitRate := os.Getenv("DEFAULT_BIT_RATE")
	if defaultBitRate == "" {
		DefaultBitRate = 160
	} else {
		defaultBitRateInt, err := strconv.Atoi(defaultBitRate)
		if err != nil {
			DefaultBitRate = 160
		} else {
			DefaultBitRate = defaultBitRateInt
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

	ffprobeConcurrentProcesses := os.Getenv("FFPROBE_CONCURRENT_PROCESSES")
	ffprobeConcurrentProcessesInt, err := strconv.Atoi(ffprobeConcurrentProcesses)
	if err != nil {
		FfprobeConcurrentProcesses = 8 // default to 8 ffprobe concurrent processes
	} else {
		FfprobeConcurrentProcesses = ffprobeConcurrentProcessesInt
	}

	audioFileTypesEnv := cmp.Or(os.Getenv("AUDIO_FILE_TYPES"), ".aac,.alac,.flac,.m4a,.mp3,.ogg,.opus,.wav,.wma")
	AudioFileTypes = strings.Split(audioFileTypesEnv, ",")
	for i, ext := range AudioFileTypes {
		AudioFileTypes[i] = strings.TrimSpace(ext)
	}
	logger.Printf("Audio file types: %v", AudioFileTypes)

	AdminUsername = os.Getenv("ADMIN_USERNAME")
	AdminPassword = os.Getenv("ADMIN_PASSWORD")
	AdminEmail = os.Getenv("ADMIN_EMAIL")
}

func IsLocalDevEnv() bool {
	localDev := os.Getenv("LOCAL_DEV_ENV")
	localDevBool, _ := strconv.ParseBool(localDev)
	return localDevBool
}
