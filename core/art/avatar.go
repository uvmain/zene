package art

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"zene/core/config"
	"zene/core/logger"
)

func UpsertUserAvatarImage(userId int, img image.Image) error {
	fileName := fmt.Sprintf("%d.jpg", userId)
	filePath := filepath.Join(config.UserAvatarFolder, fileName)

	// if file already exists, delete it
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			logger.Printf("Error deleting existing avatar image: %v", err)
			return err
		}
	}

	go resizeImageAndSaveAsJPG(img, filePath, 512)

	return nil
}

func GetUserAvatarImage(userId int) ([]byte, error) {
	fileName := fmt.Sprintf("%d.jpg", userId)
	filePath := filepath.Join(config.UserAvatarFolder, fileName)

	_, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("file does not exist: %s:  %s", filePath, err)
	}

	blob, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading image for filepath %s: %s", filePath, err)
	}
	return blob, nil
}

func DeleteUserAvatarImage(userId int) error {
	fileName := fmt.Sprintf("%d.jpg", userId)
	filePath := filepath.Join(config.UserAvatarFolder, fileName)

	_, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("file does not exist: %s:  %s", filePath, err)
	}

	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("error deleting image for filepath %s: %s", filePath, err)
	}
	return nil
}
