package net

import (
	"net/http"
	"time"
)

func EnableCdnCaching(w http.ResponseWriter) {
	expiryDate := time.Now().AddDate(1, 0, 0)
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("Expires", expiryDate.String())
}

func AddUserAgentHeaderToRequest(req *http.Request) {
	var userAgent = "Zene/1.0 (https://github.com/uvmain/zene)"
	req.Header.Set("User-Agent", userAgent)
}
