package database

import (
	"context"
	"zene/core/logic"
	"zene/core/types"
)

func GetIndexes(ctx context.Context, musicFolderIds []int64, ifModifiedSince int64) (types.SubsonicIndexes, error) {
	latestScan, err := GetLatestCompletedScan(ctx)
	if err != nil {
		return types.SubsonicIndexes{}, err
	}

	latestScanTime := logic.GetStringTimeFormatted(latestScan.CompletedDate)
	latestScanTimeUnix := latestScanTime.UnixMilli()

	response := types.SubsonicIndexes{}
	response.IgnoredArticles = ""
	response.LastModified = latestScanTimeUnix

	if ifModifiedSince != 0 && latestScanTimeUnix > ifModifiedSince {
		return response, nil
	}

	return response, nil
}
