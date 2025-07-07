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
	startTemporaryTokensCleanupRoutine(ctx)
}

func startSessionCleanupRoutine(ctx context.Context) {
	logger.Println("Starting session cleanup routine")
	go func() {
		for {
			database.CleanupExpiredSessions(ctx)
			time.Sleep(30 * time.Minute)
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

func startTemporaryTokensCleanupRoutine(ctx context.Context) {
	logger.Println("Starting temporary_tokens cleanup routine")
	go func() {
		for {
			database.CleanupExpiredTemporaryTokens(ctx)
			time.Sleep(30 * time.Minute)
		}
	}()
}
