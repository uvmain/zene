package logic

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"
	"zene/core/config"
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

func GetStringTimeFormatted(timeString string) time.Time {
	timeTime, err := time.Parse(time.RFC3339Nano, timeString)
	if err != nil {
		logger.Printf("Error parsing time string '%s': %v", timeString, err)
		return time.Time{}
	}
	return timeTime
}

func IntSliceToInt64Slice(ints []int) []int64 {
	int64s := make([]int64, len(ints))
	for i, v := range ints {
		int64s[i] = int64(v)
	}
	return int64s
}

func GenerateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
	password := make([]byte, length)

	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}

	return string(password), nil
}

func StringToIntSlice(folderString string) []int {
	if folderString == "" {
		return nil
	}

	var folderIds []int
	for _, idStr := range strings.Split(folderString, ",") {
		if id, err := strconv.Atoi(idStr); err == nil {
			folderIds = append(folderIds, id)
		} else {
			logger.Printf("Error parsing folder ID from string '%s': %v", idStr, err)
		}
	}
	return folderIds
}

func GetDefaultRoleValue(roleName string) bool {
	switch roleName {
	case "adminRole":
		return false
	case "settingsRole":
		return true
	case "streamRole":
		return true
	case "jukeboxRole":
		return false
	case "downloadRole":
		return false
	case "uploadRole":
		return false
	case "playlistRole":
		return false
	case "coverArtRole":
		return false
	case "commentRole":
		return false
	case "podcastRole":
		return false
	case "shareRole":
		return false
	case "videoConversionRole":
		return false
	case "scrobblingEnabled":
		return true
	case "ldapAuthenticated":
		return false
	default:
		return true
	}
}

func GetArtistImageUrl(artistID string) string {
	return fmt.Sprintf("%s/share/img/%s", config.BaseUrl, artistID)
}
