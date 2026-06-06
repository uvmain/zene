package ffprobe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"zene/core/config"
	"zene/core/net"
)

const (
	fileName = "ffprobe.zip"
	mainUrl  = "https://ffbinaries.com/api/v1/version/latest"
)

var (
	target   string
	platform = runtime.GOOS
	arch     = runtime.GOARCH
)

func DownloadFfprobeBinary() error {
	if err := getArch(); err != nil {
		return err
	}

	if err := downloadFfBinariesFile(); err != nil {
		return err
	}
	return nil
}

func getArch() error {
	switch platform {
	case "windows":
		if arch == "amd64" {
			target = "windows-64"
		} else {
			target = "windows-32"
		}
	case "darwin":
		target = "osx-64"
	case "linux":
		if arch == "amd64" {
			target = "linux-64"
		} else {
			target = "linux-32"
		}
	default:
		return fmt.Errorf("unsupported platform/architecture for ffprobe: %s/%s", platform, arch)
	}
	return nil
}

func downloadFfBinariesFile() error {
	response, err := http.Get(mainUrl)
	if err != nil {
		return fmt.Errorf("downloading ffprobe from %s: %v", mainUrl, err)
	}
	defer response.Body.Close()

	type BinInfo struct {
		Bin map[string]struct {
			FFProbe string `json:"ffprobe"`
		} `json:"bin"`
	}

	var info BinInfo
	if err := json.NewDecoder(response.Body).Decode(&info); err != nil {
		return fmt.Errorf("decoding JSON response from %s: %v", mainUrl, err)
	}

	url := info.Bin[target].FFProbe
	if url == "" {
		return fmt.Errorf("ffprobe download url not found at %s", mainUrl)
	}

	return net.DownloadZip(url, fileName, config.LibraryDirectory, "ffprobe")
}
