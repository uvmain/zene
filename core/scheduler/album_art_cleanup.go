package scheduler

import (
	"context"
)

func cleanupAlbumArt(ctx context.Context) {

	// dir := config.AlbumArtFolder

	// albumIds, err := database.SelectStaleAlbumArtEntries(ctx, maxAge)
	// if err != nil {
	// 	logger.Printf("Error selecting stale audio cache entries: %v", err)
	// }

	// for _, key := range staleKeys {
	// 	fullPath := filepath.Join(dir, key)
	// 	err := os.Remove(fullPath)
	// 	if err != nil {
	// 		logger.Printf("Failed to delete stale cache file %s: %v", fullPath, err)
	// 		continue
	// 	}
	// 	logger.Printf("Deleted stale cache file: %s", key)

	// 	err = database.DeleteAudioCacheEntry(key)
	// 	if err != nil {
	// 		logger.Printf("Failed to delete DB entry for %s: %v", key, err)
	// 	}
	// }

	// // enforce config.AudioCacheMaxMB
	// files, err := os.ReadDir(dir)
	// if err != nil {
	// 	logger.Printf("Failed to read audio cache directory: %v", err)
	// 	return
	// }

	// type fileInfo struct {
	// 	path string
	// 	size int64
	// 	mod  time.Time
	// }

	// var totalSize int64
	// var infos []fileInfo

	// for _, entry := range files {
	// 	if entry.IsDir() {
	// 		continue
	// 	}
	// 	fullPath := filepath.Join(dir, entry.Name())
	// 	info, err := entry.Info()
	// 	if err != nil {
	// 		logger.Printf("Skipping file %s: %v", entry.Name(), err)
	// 		continue
	// 	}
	// 	totalSize += info.Size()
	// 	infos = append(infos, fileInfo{
	// 		path: fullPath,
	// 		size: info.Size(),
	// 		mod:  info.ModTime(),
	// 	})
	// }

	// if totalSize <= maxCacheSizeBytes() {
	// 	return
	// }

	// logger.Printf("Audio cache is %d bytes, cleaning up based on size...", totalSize)

	// // Sort by mod time (oldest first)
	// sort.Slice(infos, func(i, j int) bool {
	// 	return infos[i].mod.Before(infos[j].mod)
	// })

	// for _, fi := range infos {
	// 	err := os.Remove(fi.path)
	// 	if err != nil {
	// 		logger.Printf("Failed to delete cache file %s: %v", fi.path, err)
	// 		continue
	// 	}
	// 	totalSize -= fi.size
	// 	logger.Printf("Deleted %s (%d bytes)", filepath.Base(fi.path), fi.size)

	// 	cacheKey := filepath.Base(fi.path)
	// 	err = database.DeleteAudioCacheEntry(cacheKey)
	// 	if err != nil {
	// 		logger.Printf("Failed to delete cache row from audio_cache %s: %v", fi.path, err)
	// 		continue
	// 	}

	// 	if totalSize <= maxCacheSizeBytes() {
	// 		break
	// 	}
	// }
}
