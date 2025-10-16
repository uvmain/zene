package art

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
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

func getBytesFromFilePath(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func GetImageFromInternet(imageUrl string) (image.Image, error) {
	res, err := http.Get(imageUrl)
	if err != nil {
		return nil, fmt.Errorf("downloading image: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status downloading image: %s", res.Status)
	}

	img, _, err := image.Decode(res.Body)
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}
	return img, nil
}

func resizeFileAndSaveAsJPG(imagePath string, outputPath string, pixelSize int) {
	if filepath.Ext(outputPath) != ".jpg" {
		outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".jpg"
	}

	img, err := getImageFromFilePath(imagePath)
	if err != nil {
		logger.Printf("Failed to decode image in resizeFileAndSaveAsJPG: %v", err)
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

func ResizeImageAndSaveAsJPG(img image.Image, outputPath string, pixelSize int) error {
	if filepath.Ext(outputPath) != ".jpg" {
		outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".jpg"
	}

	resizedImg := resize.Thumbnail(uint(pixelSize), uint(pixelSize), img, resize.Lanczos3)

	outFile, err := os.Create(outputPath)
	if err != nil {
		logger.Printf("Failed to create output file: %v", err)
		return err
	}
	defer outFile.Close()

	opts := jpeg.Options{Quality: 90}
	if err := jpeg.Encode(outFile, resizedImg, &opts); err != nil {
		logger.Printf("Failed to encode image to jpg: %v", err)
		return err
	}

	return nil
}

func resizeBytesAndSaveAsJPG(imgBytes []byte, outputPath string, pixelSize int) {
	if filepath.Ext(outputPath) != ".jpg" {
		outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".jpg"
	}

	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		logger.Printf("Failed to decode image in resizeBytesAndSaveAsJPG: %v", err)
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

func getImageFromFile(file multipart.File) (image.Image, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func ResizeMultipartFileAndSaveAsJPG(file multipart.File, filepath string, pixelSize int) error {
	defer file.Close()

	img, err := getImageFromFile(file)
	if err != nil {
		logger.Printf("Failed to decode image in ResizeMultipartFileAndSaveAsJPG: %v", err)
		return err
	}

	return ResizeImageAndSaveAsJPG(img, filepath, pixelSize)
}
