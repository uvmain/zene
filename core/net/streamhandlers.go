package net

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"zene/core/database"
)

func HandleStreamFile(w http.ResponseWriter, r *http.Request) {
	fileId := r.PathValue("fileId")
	file, err := database.SelectFileByFileId(fileId)

	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	f, err := os.Open(file.FilePath)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		http.Error(w, "file stat error", http.StatusInternalServerError)
		return
	}

	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		// eg: bytes=0-1023
		if !strings.HasPrefix(rangeHeader, "bytes=") {
			http.Error(w, "Invalid range header format", http.StatusBadRequest)
			return
		}

		rangeParts := strings.SplitN(strings.TrimPrefix(rangeHeader, "bytes="), "-", 2)
		if len(rangeParts) != 2 {
			http.Error(w, "Invalid range format", http.StatusBadRequest)
			return
		}

		startStr := rangeParts[0]
		endStr := rangeParts[1]

		start, err := strconv.ParseInt(startStr, 10, 64)
		if err != nil && startStr != "" { // startStr can be empty for "bytes=-N"
			http.Error(w, "Invalid start range", http.StatusBadRequest)
			return
		}

		fileSize := fi.Size()
		var end int64

		if endStr == "" { // "bytes=N-"
			end = fileSize - 1
		} else {
			end, err = strconv.ParseInt(endStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid end range", http.StatusBadRequest)
				return
			}
		}

		if startStr == "" { // "bytes=-N", means last N bytes
			if end <= 0 || end > fileSize { // end here is N from "bytes=-N"
				http.Error(w, "Invalid suffix range", http.StatusBadRequest)
				return
			}
			start = fileSize - end
			end = fileSize - 1
		}

		if start < 0 || start > fileSize-1 || end < start || end >= fileSize {
			http.Error(w, "Range not satisfiable", http.StatusRequestedRangeNotSatisfiable)
			return
		}

		_, err = f.Seek(start, io.SeekStart)
		if err != nil {
			http.Error(w, "Failed to seek file", http.StatusInternalServerError)
			return
		}

		contentLength := end - start + 1
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
		w.Header().Set("Accept-Ranges", "bytes")
		contentType := mime.TypeByExtension(filepath.Ext(file.FilePath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		w.Header().Set("Content-Type", contentType)

		w.WriteHeader(http.StatusPartialContent)
		_, err = io.CopyN(w, f, contentLength)
		if err != nil {
			log.Printf("Error copying range to response: %v", err)
			return
		}
	} else {
		http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
	}
}
