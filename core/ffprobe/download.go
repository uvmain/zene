package ffprobe

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

func setLatestFfprobeDownloadUrl() error {
	switch runtime.GOOS {

	case "darwin":
		targetUrl = "https://evermeet.cx/ffmpeg/getrelease/ffprobe/zip"
		fileName = "ffprobe.zip"
		ext = ".zip"
		return nil

	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			targetUrl = "https://johnvansickle.com/ffmpeg/builds/ffmpeg-git-amd64-static.tar.xz"
			ext = ".tar.xz"
			fileName = "ffprobe.tar.xz"
			return nil
		case "arm64":
			targetUrl = "https://johnvansickle.com/ffmpeg/builds/ffmpeg-git-arm64-static.tar.xz"
			ext = ".tar.xz"
			fileName = "ffprobe.tar.xz"
			return nil
		}

	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			targetUrl = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip"
			ext = ".zip"
			fileName = "ffprobe.zip"
			return nil
		case "arm64":
			targetUrl = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-winarm64-gpl.zip"
			ext = ".zip"
			fileName = "ffprobe.zip"
			return nil
		}
	}

	return fmt.Errorf("unsupported platform %s/%s", runtime.GOOS, runtime.GOARCH)
}

func DownloadFfprobeBinary() error {
	if err := setLatestFfprobeDownloadUrl(); err != nil {
		return err
	}

	if err := downloadFfprobe(); err != nil {
		return err
	}

	return nil
}

func downloadFfprobe() error {
	targetPath := filepath.Join(config.TempDirectory, fileName)
	// delete file if it already exists
	if io.FileExists(targetPath) {
		io.Cleanup(targetPath)
	}
	err := net.DownloadBinaryFile(targetUrl, targetPath)
	if err != nil {
		return fmt.Errorf("downloading ffprobe binary: %v", err)
	}
	if !io.FileExists(targetPath) {
		return fmt.Errorf("ffprobe binary file not found after download: %s", targetPath)
	}
	logger.Printf("ffprobe download from %s", targetUrl)

	filters := []string{"ffprobe", "ffprobe.exe"}

	switch ext {
	case ".zip":
		io.Unzip(targetPath, config.LibraryDirectory, filters)
	case ".tar.xz":
		io.UnTarXz(targetPath, config.LibraryDirectory, filters)
	}

	io.Cleanup(targetPath)
	return nil
}
