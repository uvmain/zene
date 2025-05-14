package config

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var MusicDir string
var DatabaseDirectory string
var FfmpegPath string
var FfprobePath string
var AudioFileTypes []string
var LastFmApiKey string
var ArtworkFolder string
var AlbumArtFolder string
var ArtistArtFolder string

func LoadConfig() {

	godotenv.Load(".env")

	musicDir := os.Getenv("MUSIC_DIR")
	if musicDir == "" {
		musicDir = "./"
	}
	MusicDir, _ = filepath.Abs(musicDir)
	log.Printf("Using music directory: %s", MusicDir)

	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		dataPath = "./data"
	}

	DatabaseDirectory, _ = filepath.Abs(dataPath)

	ArtworkFolder = filepath.Join(dataPath, "artwork")
	AlbumArtFolder = filepath.Join(ArtworkFolder, "album")
	ArtistArtFolder = filepath.Join(ArtworkFolder, "artist")

	ffmpegPath := os.Getenv("FFMPEG_PATH")
	if ffmpegPath == "" {
		FfmpegPath = "ffmpeg"
	} else {
		FfmpegPath, _ = filepath.Abs(ffmpegPath)
	}

	log.Printf("FFMPEG_PATH: %s", FfmpegPath)
	version, err := exec.Command(FfmpegPath, "-version").Output()
	if err != nil {
		log.Printf("ffmpeg not found at %s", FfmpegPath)
	} else {
		log.Printf("ffmpeg version is %v", strings.Split(string(version), "\n")[0])
	}

	ffprobePath := os.Getenv("FFPROBE_PATH")
	if ffprobePath == "" {
		FfprobePath = "ffprobe"
	} else {
		FfprobePath, _ = filepath.Abs(ffprobePath)
	}

	log.Printf("FFPROBE_PATH: %s", FfprobePath)
	version, err = exec.Command(FfprobePath, "-version").Output()
	if err != nil {
		log.Printf("ffprobe not found at %s: %v", FfprobePath, err)
	} else {
		log.Printf("ffprobe version is %v", strings.Split(string(version), "\n")[0])
	}

	audioFileTypesEnv := os.Getenv("AUDIO_FILE_TYPES")
	if audioFileTypesEnv == "" {
		AudioFileTypes = []string{
			".aac", ".alac", ".flac", ".m4a", ".mp3", ".ogg",
		}
	} else {
		AudioFileTypes = strings.Split(audioFileTypesEnv, ",")
		// Trim whitespace from each element (optional but recommended)
		for i, ext := range AudioFileTypes {
			AudioFileTypes[i] = strings.TrimSpace(ext)
		}
	}
	log.Printf("Audio file types: %v", AudioFileTypes)

	lastFmApiKey := os.Getenv("LASTFM_API_KEY")
	if lastFmApiKey == "" {
		log.Println("Warning: last.fm API key not defined")
	} else {
		LastFmApiKey = lastFmApiKey
		log.Println("last.fm API key is defined")
	}
}

func IsLocalDevEnv() bool {
	localDev := os.Getenv("LOCAL_DEV_ENV")
	localDevBool, _ := strconv.ParseBool(localDev)
	return localDevBool
}
