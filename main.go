package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"zene/config"
	"zene/database"
	"zene/router"
)

func main() {
	config.LoadConfig()

	database.Initialise()
	defer database.CloseDatabase()

	router.CreateRoutes()

	err := scanDirectory(config.MusicDir)
	if err != nil {
		log.Fatalf("Failed to scan directory: %v", err)
	}

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func scanDirectory(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return database.InsertIntoFiles(filepath.Dir(path), info.Name(), time.Now().Format(time.RFC3339), info.ModTime().Format(time.RFC3339))
	})
}
