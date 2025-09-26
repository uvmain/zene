package scheduler

import (
	"context"
	"time"
	"zene/core/database"
	"zene/core/logger"
)

func Initialise(ctx context.Context) {
	startAudioCacheCleanupRoutine(ctx)
	startNowPlayingCleanupRoutine(ctx)
	startOrphanedPlaylistEntriesCleanupRoutine(ctx)
	startPodcastCleanupRoutine(ctx)
}

func startAudioCacheCleanupRoutine(ctx context.Context) {
	logger.Println("Scheduler: starting audio cache cleanup routine")
	go func() {
		cleanupAudioCache(ctx)
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
		database.CleanupNowPlaying(ctx)
		ticker := time.NewTicker(5 * time.Minute)
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

func startOrphanedPlaylistEntriesCleanupRoutine(ctx context.Context) {
	logger.Println("Scheduler: starting orphaned playlist entries cleanup routine")
	go func() {
		database.RemoveOrphanedPlaylistEntries(ctx)
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Println("Scheduler: stopping orphaned playlist entries cleanup routine")
				return
			case <-ticker.C:
				database.RemoveOrphanedPlaylistEntries(ctx)
			}
		}
	}()
}

func startPodcastCleanupRoutine(ctx context.Context) {
	logger.Println("Scheduler: starting podcast cleanup routine")
	go func() {
		cleanupMissingPodcasts(ctx)
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Println("Scheduler: stopping podcast cleanup routine")
				return
			case <-ticker.C:
				cleanupMissingPodcasts(ctx)
			}
		}
	}()
}
