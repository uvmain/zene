package ffmpeg

import (
	"fmt"
	"path/filepath"
	"runtime"
	"zene/core/config"
	"zene/core/io"
	"zene/core/logger"
	"zene/core/net"
)

var (
	targetUrl string
	platform  = runtime.GOOS
	arch      = runtime.GOARCH
	ext       string
	fileName  string
)

func setLatestFfmpegDownloadUrl() error {
	switch runtime.GOOS {

	case "darwin":
		targetUrl = "https://evermeet.cx/ffmpeg/getrelease/ffmpeg/zip"
		fileName = "ffmpeg.zip"
		ext = ".zip"
		return nil

	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			targetUrl = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-linux64-gpl.tar.xz"
			ext = ".tar.xz"
			fileName = "ffmpeg.tar.xz"
			return nil
		case "arm64":
			targetUrl = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-linuxarm64-gpl.tar.xz"
			ext = ".tar.xz"
			fileName = "ffmpeg.tar.xz"
			return nil
		}

	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			targetUrl = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip"
			ext = ".zip"
			fileName = "ffmpeg.zip"
			return nil
		case "arm64":
			targetUrl = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-winarm64-gpl.zip"
			ext = ".zip"
			fileName = "ffmpeg.zip"
			return nil
		}
	}

	return fmt.Errorf("unsupported platform %s/%s", runtime.GOOS, runtime.GOARCH)
}

func DownloadFfmpegBinary() error {
	if err := setLatestFfmpegDownloadUrl(); err != nil {
		return err
	}

	if err := downloadFfmpeg(); err != nil {
		return err
	}

	return nil
}

func downloadFfmpeg() error {
	targetPath := filepath.Join(config.TempDirectory, fileName)
	// delete file if it already exists
	if io.FileExists(targetPath) {
		io.Cleanup(targetPath)
	}
	err := net.DownloadBinaryFile(targetUrl, targetPath)
	if err != nil {
		return fmt.Errorf("downloading ffmpeg binary: %v", err)
	}
	if !io.FileExists(targetPath) {
		return fmt.Errorf("ffmpeg binary file not found after download: %s", targetPath)
	}
	logger.Printf("ffmpeg download from %s", targetUrl)

	filters := []string{"ffmpeg", "ffmpeg.exe"}

	switch ext {
	case ".zip":
		io.Unzip(targetPath, config.LibraryDirectory, filters)
	case ".tar.xz":
		io.UnTarXz(targetPath, config.LibraryDirectory, filters)
	}

	io.Cleanup(targetPath)
	return nil
}
