package ffmpeg

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"zene/core/config"
	"zene/core/io"
	"zene/core/logger"
)

func InitializeFfmpeg(ctx context.Context) {
	logger.Printf("FFMPEG_PATH: %s", config.FfmpegPath)

	if io.FileExists(config.FfmpegPath) {
		logger.Printf("ffmpeg binary found at %s", config.FfmpegPath)
	} else {
		err := downloadFfmpegBinary()
		if err != nil {
			log.Fatalf("downloading ffmpeg binary: %v", err)
		}
	}

	version, err := exec.CommandContext(ctx, config.FfmpegPath, "-version").Output()
	if err != nil {
		log.Fatalf("ffmpeg not found at %s: %v", config.FfmpegPath, err)
	} else {
		logger.Printf("ffmpeg version is %v", strings.Split(string(version), "\n")[0])
	}
}

func GetCoverArtFromTrack(ctx context.Context, audiofilePath string) ([]byte, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, config.FfmpegPath,
		"-i", audiofilePath,
		"-an",
		"-map", "0:v:0",
		"-vframes", "1",
		"-f", "mjpeg", // force JPEG output
		"pipe:1",
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg failed for %s: %v\nstderr: %s", audiofilePath, err, stderr.String())
	}

	return stdout.Bytes(), nil
}
