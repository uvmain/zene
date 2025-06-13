package scheduler

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
	"zene/core/config"
)

func maxCacheSizeBytes() int64 {
	return int64(config.AudioCacheMaxMB) * 1024 * 1024
}

func cleanupAudioCache() {
	dir := config.AudioCacheFolder
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

	log.Printf("Audio cache is %d bytes, cleaning up...", totalSize)

	// sort oldest first
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

		if totalSize <= maxCacheSizeBytes() {
			break
		}
	}
}
