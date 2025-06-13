package main

import (
	"context"
	"zene/core/config"
	"zene/core/database"
	"zene/core/io"
	"zene/core/scanner"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config.LoadConfig()
	io.CreateDirs()
	database.Initialise(ctx)
	defer database.CloseDatabase()

	go func() {
		scanner.RunScan(ctx)
	}()

	StartServer()
}
