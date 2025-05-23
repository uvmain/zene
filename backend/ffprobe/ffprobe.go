package ffprobe

import (
	"path/filepath"
	"strings"

	"zene/types"
)

func GetTags(audiofilePath string) (types.TrackMetadata, error) {
	var err error
	var result types.TrackMetadata

	if filepath.Ext(audiofilePath) == ".opus" {
		result, err = GetOpusTags(audiofilePath)
	} else {
		result, err = GetCommonTags(audiofilePath)
	}

	return result, err
}

func getTagStringValue(tags map[string]string, inputs []string) string {
	for _, input := range inputs {
		value := tags[input]
		if value != "" {
			return value
		}
		value = tags[strings.ToUpper(input)]
		if value != "" {
			return value
		}
		value = tags[strings.ToLower(input)]
		if value != "" {
			return value
		}
	}
	return ""
}
