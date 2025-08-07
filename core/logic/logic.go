package logic

import (
	"context"
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"
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
	case "admin_role":
		return false
	case "settings_role":
		return true
	case "stream_role":
		return true
	case "jukebox_role":
		return false
	case "download_role":
		return false
	case "upload_role":
		return false
	case "playlist_role":
		return false
	case "cover_art_role":
		return false
	case "comment_role":
		return false
	case "podcast_role":
		return false
	case "share_role":
		return false
	case "video_conversion_role":
		return false
	case "scrobbling_enabled":
		return true
	case "ldap_authenticated":
		return false
	default:
		return true
	}
}
