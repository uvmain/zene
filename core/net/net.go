package net

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"zene/core/config"
	zene_io "zene/core/io"
	"zene/core/logger"
	"zene/core/subsonic"
	"zene/core/types"
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

func AddUserAgentHeaderToRequest(req *http.Request) {
	var userAgent = "zene/core/1.0 (https://github.com/uvmain/zene)"
	req.Header.Set("User-Agent", userAgent)
}

func DownloadZip(url string, fileName string, targetDirectory string, fileNameFilter string) error {
	logger.Println("Downloading:", url)
	response, err := http.Get(url)
	if err != nil {
		zene_io.Cleanup(fileName)
		return fmt.Errorf("downloading zip from %s: %v", url, err)
	}
	defer response.Body.Close()

	fileName = filepath.Join(config.TempDirectory, fileName)
	out, err := os.Create(fileName)
	if err != nil {
		out.Close()
		zene_io.Cleanup(fileName)
		return err
	}

	_, err = io.Copy(out, response.Body)
	if err != nil {
		out.Close()
		zene_io.Cleanup(fileName)
		return err
	}

	out.Close()

	if err := zene_io.Unzip(fileName, targetDirectory, fileNameFilter); err != nil {
		zene_io.Cleanup(fileName)
		return fmt.Errorf("unzipping %s: %v", fileName, err)
	}

	return nil
}

// WriteSubsonicError writes a Subsonic API error response in XML or JSON format, defaulting to XML.
// It always returns HTTP status 200 OK, as per Subsonic API specification.
// The response includes the error code and message if there is an error.
func WriteSubsonicError(w http.ResponseWriter, r *http.Request, code int, message string, helpUrl string) {

	response := subsonic.GetPopulatedSubsonicResponse(r.Context(), true)
	response.SubsonicResponse.Error.Code = code
	response.SubsonicResponse.Error.Message = message
	if helpUrl != "" {
		response.SubsonicResponse.Error.HelpUrl = helpUrl
	}

	format := r.FormValue("f")
	if format == "json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))
		xml.NewEncoder(w).Encode(response)
	}
}

func ParseBooleanFormValue(w http.ResponseWriter, r *http.Request, key string) bool {
	formValue := r.FormValue(key)
	if formValue != "" {
		parsedBool, err := strconv.ParseBool(formValue)
		if err != nil {
			errString := fmt.Sprintf("%s must be true or false", key)
			WriteSubsonicError(w, r, types.ErrorMissingParameter, errString, "")
			return false
		}
		return parsedBool
	} else {
		return false
	}
}

func ParseDuplicateFormKeys(r *http.Request, key string, intArray bool) ([]int, []string, error) { // returns []int and []string, parses []int only if intArray is true
	if err := r.ParseForm(); err != nil {
		logger.Printf("Error parsing form: %v", err)
		return nil, nil, fmt.Errorf("error parsing form: %w", err)
	}

	intSlice := []int{}
	stringSlice := r.Form[key]

	if intArray {
		for _, idStr := range stringSlice {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				logger.Printf("Error parsing %s in parseDuplicateFormKeys: %v", key, err)
				return intSlice, []string{}, fmt.Errorf("error parsing %s: %w", key, err)
			}
			intSlice = append(intSlice, id)
		}
	}
	return intSlice, stringSlice, nil
}

func GetImageFromRequest(r *http.Request, key string) (image.Image, error) {
	if err := r.ParseMultipartForm(10); err != nil {
		return nil, fmt.Errorf("error parsing multipart form: %w", err)
	}

	file, _, err := r.FormFile(key)
	if err != nil {
		return nil, fmt.Errorf("error getting image from request: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}
	return img, nil
}
