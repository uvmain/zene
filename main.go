package main

import (
	"context"
	"zene/core/config"
	"zene/core/database"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/scanner"
	"zene/core/scheduler"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Initialise()
	defer logger.Shutdown()

	config.LoadConfig()
	io.CreateDirs()

	database.Initialise(ctx)
	defer database.CloseDatabase()

	scheduler.Initialise(ctx)

	go func() {
		scanner.RunScan(ctx)
	}()

	StartServer()
}
