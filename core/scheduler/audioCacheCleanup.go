package scheduler

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
	"zene/core/config"
	"zene/core/database"
)

func maxCacheSizeBytes() int64 {
	return int64(config.AudioCacheMaxMB) * 1024 * 1024
}

func cleanupAudioCache(ctx context.Context) {
	dir := config.AudioCacheFolder

	// enforce config.AudioCacheMaxDays
	maxAge := time.Now().Add(-time.Duration(config.AudioCacheMaxDays) * 24 * time.Hour)

	staleKeys, err := database.SelectStaleAudioCacheEntries(ctx, maxAge)
	if err != nil {
		log.Printf("Error selecting stale audio cache entries: %v", err)
	}

	for _, key := range staleKeys {
		fullPath := filepath.Join(dir, key)
		err := os.Remove(fullPath)
		if err != nil {
			log.Printf("Failed to delete stale cache file %s: %v", fullPath, err)
			continue
		}
		log.Printf("Deleted stale cache file: %s", key)

		err = database.DeleteAudioCacheEntry(ctx, key)
		if err != nil {
			log.Printf("Failed to delete DB entry for %s: %v", key, err)
		}
	}

	// enforce config.AudioCacheMaxMB
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Failed to read audio cache directory: %v", err)
		return
	}

	type fileInfo struct {
		path string
		size int64
		mod  time.Time
	}

	var totalSize int64
	var infos []fileInfo

	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		fullPath := filepath.Join(dir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			log.Printf("Skipping file %s: %v", entry.Name(), err)
			continue
		}
		totalSize += info.Size()
		infos = append(infos, fileInfo{
			path: fullPath,
			size: info.Size(),
			mod:  info.ModTime(),
		})
	}

	if totalSize <= maxCacheSizeBytes() {
		return
	}

	log.Printf("Audio cache is %d bytes, cleaning up based on size...", totalSize)

	// Sort by mod time (oldest first)
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].mod.Before(infos[j].mod)
	})

	for _, fi := range infos {
		err := os.Remove(fi.path)
		if err != nil {
			log.Printf("Failed to delete cache file %s: %v", fi.path, err)
			continue
		}
		totalSize -= fi.size
		log.Printf("Deleted %s (%d bytes)", filepath.Base(fi.path), fi.size)

		cacheKey := filepath.Base(fi.path)
		err = database.DeleteAudioCacheEntry(ctx, cacheKey)
		if err != nil {
			log.Printf("Failed to delete cache row from audio_cache %s: %v", fi.path, err)
			continue
		}

		if totalSize <= maxCacheSizeBytes() {
			break
		}
	}
}
