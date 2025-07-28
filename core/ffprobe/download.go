package ffprobe

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"zene/core/config"
	"zene/core/net"
)

const (
	fileName = "ffprobe.zip"
	mainUrl  = "https://ffbinaries.com/api/v1/version/latest"
	macUrl   = "https://www.osxexperts.net"
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

	if platform == "darwin" {
		if err := downloadOsxExpertsBinariesFile(); err != nil {
			return err
		}
	} else {
		if err := downloadFfBinariesFile(); err != nil {
			return err
		}
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
		if arch == "arm64" {
			target = "arm"
		} else {
			target = "intel"
		}
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

func downloadOsxExpertsBinariesFile() error {
	response, err := http.Get(macUrl)
	if err != nil {
		return fmt.Errorf("downloading ffprobe from %s: %v", macUrl, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	html := string(body)
	prefix := `href="`
	suffix := `"`

	elementIndex := strings.Index(html, "ffprobe")
	if elementIndex == -1 {
		return fmt.Errorf("finding ffprobe link from html")
	}

	start := strings.LastIndex(html[:elementIndex], prefix) + len(prefix)
	end := strings.Index(html[start:], suffix)
	if end == -1 {
		return fmt.Errorf("extracting ffprobe link from html")
	}
	url := html[start : start+end]

	return net.DownloadZip(url, fileName, config.LibraryDirectory, "ffprobe")
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
		return fmt.Errorf("ffprobe not found at %s", mainUrl)
	}

	return net.DownloadZip(url, fileName, config.LibraryDirectory, "ffprobe")
}
