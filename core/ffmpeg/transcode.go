package ffmpeg

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"zene/core/config"
	"zene/core/database"
	"zene/core/logger"
)

func TranscodeFile(ctx context.Context, filePath string, trackId string, quality int) (string, error) {
	filePathAbs, _ := filepath.Abs(filePath)

	if _, err := os.Stat(filePathAbs); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s:  %s", filePathAbs, err)
	}

	cacheKey := fmt.Sprintf("%s-%d.aac", trackId, quality)
	cachePath := filepath.Join(config.AudioCacheFolder, cacheKey)

	if _, err := os.Stat(cachePath); err == nil {
		return cachePath, nil
	}

	logger.Printf("Transcoding %s at %dk at %s", filePath, quality, cachePath)

	cmd := exec.CommandContext(ctx, config.FfmpegPath,
		"-loglevel", "error",
		"-i", filePathAbs,
		"-vn",
		"-c:a", "aac",
		"-b:a", fmt.Sprintf("%dk", quality),
		"-f", "adts",
		cachePath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		cleanupIncompleteCache(ctx, cachePath, cacheKey)
		return "", fmt.Errorf("running ffprobe: %s", output)
	} else {
		logger.Printf("Transcoding %s complete", filePath)
	}

	err = database.UpsertAudioCacheEntry(ctx, cacheKey)
	if err != nil {
		cleanupIncompleteCache(ctx, cachePath, cacheKey)
		logger.Printf("Error upserting audiocache entry for: %s", cacheKey)
		return "", err
	}

	return cachePath, nil
}

func cleanupIncompleteCache(ctx context.Context, cachePath string, cacheKey string) {
	if err := os.Remove(cachePath); err != nil {
		logger.Printf("Failed to remove incomplete cache file %s: %v", cachePath, err)
	} else {
		logger.Printf("Removed incomplete cache file: %s", cachePath)
	}

	if err := database.DeleteAudioCacheEntry(ctx, cacheKey); err != nil {
		logger.Printf("Failed to delete audio cache entry for %s: %v", cacheKey, err)
	} else {
		logger.Printf("Deleted audio cache entry for %s", cacheKey)
	}
}

func TranscodeAndStream(ctx context.Context, w http.ResponseWriter, r *http.Request, filePathAbs string, trackId string, quality int) error {
	cacheKey := fmt.Sprintf("%s-%d.aac", trackId, quality)
	cachePath := filepath.Join(config.AudioCacheFolder, cacheKey)

	// serve from cache if it exists
	if _, err := os.Stat(cachePath); err == nil {
		logger.Printf("Serving transcoded file from cache: %s", cachePath)
		if err := database.UpsertAudioCacheEntry(ctx, cacheKey); err != nil {
			logger.Printf("Failed to update last_accessed for %s: %v", cacheKey, err)
		}
		f, err := os.Open(cachePath)
		if err != nil {
			return fmt.Errorf("opening cached file: %w", err)
		}
		defer f.Close()
		w.Header().Set("Content-Type", "audio/aac")
		_, err = io.Copy(w, f)
		return err
	}

	logger.Printf("Transcoding %s to stream at %dk", filePathAbs, quality)

	cmd := exec.CommandContext(ctx, config.FfmpegPath,
		"-loglevel", "error",
		"-i", filePathAbs,
		"-vn",
		"-c:a", "aac",
		"-b:a", fmt.Sprintf("%dk", quality),
		"-f", "adts",
		"pipe:1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("getting ffmpeg stdout: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("getting ffmpeg stderr: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting ffmpeg: %w", err)
	}

	go func() {
		slurp, _ := io.ReadAll(stderr)
		if len(slurp) > 0 {
			logger.Printf("ffmpeg stderr: %s", slurp)
		}
	}()

	w.Header().Set("Content-Type", "audio/aac")

	cacheFile, err := os.Create(cachePath)
	if err != nil {
		return fmt.Errorf("creating cache file: %w", err)
	}
	defer cacheFile.Close()

	mw := io.MultiWriter(w, cacheFile)
	_, err = io.Copy(mw, stdout)

	waitErr := cmd.Wait()

	if err != nil {
		cleanupIncompleteCache(ctx, cachePath, cacheKey)
		return fmt.Errorf("copy failed: %w", err)
	} else {
		logger.Printf("Transcoding %s complete, cached at %s", filePathAbs, cachePath)
		err = database.UpsertAudioCacheEntry(ctx, cacheKey)
		if err != nil {
			cleanupIncompleteCache(ctx, cachePath, cacheKey)
			logger.Printf("Error upserting audiocache entry for: %s", cacheKey)
			return err
		}
	}
	if waitErr != nil {
		cleanupIncompleteCache(ctx, cachePath, cacheKey)
		return fmt.Errorf("ffmpeg exited with error: %w", waitErr)
	}

	return nil
}
