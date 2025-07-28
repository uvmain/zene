package ffprobe

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"zene/core/config"
	"zene/core/logger"
)

const (
	fileName  = "ffprobe.zip"
	macosxDir = "__MACOSX"
	mainUrl   = "https://ffbinaries.com/api/v1/version/latest"
	macUrl    = "https://www.osxexperts.net"
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

	return downloadZip(url)
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
	return downloadZip(url)
}

func downloadZip(url string) error {
	logger.Println("Downloading:", url)
	response, err := http.Get(url)
	if err != nil {
		cleanup()
		return fmt.Errorf("downloading zip from %s: %v", url, err)
	}
	defer response.Body.Close()

	out, err := os.Create(fileName)
	if err != nil {
		out.Close()
		cleanup()
		return err
	}

	_, err = io.Copy(out, response.Body)
	if err != nil {
		out.Close()
		cleanup()
		return err
	}

	out.Close()

	if err := unzip(fileName); err != nil {
		cleanup()
		return fmt.Errorf("unzipping %s: %v", fileName, err)
	}

	return nil
}

func unzip(src string) error {
	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		if strings.Contains(file.Name, "ffprobe") {
			fileReadCloser, err := file.Open()
			if err != nil {
				return err
			}

			outFile, err := os.Create(file.Name)
			if err != nil {
				fileReadCloser.Close()
				return err
			}

			_, err = io.Copy(outFile, fileReadCloser)
			fileReadCloser.Close()
			outFile.Close()
			if err != nil {
				return err
			}

			// move file to LibraryDirectory
			if err := os.Rename(file.Name, filepath.Join(config.LibraryDirectory, file.Name)); err != nil {
				return fmt.Errorf("moving %s to %s: %v", file.Name, config.LibraryDirectory, err)
			}
			logger.Printf("ffprobe extracted %s to %s", file.Name, config.LibraryDirectory)
		}
	}

	zipReader.Close()
	cleanup()
	return nil
}

func cleanup() {
	err := os.Remove(fileName)
	if err != nil {
		logger.Printf("Error removing file %s: %v", fileName, err)
	}

	err = os.RemoveAll(macosxDir)
	if err != nil {
		logger.Printf("Error removing file %s: %v", fileName, err)
	}
}
