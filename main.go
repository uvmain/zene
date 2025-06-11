package main

import (
	"context"
	"zene/core/art"
	"zene/core/config"
	"zene/core/database"
	"zene/core/scanner"
)

func main() {
	ctx := context.Background()

	config.LoadConfig()

	database.Initialise(ctx)
	defer database.CloseDatabase()

	art.Initialise()

	go scanner.RunScan(ctx)

	StartServer()
}
