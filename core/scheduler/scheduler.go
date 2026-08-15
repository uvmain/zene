package scheduler

import (
	"context"
	"time"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/scanner"
)

func startSchedule(ctx context.Context, name string, interval time.Duration, task func(context.Context) error) {
	logger.Printf("[Scheduler] starting %s routine", name)
	err := task(ctx)
	if err != nil {
		logger.Printf("[Scheduler] error running %s routine: %v", name, err)
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Printf("[Scheduler] stopping %s routine", name)
				return
			case <-ticker.C:
				err := task(ctx)
				if err != nil {
					logger.Printf("[Scheduler] error running %s routine: %v", name, err)
				}
			}
		}
	}()
}

func Initialise(ctx context.Context) {
	startSchedule(ctx, "audio cache cleanup", 1*time.Hour, cleanupAudioCache)
	startSchedule(ctx, "now playing cleanup", 5*time.Minute, database.CleanupNowPlaying)
	startSchedule(ctx, "orphaned playlist entries cleanup", 1*time.Hour, database.RemoveOrphanedPlaylistEntries)
	startSchedule(ctx, "album art cleanup", 1*time.Hour, cleanupAlbumArt)
	startSchedule(ctx, "artist art cleanup", 1*time.Hour, cleanupArtistArt)
	startSchedule(ctx, "podcast cleanup", 1*time.Hour, cleanupMissingPodcasts)
	startSchedule(ctx, "podcast episode refresh", 2*time.Hour, fetchNewPodcastEpisodes)
	startSchedule(ctx, "playback reports cleanup", 5*time.Minute, database.CleanupPlaybackReports)
	startSchedule(ctx, "transcode params cleanup", 1*time.Hour, database.ClearOldTranscodeParams)
	startSchedule(ctx, "share cleanup", 1*time.Hour, database.ClearExpiredShares)
	startSchedule(ctx, "scheduled scan", 45*time.Minute, scanner.RunScheduledScan)
}
