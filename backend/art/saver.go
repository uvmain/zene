package art

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

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

func resizeAndSaveAsJPG(imagePath string, outputPath string, pixelSize int) error {
	if filepath.Ext(outputPath) != ".jpg" {
		outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".jpg"
	}

	img, err := getImageFromFilePath(imagePath)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize the image while maintaining aspect ratio
	resizedImg := resize.Thumbnail(uint(pixelSize), uint(pixelSize), img, resize.Lanczos3)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	opts := jpeg.Options{Quality: 90}
	if err := jpeg.Encode(outFile, resizedImg, &opts); err != nil {
		return fmt.Errorf("failed to encode image to jpg: %w", err)
	}

	return nil
}
