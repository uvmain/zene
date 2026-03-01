package butterchurn

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"math/rand/v2"
	"path/filepath"
	"strings"
	"zene/core/types"
)

//go:embed all:presets
var presetFiles embed.FS

func getFilePaths() ([]string, error) {
	subFS, err := fs.Sub(presetFiles, "presets")
	if err != nil {
		return nil, err
	}

	var paths []string
	err = fs.WalkDir(subFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func GetPresets(count int, random bool) ([]types.ButterchurnPreset, error) {
	subFS, err := fs.Sub(presetFiles, "presets")
	if err != nil {
		return nil, err
	}

	filePaths, err := getFilePaths()
	if err != nil {
		return nil, err
	}

	if len(filePaths) == 0 {
		return nil, fmt.Errorf("no preset files found")
	}

	fileCount := len(filePaths)
	if count > fileCount {
		count = fileCount
	}
	if count == 0 {
		count = fileCount
	}

	presets := make([]types.ButterchurnPreset, 0, count)

	if random {
		rand.Shuffle(fileCount, func(i, j int) {
			filePaths[i], filePaths[j] = filePaths[j], filePaths[i]
		})
	}

	for i := 0; i < count; i++ {
		path := filePaths[i]
		content, err := fs.ReadFile(subFS, path)
		if err != nil {
			return nil, err
		}

		var jsonData interface{}
		if err := json.Unmarshal(content, &jsonData); err != nil {
			return nil, err
		}

		name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

		presets = append(presets, types.ButterchurnPreset{
			Name:   name,
			Preset: jsonData,
		})
	}

	return presets, nil
}
