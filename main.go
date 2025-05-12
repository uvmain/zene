package main

import (
	"zene/config"
	"zene/database"
	"zene/scanner"
)

func main() {
	config.LoadConfig()

	database.Initialise()
	defer database.CloseDatabase()

	go scanner.ScanMusicDirectory()

	StartServer()

}
