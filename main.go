package main

import (
	"context"
	"zene/core/config"
	"zene/core/database"
	"zene/core/io"
	"zene/core/scanner"
	"zene/core/scheduler"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config.LoadConfig()
	io.CreateDirs()
	database.Initialise(ctx)
	defer database.CloseDatabase()

	scheduler.Init(ctx)

	go func() {
		scanner.RunScan(ctx)
	}()

	StartServer()
}
