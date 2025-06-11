package logic

import (
	"context"
	"log"
	"strconv"
	"time"
)

func GenerateSlug() string {
	unixTime := time.Now().Unix()
	unixTimeString := strconv.FormatInt(unixTime, 10)

	nanoTime := time.Now().Nanosecond()
	nanoTimeString := strconv.Itoa(nanoTime)
	return unixTimeString + nanoTimeString
}

// CheckContext returns an error if the context is done/cancelled
// For example, if the http session is closed
// usage:
//
//	if err := logic.CheckContext(ctx); err != nil {
//		return []types.TrackMetadata{}, err
//	}
func CheckContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		log.Printf("Context Done: %s", ctx.Err().Error())
		return ctx.Err()
	default:
		return nil
	}
}
