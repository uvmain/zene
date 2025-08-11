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
	logger.Println("Scheduler: starting audio cache cleanup routine")
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Println("Scheduler: stopping audio cache cleanup routine")
				return
			case <-ticker.C:
				cleanupAudioCache(ctx)
			}
		}
	}()
}
