package logic

import (
	"context"
	"sync"
	"time"
	"zene/core/logger"
	"zene/core/types"
)

var (
	bootTime     time.Time
	bootTimeOnce sync.Once
)

func GetBootTime() time.Time {
	bootTimeOnce.Do(func() {
		bootTime = time.Now().UTC().Truncate(time.Second)
	})
	return bootTime
}

// CheckContext returns an error if the context is done/cancelled
// For example, if the http session is closed
// usage:
//
//	if err := logic.CheckContext(ctx); err != nil {
//		return []types.Metadata{}, err
//	}
func CheckContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		logger.Printf("Context Done: %s", ctx.Err().Error())
		return ctx.Err()
	default:
		return nil
	}
}

func FilesInSliceOnceNotInSliceTwo(slice1, slice2 []types.File) []types.File {
	slice2Map := make(map[string]bool)
	for _, f := range slice2 {
		slice2Map[f.FilePathAbs] = true
	}

	var diff []types.File
	for _, f := range slice1 {
		if !slice2Map[f.FilePathAbs] {
			diff = append(diff, f)
		}
	}

	return diff
}

func GetCurrentTimeFormatted() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}
