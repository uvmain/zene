package ffprobe

import (
	"context"
	"path/filepath"
	"strings"

	"zene/core/types"
)

func GetTags(ctx context.Context, audiofilePath string) (types.Tags, error) {
	var err error
	var result types.Tags

	if filepath.Ext(audiofilePath) == ".opus" {
		result, err = GetOpusTags(ctx, audiofilePath)
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
