package auth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestWriteSubsonicErrorXML(t *testing.T) {
	// Create a request with default format (XML)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	writeSubsonicError(w, req, ErrorWrongCredentials, "Test error message")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/xml" {
		t.Errorf("Expected Content-Type application/xml, got %s", contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>") {
		t.Error("Expected XML declaration in response")
	}
	if !strings.Contains(body, "status=\"failed\"") {
		t.Error("Expected status=\"failed\" in response")
	}
	if !strings.Contains(body, "code=\"40\"") {
		t.Error("Expected code=\"40\" in response")
	}
	if !strings.Contains(body, "message=\"Test error message\"") {
		t.Error("Expected error message in response")
	}
}

func TestWriteSubsonicErrorJSON(t *testing.T) {
	// Create a request with JSON format
	form := url.Values{}
	form.Add("f", "json")
	req := httptest.NewRequest("POST", "/test", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	writeSubsonicError(w, req, ErrorMissingParameter, "Test JSON error")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, `"status":"failed"`) {
		t.Error("Expected status:failed in JSON response")
	}
	if !strings.Contains(body, `"code":10`) {
		t.Error("Expected code:10 in JSON response")
	}
	if !strings.Contains(body, `"Test JSON error"`) {
		t.Error("Expected error message in JSON response")
	}
}

func TestSubsonicErrorCodes(t *testing.T) {
	// Test that error code constants have expected values
	expectedCodes := map[string]int{
		"ErrorGeneric":               0,
		"ErrorMissingParameter":      10,
		"ErrorIncompatibleVersion":   20,
		"ErrorIncompatibleClient":    30,
		"ErrorWrongCredentials":      40,
		"ErrorTokenAuthNotSupported": 41,
		"ErrorNotAuthorized":         50,
		"ErrorTrialExpired":          60,
		"ErrorDataNotFound":          70,
	}

	actualCodes := map[string]int{
		"ErrorGeneric":               ErrorGeneric,
		"ErrorMissingParameter":      ErrorMissingParameter,
		"ErrorIncompatibleVersion":   ErrorIncompatibleVersion,
		"ErrorIncompatibleClient":    ErrorIncompatibleClient,
		"ErrorWrongCredentials":      ErrorWrongCredentials,
		"ErrorTokenAuthNotSupported": ErrorTokenAuthNotSupported,
		"ErrorNotAuthorized":         ErrorNotAuthorized,
		"ErrorTrialExpired":          ErrorTrialExpired,
		"ErrorDataNotFound":          ErrorDataNotFound,
	}

	for name, expected := range expectedCodes {
		if actual := actualCodes[name]; actual != expected {
			t.Errorf("Expected %s to be %d, got %d", name, expected, actual)
		}
	}
}