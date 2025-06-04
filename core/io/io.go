package io

import (
	"log"
	"os"
	"time"

	"github.com/djherbis/times"
)

func FileExists(absoluteFilePath string) bool {
	if _, err := os.Stat(absoluteFilePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDir(directoryPath string) {
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		log.Printf("Creating directory: %s", directoryPath)
		err := os.MkdirAll(directoryPath, 0755)
		if err != nil {
			log.Printf("Error creating directory%s: %s", directoryPath, err)
		} else {
			log.Printf("Directory created: %s", directoryPath)
		}
	} else {
		log.Printf("Directory already exists: %s", directoryPath)
	}
}

func GetChangedTime(path string) time.Time {
	t, err := times.Stat(path)
	if err != nil {
		log.Printf("Error retrieving file times for %s: %v", path, err)
		return time.Time{}
	}

	modTime := t.ModTime()
	var changeTime time.Time

	if t.HasChangeTime() {
		changeTime = t.ChangeTime()
	} else {
		changeTime = t.ModTime()
	}

	if changeTime.After(modTime) {
		modTime = changeTime
	}
	return modTime
}
