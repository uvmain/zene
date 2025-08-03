package ffmpeg

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"zene/core/config"
	"zene/core/io"
	"zene/core/logger"
)

func InitializeFfmpeg(ctx context.Context) error {
	logger.Printf("FFMPEG_PATH: %s", config.FfmpegPath)

	if io.FileExists(config.FfmpegPath) {
		logger.Printf("ffmpeg binary found at %s", config.FfmpegPath)
	} else {
		err := downloadFfmpegBinary()
		if err != nil {
			return fmt.Errorf("downloading ffmpeg binary: %v", err)
		}
	}

	version, err := exec.CommandContext(ctx, config.FfmpegPath, "-version").Output()
	if err != nil {
		return fmt.Errorf("ffmpeg not found at %s: %v", config.FfmpegPath, err)
	} else {
		logger.Printf("ffmpeg version is %v", strings.Split(string(version), "\n")[0])
		return nil
	}
}

func GetCoverArtFromTrack(ctx context.Context, audiofilePath string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, config.FfmpegPath,
		"-i", audiofilePath,
		"-f", "image2",
		"-vcodec", "copy",
		"-an",
		"pipe:1",
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg failed in GetCoverArtFromTrack for %s: %v\nOutput: %s", audiofilePath, err, out.String())
	}

	return out.Bytes(), nil
}
