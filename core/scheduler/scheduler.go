package scheduler

import (
	"context"
	"time"
	"zene/core/database"
	"zene/core/deezer"
	"zene/core/logger"
)

func Initialise(ctx context.Context) {
	startAudioCacheCleanupRoutine(ctx)
	startNowPlayingCleanupRoutine(ctx)
	startDeezerCacheCleanupRoutine(ctx)
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

func startNowPlayingCleanupRoutine(ctx context.Context) {
	logger.Println("Scheduler: starting now playing cleanup routine")
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Println("Scheduler: stopping now playing cleanup routine")
				return
			case <-ticker.C:
				database.CleanupNowPlaying(ctx)
			}
		}
	}()
}

func startDeezerCacheCleanupRoutine(ctx context.Context) {
	logger.Println("Scheduler: starting deezer cache cleanup routine")
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Println("Scheduler: stopping deezer cache cleanup routine")
				return
			case <-ticker.C:
				deezer.CleanupSimilarArtistsCache(ctx)
				deezer.CleanupTopSongsCache(ctx)
			}
		}
	}()
}
