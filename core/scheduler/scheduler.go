package scheduler

import (
	"context"
	"time"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/scanner"
)

func Initialise(ctx context.Context) {
	startAudioCacheCleanupRoutine(ctx)
	startNowPlayingCleanupRoutine(ctx)
	startAlbumArtCleanupRoutine(ctx)
	startOrphanedPlaylistEntriesCleanupRoutine(ctx)
	startPodcastCleanupRoutine(ctx)
	startPodcastEpisodeRefreshRoutine(ctx)
	startScanScheduleRoutine(ctx)
}

func startAudioCacheCleanupRoutine(ctx context.Context) {
	logger.Println("Scheduler: starting audio cache cleanup routine")
	cleanupAudioCache(ctx)
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
	database.CleanupNowPlaying(ctx)
	go func() {
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
	database.RemoveOrphanedPlaylistEntries(ctx)
	go func() {
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
	cleanupMissingPodcasts(ctx)
	go func() {
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

func startPodcastEpisodeRefreshRoutine(ctx context.Context) {
	logger.Println("Scheduler: starting podcast episode refresh routine")
	go func() {
		fetchNewPodcastEpisodes(ctx)
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Println("Scheduler: stopping podcast episode refresh routine")
				return
			case <-ticker.C:
				fetchNewPodcastEpisodes(ctx)
			}
		}
	}()
}

func startAlbumArtCleanupRoutine(ctx context.Context) {
	logger.Println("Scheduler: starting album art cleanup routine")
	cleanupAlbumArt(ctx)
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Println("Scheduler: stopping album art cleanup routine")
				return
			case <-ticker.C:
				cleanupAlbumArt(ctx)
			}
		}
	}()
}

func startScanScheduleRoutine(ctx context.Context) {
	logger.Println("Scheduler: starting scan schedule routine")
	scanner.RunScan(ctx)
	go func() {
		ticker := time.NewTicker(45 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Println("Scheduler: stopping album art cleanup routine")
				return
			case <-ticker.C:
				scanner.RunScan(ctx)
			}
		}
	}()
}
