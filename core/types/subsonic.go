package types

import (
	"encoding/xml"
)

type SubsonicStandard struct {
	XMLName       xml.Name       `xml:"subsonic-response" json:"-"`
	Xmlns         string         `xml:"xmlns,attr" json:"-"`
	Status        string         `xml:"status,attr" json:"status"`
	Version       string         `xml:"version,attr" json:"version"`
	Type          string         `xml:"type,attr" json:"type"`
	ServerVersion string         `xml:"serverVersion,attr" json:"serverVersion"`
	OpenSubsonic  bool           `xml:"openSubsonic,attr" json:"openSubsonic"`
	Error         *SubsonicError `xml:"error,omitempty" json:"error,omitempty"`
}

type SubsonicResponse struct {
	SubsonicResponse SubsonicStandard `json:"subsonic-response"`
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

func GetPopulatedSubsonicResponse(withError bool) SubsonicResponse {
	response := SubsonicResponse{
		SubsonicResponse: SubsonicStandard{
			Status:        "ok",
			Version:       "1.16.1",
			Type:          "zene",
			ServerVersion: "0.1.0",
			OpenSubsonic:  true,
			Xmlns:         "http://subsonic.org/restapi",
		},
	}

	if withError {
		response.SubsonicResponse.Status = "error"
		response.SubsonicResponse.Error = &SubsonicError{
			Code:    ErrorGeneric,
			Message: "An error occurred",
		}
	}
	return response
}

type LicenseInfo struct {
	Valid          bool   `xml:"valid,attr" json:"valid"`
	Email          string `xml:"email,attr,omitempty" json:"email,omitempty"`
	LicenseExpires string `xml:"licenseExpires,attr,omitempty" json:"licenseExpires,omitempty"`
	TrialExpires   string `xml:"trialExpires,attr,omitempty" json:"trialExpires,omitempty"`
}

type SubsonicLicense struct {
	XMLName       xml.Name       `xml:"subsonic-response" json:"-"`
	Xmlns         string         `xml:"xmlns,attr" json:"-"`
	Status        string         `xml:"status,attr" json:"status"`
	Version       string         `xml:"version,attr" json:"version"`
	Type          string         `xml:"type,attr" json:"type"`
	ServerVersion string         `xml:"serverVersion,attr" json:"serverVersion"`
	OpenSubsonic  bool           `xml:"openSubsonic,attr" json:"openSubsonic"`
	Error         *SubsonicError `xml:"error,omitempty" json:"error,omitempty"`
	License       *LicenseInfo   `xml:"license" json:"license"`
}

type SubsonicLicenseResponse struct {
	SubsonicResponse SubsonicLicense `json:"subsonic-response"`
}

type TokenInfo struct {
	Username string `xml:"username" json:"username"`
}

type SubsonicTokenInfo struct {
	XMLName       xml.Name       `xml:"subsonic-response" json:"-"`
	Xmlns         string         `xml:"xmlns,attr" json:"-"`
	Status        string         `xml:"status,attr" json:"status"`
	Version       string         `xml:"version,attr" json:"version"`
	Type          string         `xml:"type,attr" json:"type"`
	ServerVersion string         `xml:"serverVersion,attr" json:"serverVersion"`
	OpenSubsonic  bool           `xml:"openSubsonic,attr" json:"openSubsonic"`
	Error         *SubsonicError `xml:"error,omitempty" json:"error,omitempty"`
	TokenInfo     *TokenInfo     `xml:"tokenInfo" json:"tokenInfo"`
}

type SubsonicTokenInfoResponse struct {
	SubsonicResponse SubsonicTokenInfo `json:"subsonic-response"`
}

type OpenSubsonicExtensions struct {
	Name     string `xml:"name" json:"name"`
	Versions []int  `xml:"versions" json:"versions"`
}

type SubsonicOpenSubsonicExtensions struct {
	XMLName                xml.Name                  `xml:"subsonic-response" json:"-"`
	Xmlns                  string                    `xml:"xmlns,attr" json:"-"`
	Status                 string                    `xml:"status,attr" json:"status"`
	Version                string                    `xml:"version,attr" json:"version"`
	Type                   string                    `xml:"type,attr" json:"type"`
	ServerVersion          string                    `xml:"serverVersion,attr" json:"serverVersion"`
	OpenSubsonic           bool                      `xml:"openSubsonic,attr" json:"openSubsonic"`
	Error                  *SubsonicError            `xml:"error,omitempty" json:"error,omitempty"`
	OpenSubsonicExtensions []*OpenSubsonicExtensions `xml:"openSubsonicExtensions" json:"openSubsonicExtensions"`
}

type SubsonicOpenSubsonicExtensionsResponse struct {
	SubsonicResponse SubsonicOpenSubsonicExtensions `json:"subsonic-response"`
}
