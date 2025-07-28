package ffmpeg

import (
	"fmt"
	"os/exec"
	"strings"

	"zene/core/config"
	"zene/core/io"
	"zene/core/logger"
)

func InitializeFfmpeg() error {
	logger.Printf("FFMPEG_PATH: %s", config.FfmpegPath)

	if io.FileExists(config.FfmpegPath) {
		logger.Printf("ffmpeg binary found at %s", config.FfmpegPath)
	} else {
		err := downloadFfmpegBinary()
		if err != nil {
			return fmt.Errorf("downloading ffmpeg binary: %v", err)
		}
	}

	version, err := exec.Command(config.FfmpegPath, "-version").Output()
	if err != nil {
		return fmt.Errorf("ffmpeg not found at %s: %v", config.FfmpegPath, err)
	} else {
		logger.Printf("ffmpeg version is %v", strings.Split(string(version), "\n")[0])
		return nil
	}
}
