package logic

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"os"

	"golang.org/x/image/draw"
)

func ResizeJpegImage(ctx context.Context, filePath string, maxPixels int, quality int) ([]byte, error) {
	if quality < 1 || quality > 100 {
		quality = 90
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sourceFile, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := sourceFile.Bounds()
	originalWidth, originalHeight := bounds.Dx(), bounds.Dy()

	var newWidth, newHeight int
	if originalWidth >= originalHeight {
		newWidth = maxPixels
		newHeight = originalHeight * maxPixels / originalWidth
	} else {
		newHeight = maxPixels
		newWidth = originalWidth * maxPixels / originalHeight
	}

	destinationImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	resultChannel := make(chan []byte, 1)
	errChannel := make(chan error, 1)

	go func() {
		draw.CatmullRom.Scale(destinationImage, destinationImage.Bounds(), sourceFile, bounds, draw.Over, nil)

		var imageBuffer bytes.Buffer
		err := jpeg.Encode(&imageBuffer, destinationImage, &jpeg.Options{Quality: quality})
		if err != nil {
			errChannel <- err
			return
		}
		resultChannel <- imageBuffer.Bytes()
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errChannel:
		return nil, err
	case result := <-resultChannel:
		return result, nil
	}
}
