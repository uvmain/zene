package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"zene/core/logger"

	"github.com/joho/godotenv"
)

var MusicDir string
var DatabaseDirectory string
var FfmpegPath string
var FfprobePath string
var AudioFileTypes []string
var ArtworkFolder string
var AlbumArtFolder string
var ArtistArtFolder string
var AudioCacheFolder string
var AudioCacheMaxMB int
var AudioCacheMaxDays int

func LoadConfig() {

	godotenv.Load(".env")

	musicDir := os.Getenv("MUSIC_DIR")
	if musicDir == "" {
		musicDir = "./music"
	}
	MusicDir, _ = filepath.Abs(musicDir)
	logger.Printf("Using music directory: %s", MusicDir)

	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		dataPath = "./data"
	}

	DatabaseDirectory = filepath.Join(dataPath, "database")
	AudioCacheFolder = filepath.Join(dataPath, "audio-cache")
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
		FfmpegPath = "./bin/ffmpeg"
	} else {
		FfmpegPath, _ = filepath.Abs(ffmpegPath)
	}

	logger.Printf("FFMPEG_PATH: %s", FfmpegPath)
	version, err := exec.Command(FfmpegPath, "-version").Output()
	if err != nil {
		logger.Printf("ffmpeg not found at %s", FfmpegPath)
	} else {
		logger.Printf("ffmpeg version is %v", strings.Split(string(version), "\n")[0])
	}

	ffprobePath := os.Getenv("FFPROBE_PATH")
	if ffprobePath == "" {
		FfprobePath = "./bin/ffprobe"
	} else {
		FfprobePath, _ = filepath.Abs(ffprobePath)
	}

	logger.Printf("FFPROBE_PATH: %s", FfprobePath)
	version, err = exec.Command(FfprobePath, "-version").Output()
	if err != nil {
		logger.Printf("ffprobe not found at %s: %v", FfprobePath, err)
	} else {
		logger.Printf("ffprobe version is %v", strings.Split(string(version), "\n")[0])
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
}

func IsLocalDevEnv() bool {
	localDev := os.Getenv("LOCAL_DEV_ENV")
	localDevBool, _ := strconv.ParseBool(localDev)
	return localDevBool
}
