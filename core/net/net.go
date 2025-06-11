package net

import (
	"net/http"
	"time"
)

func IfModifiedResponse(w http.ResponseWriter, r *http.Request, lastModified time.Time) bool {
	w.Header().Set("Last-Modified", lastModified.Truncate(time.Second).UTC().Format(http.TimeFormat))
	w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
	ifModifiedSince := r.Header.Get("If-Modified-Since")
	if ifModifiedSince != "" {
		ifTime, err := time.Parse(http.TimeFormat, ifModifiedSince)
		if err == nil && !lastModified.Truncate(time.Second).After(ifTime) {
			w.WriteHeader(http.StatusNotModified)
			return true
		}
	}
	return false
}
