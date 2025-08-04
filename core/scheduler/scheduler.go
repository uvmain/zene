package scheduler

import (
	"context"
	"time"
	"zene/core/logger"
)

func Initialise(ctx context.Context) {
	startAudioCacheCleanupRoutine(ctx)
}

func startAudioCacheCleanupRoutine(ctx context.Context) {
	logger.Println("Starting audio cache cleanup routine")
	go func() {
		for {
			cleanupAudioCache(ctx)
			time.Sleep(1 * time.Hour)
		}
	}()
}
