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
		return "", fmt.Errorf("File does not exist: %s:  %s", filePathAbs, err)
	}

	cacheKey := fmt.Sprintf("%s-%d.aac", trackId, quality)
	cachePath := filepath.Join(config.AudioCacheFolder, cacheKey)

	if _, err := os.Stat(cachePath); err == nil {
		return cachePath, nil
	}

	logger.Printf("Transcoding %s at %dk at %s", filePath, quality, cachePath)

	cmd := exec.Command("ffmpeg",
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
		return "", fmt.Errorf("Error running ffprobe: %s", output)
	} else {
		logger.Printf("Transcoding %s complete", filePath)
	}

	err = database.UpsertAudioCacheEntry(ctx, cacheKey)
	if err != nil {
		logger.Printf("Error upserting audiocache entry for: %s", cacheKey)
		return "", err
	}

	return cachePath, nil
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
			return fmt.Errorf("failed to open cached file: %w", err)
		}
		defer f.Close()
		w.Header().Set("Content-Type", "audio/aac")
		_, err = io.Copy(w, f)
		return err
	}

	logger.Printf("Transcoding %s to stream at %dk", filePathAbs, quality)

	cmd := exec.Command("ffmpeg",
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
		return fmt.Errorf("failed to get ffmpeg stdout: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get ffmpeg stderr: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
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
		return fmt.Errorf("failed to create cache file: %w", err)
	}
	defer cacheFile.Close()

	mw := io.MultiWriter(w, cacheFile)
	_, err = io.Copy(mw, stdout)

	waitErr := cmd.Wait()

	if err != nil {
		return fmt.Errorf("copy failed: %w", err)
	} else {
		logger.Printf("Transcoding %s complete, cached at %s", filePathAbs, cachePath)
		err = database.UpsertAudioCacheEntry(ctx, cacheKey)
		if err != nil {
			logger.Printf("Error upserting audiocache entry for: %s", cacheKey)
			return err
		}
	}
	if waitErr != nil {
		return fmt.Errorf("ffmpeg exited with error: %w", waitErr)
	}

	return nil
}
