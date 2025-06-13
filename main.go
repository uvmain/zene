package main

import (
	"context"
	"zene/core/config"
	"zene/core/database"
	"zene/core/io"
	// "zene/core/scanner"
)

func main() {
	ctx := context.Background()

	config.LoadConfig()

	io.CreateDirs()

	database.Initialise(ctx)
	defer database.CloseDatabase()

	// go scanner.RunScan(ctx)

	StartServer()
}
