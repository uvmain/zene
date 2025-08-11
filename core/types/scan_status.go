package types

import (
	"encoding/xml"
)

type ScanStatus struct {
	Scanning    bool   `xml:"scanning,attr" json:"scanning"`
	Count       int    `xml:"count,attr" json:"count"`
	FolderCount int    `xml:"folderCount,attr" json:"folderCount"`
	LastScan    string `xml:"lastScan,attr" json:"lastScan"`
	ScanType    string `xml:"scanType,attr" json:"scanType"`
	ElapsedTime int    `xml:"elapsedTime,attr" json:"elapsedTime"`
}

type SubsonicScanStatus struct {
	XMLName       xml.Name       `xml:"subsonic-response" json:"-"`
	Xmlns         string         `xml:"xmlns,attr" json:"-"`
	Status        string         `xml:"status,attr" json:"status"`
	Version       string         `xml:"version,attr" json:"version"`
	Type          string         `xml:"type,attr" json:"type"`
	ServerVersion string         `xml:"serverVersion,attr" json:"serverVersion"`
	OpenSubsonic  bool           `xml:"openSubsonic,attr" json:"openSubsonic"`
	Error         *SubsonicError `xml:"error,omitempty" json:"error,omitempty"`
	ScanStatus    *ScanStatus    `xml:"scanStatus" json:"scanStatus"`
}

type SubsonicScanStatusResponse struct {
	SubsonicResponse SubsonicScanStatus `json:"subsonic-response"`
}
