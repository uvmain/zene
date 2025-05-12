package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

var MusicDir string
var DatabaseDirectory string

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
}

func IsLocalDevEnv() bool {
	localDev := os.Getenv("LOCAL_DEV_ENV")
	localDevBool, _ := strconv.ParseBool(localDev)
	return localDevBool
}
