package types

import "encoding/xml"

// SubsonicResponse represents the top-level response structure for Subsonic API
type SubsonicResponse struct {
	XMLName       xml.Name       `xml:"subsonic-response" json:"-"`
	Status        string         `xml:"status,attr" json:"status"`
	Version       string         `xml:"version,attr" json:"version"`
	Type          string         `xml:"type,attr" json:"type"`
	ServerVersion string         `xml:"serverVersion,attr" json:"serverVersion"`
	OpenSubsonic  bool           `xml:"openSubsonic,attr" json:"openSubsonic"`
	Error         *SubsonicError `xml:"error,omitempty" json:"error,omitempty"`
}

// SubsonicError represents a Subsonic API error
type SubsonicError struct {
	Code    int    `xml:"code,attr" json:"code"`
	Message string `xml:"message,attr" json:"message"`
	HelpUrl string `xml:"helpUrl,attr,omitempty" json:"helpUrl,omitempty"`
}

// Subsonic error codes as defined in the OpenSubsonic API specification
const (
	ErrorGeneric                   = 0
	ErrorMissingParameter          = 10
	ErrorIncompatibleVersion       = 20
	ErrorIncompatibleClient        = 30
	ErrorWrongCredentials          = 40
	ErrorTokenAuthNotSupported     = 41
	ErrorAuthMechanismNotSupported = 42
	ErrorTooManyAuthMechanisms     = 43
	ErrorInvalidApiKey             = 44
	ErrorNotAuthorized             = 50
	ErrorTrialExpired              = 60
	ErrorDataNotFound              = 70
)
