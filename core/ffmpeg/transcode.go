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

func cleanupIncompleteCache(cachePath string, cacheKey string) {
	if err := os.Remove(cachePath); err != nil {
		logger.Printf("Failed to remove incomplete cache file %s: %v", cachePath, err)
	} else {
		logger.Printf("Removed incomplete cache file: %s", cachePath)
	}

	if err := database.DeleteAudioCacheEntry(cacheKey); err != nil {
		logger.Printf("Failed to delete audio cache entry for %s: %v", cacheKey, err)
	} else {
		logger.Printf("Deleted audio cache entry for %s", cacheKey)
	}
}

func TranscodeAndStream(ctx context.Context, w http.ResponseWriter, r *http.Request, filePathAbs string, trackId string, maxBitRate int, timeOffset int, format string) error {
	cacheKey := fmt.Sprintf("%s-%d.%s", trackId, maxBitRate, format)
	cachePath := filepath.Join(config.AudioCacheFolder, cacheKey)

	useCache := timeOffset <= 0

	if useCache {
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
			fileInfo, err := f.Stat()
			if err != nil {
				return fmt.Errorf("getting file info: %w", err)
			}
			w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
			w.Header().Set("Cache-Control", "public, max-age=31536000")
			w.Header().Set("Content-Type", fmt.Sprintf("audio/%s", format))
			_, err = io.Copy(w, f)
			return err
		}
	}

	if timeOffset > 0 {
		logger.Printf("Transcoding %s to stream at %s %dk starting from %ds (no cache)", filePathAbs, format, maxBitRate, timeOffset)
	} else {
		logger.Printf("Transcoding %s to stream at %s %dk", filePathAbs, format, maxBitRate)
	}

	// Build ffmpeg arguments based on format
	args := []string{"-loglevel", "error"}
	if timeOffset > 0 {
		// Place -ss before -i for faster seek
		args = append(args, "-ss", fmt.Sprintf("%d", timeOffset))
	}
	args = append(args, "-i", filePathAbs, "-vn")

	var codec, muxer, contentType string
	switch format {
	case "aac":
		codec = "aac"
		muxer = "adts"
		contentType = "audio/aac"
	case "mp3":
		codec = "libmp3lame"
		muxer = "mp3"
		contentType = "audio/mpeg"
	case "opus":
		codec = "libopus"
		muxer = "opus"
		contentType = "audio/opus"
	case "flac":
		codec = "flac"
		muxer = "flac"
		contentType = "audio/flac"
	case "vorbis":
		codec = "libvorbis"
		muxer = "ogg"
		contentType = "audio/vorbis"
	case "wav":
		codec = "pcm_s16le"
		muxer = "wav"
		contentType = "audio/wav"
	case "alac":
		codec = "alac"
		muxer = "mp4"
		contentType = "audio/alac"
	case "wma":
		codec = "wmav2"
		muxer = "asf"
		contentType = "audio/x-ms-wma"
	case "aac_latm":
		codec = "aac"
		muxer = "latm"
		contentType = "audio/aac"
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	args = append(args, "-c:a", codec, "-b:a", fmt.Sprintf("%dk", maxBitRate), "-f", muxer, "pipe:1")

	cmd := exec.CommandContext(ctx, config.FfmpegPath, args...)

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

	w.Header().Set("Content-Type", contentType)

	var mw io.Writer = w
	var cacheFile *os.File
	if useCache {
		cacheFile, err = os.Create(cachePath)
		if err != nil {
			return fmt.Errorf("creating cache file: %w", err)
		}
		defer cacheFile.Close()
		mw = io.MultiWriter(w, cacheFile)
	}

	_, err = io.Copy(mw, stdout)
	waitErr := cmd.Wait()

	if err != nil {
		if useCache {
			cacheFile.Close()
			cleanupIncompleteCache(cachePath, cacheKey)
		}
		return fmt.Errorf("copy failed: %w", err)
	}

	if waitErr != nil {
		if useCache {
			cacheFile.Close()
			cleanupIncompleteCache(cachePath, cacheKey)
		}
		return fmt.Errorf("ffmpeg exited with error: %w", waitErr)
	}

	if useCache {
		logger.Printf("Transcoding %s complete, cached at %s", filePathAbs, cachePath)
		if err = database.UpsertAudioCacheEntry(ctx, cacheKey); err != nil {
			cacheFile.Close()
			cleanupIncompleteCache(cachePath, cacheKey)
			logger.Printf("Error upserting audiocache entry for: %s", cacheKey)
			return err
		}
		cacheFile.Close()
	} else {
		logger.Printf("Transcoding %s complete (offset stream, not cached)", filePathAbs)
	}

	return nil
}
