package downloader

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadAndSaveAsJPG(imageURL, outputPath string) error {
	res, err := http.Get(imageURL)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status downloading image: %s", res.Status)
	}

	contentType := res.Header.Get("Content-Type")

	// Ensure the output file has a .jpg extension
	if filepath.Ext(outputPath) != ".jpg" {
		outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".jpg"
	}

	// If the image is already a jpg, save it directly
	if contentType == "image/jpeg" || contentType == "image/jpg" {
		outFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, res.Body)
		if err != nil {
			return fmt.Errorf("failed to save jpg image: %w", err)
		}

		log.Printf("file saved to %s", outputPath)

		return nil
	}

	// Otherwise, decode and re-encode the image as jpg
	var img image.Image
	switch contentType {
	case "image/png":
		img, err = png.Decode(res.Body)
	default:
		img, _, err = image.Decode(res.Body)
	}

	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	opts := jpeg.Options{Quality: 90}
	if err := jpeg.Encode(outFile, img, &opts); err != nil {
		return fmt.Errorf("failed to encode image to jpg: %w", err)
	}

	return nil
}
