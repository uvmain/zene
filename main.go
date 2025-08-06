package main

import (
	"context"
	"zene/core/config"
	"zene/core/database"
	"zene/core/encryption"
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
	encryption.GetEncryptionKey()

	database.Initialise(ctx)
	defer database.CloseDatabase()

	ffprobe.InitializeFfprobe(ctx)
	ffmpeg.InitializeFfmpeg(ctx)

	scheduler.Initialise(ctx)

	go func() {
		scanner.RunScan(ctx)
	}()

	StartServer()
}
