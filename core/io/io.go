package io

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"zene/core/logger"
	"zene/core/types"

	"github.com/djherbis/times"
)

func FileExists(absoluteFilePath string) bool {
	cleanPath, err := PathWithoutTraversal(absoluteFilePath)
	if err != nil {
		logger.Printf("error checking file existence - path traversal detected: %s", absoluteFilePath)
		return false
	}
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDir(directoryPath string) {
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		err := os.MkdirAll(directoryPath, 0755)
		if err != nil {
			logger.Printf("Error creating directory %s: %s", directoryPath, err)
		} else {
			logger.Printf("Directory created: %s", directoryPath)
		}
	} else {
		logger.Printf("Directory already exists: %s", directoryPath)
	}
}

func GetChangedTime(path string) (time.Time, error) {
	// check for invalid Windows filename characters
	if strings.ContainsAny(path, "?*<>|\"") && os.PathSeparator == '\\' {
		// fall back to os.Stat for paths with invalid characters
		info, err := os.Stat(path)
		if err != nil {
			return time.Time{}, fmt.Errorf("retrieving file times for %s: %v", path, err)
		}
		return info.ModTime(), nil
	}

	t, err := times.Stat(path)
	if err != nil {
		// if times.Stat fails, fall back to os.Stat
		info, err := os.Stat(path)
		if err != nil {
			return time.Time{}, fmt.Errorf("retrieving file times for %s: %v", path, err)
		}
		return info.ModTime(), nil
	}

	modTime := t.ModTime()
	var changeTime time.Time

	if t.HasChangeTime() {
		changeTime = t.ChangeTime()
	} else {
		changeTime = t.ModTime()
	}

	if changeTime.After(modTime) {
		modTime = changeTime
	}
	return modTime, nil
}

func GetFileBlob(ctx context.Context, filePath string) ([]byte, error) {
	filePathAbs, _ := filepath.Abs(filePath)

	if _, err := os.Stat(filePathAbs); os.IsNotExist(err) {
		logger.Printf("File does not exist: %s:  %s", filePathAbs, err)
		return nil, err
	}
	blob, err := os.ReadFile(filePathAbs)
	if err != nil {
		logger.Printf("Error reading File for filepath %s: %s", filePathAbs, err)
		return nil, err
	}

	return blob, nil
}

func GetFiles(ctx context.Context, directoryPath string, extensions []string) ([]types.File, error) {
	files := []types.File{}
	scanError := filepath.WalkDir(directoryPath, func(path string, d os.DirEntry, err error) error {

		if err != nil {
			logger.Printf("Error scanning file path %s: %v", path, err)
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if d.IsDir() {
			return nil
		}

		if len(extensions) > 0 && !slices.Contains(extensions, filepath.Ext(path)) {
			return nil
		}

		var foundFile types.File
		foundFile.FilePath = path
		foundFile.FilePathAbs, err = filepath.Abs(path)
		if err != nil {
			logger.Printf("Error getting absolute file path for %s: %v", path, err)
			return err
		}
		modTime, err := GetChangedTime(path)
		if err != nil {
			logger.Printf("Error getting changed time for %s: %v", path, err)
			return err
		}
		foundFile.DateModified = modTime.Format(time.RFC3339Nano)
		files = append(files, foundFile)
		return nil
	})

	if scanError != nil {
		logger.Printf("Error scanning files: %v", scanError)
	}

	return files, nil
}

func DeleteFile(filePath string) error {
	cleanPath, err := PathWithoutTraversal(filePath)
	if err != nil {
		return fmt.Errorf("deleting file - path traversal detected: %s", filePath)
	}
	filePathAbs, _ := filepath.Abs(cleanPath)

	if _, err := os.Stat(filePathAbs); os.IsNotExist(err) {
		return fmt.Errorf("deleting file - file does not exist: %s:  %s", filePathAbs, err)
	} else {
		err := os.Remove(filePathAbs)
		if err != nil {
			return fmt.Errorf("deleting file %s: %s", filePathAbs, err)
		}
	}
	return nil
}

func Cleanup(fileName string) error {
	macosxDir := "__MACOSX"

	err := os.Remove(fileName)
	if err != nil {
		return fmt.Errorf("removing file %s: %v", fileName, err)
	}

	err = os.RemoveAll(macosxDir)
	if err != nil {
		return fmt.Errorf("removing directory %s: %v", macosxDir, err)
	}

	return nil
}

func Unzip(srcFile string, targetDirectory string, fileNameFilter string) error {
	logger.Printf("Unzipping %s to %s", srcFile, targetDirectory)
	zipReader, err := zip.OpenReader(srcFile)
	if err != nil {
		return err
	}

	absTargetDir, err := filepath.Abs(targetDirectory)
	if err != nil {
		return fmt.Errorf("resolving target directory: %v", err)
	}

	for _, file := range zipReader.File {
		if strings.Contains(file.Name, fileNameFilter) {

			targetPath := filepath.Join(absTargetDir, filepath.Clean(file.Name))
			absTargetPath, err := filepath.Abs(targetPath)
			if err != nil {
				return fmt.Errorf("resolving target file path: %v", err)
			}

			// Ensure the target path is within the target directory
			if !strings.HasPrefix(absTargetPath, absTargetDir+string(filepath.Separator)) && absTargetPath != absTargetDir {
				return fmt.Errorf("zip slip detected: %s is outside target directory", file.Name)
			}

			fileReadCloser, err := file.Open()
			if err != nil {
				return err
			}

			// Create any necessary subdirectories
			if err := os.MkdirAll(filepath.Dir(absTargetPath), 0755); err != nil {
				fileReadCloser.Close()
				return fmt.Errorf("creating directory for %s: %v", file.Name, err)
			}

			outFile, err := os.Create(absTargetPath)
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

			logger.Printf("extracted %s to %s", file.Name, absTargetPath)
			if err := os.Chmod(absTargetPath, 0755); err != nil {
				return fmt.Errorf("setting executable permissions on %s: %v", absTargetPath, err)
			}
		}
	}

	zipReader.Close()
	Cleanup(srcFile)
	return nil
}

func PathWithoutTraversal(inputPath string) (string, error) {
	if strings.Contains(inputPath, "..") {
		return "", fmt.Errorf("path traversal detected")
	}
	cleanPath := filepath.Clean(inputPath)
	return cleanPath, nil
}
