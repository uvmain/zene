package scheduler

import (
	"context"
	"log"
	"time"
	"zene/core/database"
)

func Init(ctx context.Context) {
	startSessionCleanupRoutine(ctx)
	startAudioCacheCleanupRoutine(ctx)
}

func startSessionCleanupRoutine(ctx context.Context) {
	log.Println("Starting session cleanup routine")
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			database.CleanupExpiredSessions(ctx)
		}
	}()
}

func startAudioCacheCleanupRoutine(ctx context.Context) {
	log.Println("Starting audio cache cleanup routine")
	go func() {
		for {
			cleanupAudioCache(ctx)
			time.Sleep(1 * time.Hour)
		}
	}()
}
