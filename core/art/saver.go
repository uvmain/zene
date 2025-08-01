package art

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"zene/core/logger"

	"github.com/nfnt/resize"
)

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var img image.Image
	if strings.ToLower(filepath.Ext(filePath)) == ".jpg" || strings.ToLower(filepath.Ext(filePath)) == ".jpeg" {
		img, _, err = image.Decode(f)
	} else if strings.ToLower(filepath.Ext(filePath)) == ".png" {
		img, err = png.Decode(f)
	}
	return img, err
}

func getImageFromInternet(imageUrl string) (image.Image, error) {
	res, err := http.Get(imageUrl)
	if err != nil {
		return nil, fmt.Errorf("Failed to download image: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status downloading image: %s", res.Status)
	}

	contentType := res.Header.Get("Content-Type")

	var img image.Image
	switch contentType {
	case "image/png":
		img, err = png.Decode(res.Body)
	default:
		img, _, err = image.Decode(res.Body)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to decode image: %w", err)
	}

	return img, nil
}

func resizeFileAndSaveAsJPG(imagePath string, outputPath string, pixelSize int) {
	if filepath.Ext(outputPath) != ".jpg" {
		outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".jpg"
	}

	img, err := getImageFromFilePath(imagePath)
	if err != nil {
		logger.Printf("Failed to decode image: %v", err)
		return
	}

	resizedImg := resize.Thumbnail(uint(pixelSize), uint(pixelSize), img, resize.Lanczos3)

	outFile, err := os.Create(outputPath)
	if err != nil {
		logger.Printf("Failed to create output file: %v", err)
		return
	}
	defer outFile.Close()

	opts := jpeg.Options{Quality: 90}
	if err := jpeg.Encode(outFile, resizedImg, &opts); err != nil {
		logger.Printf("Failed to encode image to jpg: %v", err)
		return
	}
}

func resizeImageAndSaveAsJPG(img image.Image, outputPath string, pixelSize int) {
	if filepath.Ext(outputPath) != ".jpg" {
		outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".jpg"
	}

	resizedImg := resize.Thumbnail(uint(pixelSize), uint(pixelSize), img, resize.Lanczos3)

	outFile, err := os.Create(outputPath)
	if err != nil {
		logger.Printf("Failed to create output file: %v", err)
		return
	}
	defer outFile.Close()

	opts := jpeg.Options{Quality: 90}
	if err := jpeg.Encode(outFile, resizedImg, &opts); err != nil {
		logger.Printf("Failed to encode image to jpg: %v", err)
		return
	}
}

func resizeBytesAndSaveAsJPG(imgBytes []byte, outputPath string, pixelSize int) {
	if filepath.Ext(outputPath) != ".jpg" {
		outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".jpg"
	}

	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		logger.Printf("Failed to decode image: %v", err)
		return
	}

	resizedImg := resize.Thumbnail(uint(pixelSize), uint(pixelSize), img, resize.Lanczos3)

	outFile, err := os.Create(outputPath)
	if err != nil {
		logger.Printf("Failed to create output file: %v", err)
		return
	}
	defer outFile.Close()

	opts := jpeg.Options{Quality: 90}
	if err := jpeg.Encode(outFile, resizedImg, &opts); err != nil {
		logger.Printf("Failed to encode image to jpg: %v", err)
		return
	}
}
