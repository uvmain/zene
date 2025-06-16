package scheduler

import (
	"context"
	"time"
	"zene/core/database"
	"zene/core/logger"
)

func Initialise(ctx context.Context) {
	startSessionCleanupRoutine(ctx)
	startAudioCacheCleanupRoutine(ctx)
}

func startSessionCleanupRoutine(ctx context.Context) {
	logger.Println("Starting session cleanup routine")
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			database.CleanupExpiredSessions(ctx)
		}
	}()
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
