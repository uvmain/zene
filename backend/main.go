package main

import (
	"zene/art"
	"zene/config"
	"zene/database"
	"zene/scanner"
)

func main() {
	config.LoadConfig()

	database.Initialise()
	defer database.CloseDatabase()

	art.Initialise()

	go scanner.RunScan()

	StartServer()
}
