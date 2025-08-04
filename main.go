package main

import (
	"context"
	"log"
	"zene/core/auth"
	"zene/core/config"
	"zene/core/database"
	"zene/core/ffmpeg"
	"zene/core/ffprobe"
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

	err := ffprobe.InitializeFfprobe(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize ffprobe: %v", err)
		return
	}

	err = ffmpeg.InitializeFfmpeg(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize ffmpeg: %v", err)
		return
	}

	scheduler.Initialise(ctx)

	auth.Initialise(ctx)

	go func() {
		scanner.RunScan(ctx)
	}()

	StartServer()
}
