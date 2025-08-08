package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	// watch for OS signals to gracefully shut down the database
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %s, shutting down...", sig)
		cancel()
		database.CleanShutdown()
	}()

	logger.Initialise()
	defer logger.Shutdown()

	config.LoadConfig()
	io.CreateDirs()
	encryption.GetEncryptionKey()

	database.Initialise(ctx)

	ffprobe.InitializeFfprobe(ctx)
	ffmpeg.InitializeFfmpeg(ctx)

	scheduler.Initialise(ctx)

	go func() {
		scanner.RunScan(ctx)
	}()

	StartServer()
}
