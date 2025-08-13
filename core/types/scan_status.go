package types

type ScanStatus struct {
	Scanning    bool   `xml:"scanning,attr" json:"scanning"`
	Count       int    `xml:"count,attr" json:"count"`
	FolderCount int    `xml:"folderCount,attr" json:"folderCount"`
	LastScan    string `xml:"lastScan,attr" json:"lastScan"`
	ScanType    string `xml:"scanType,attr" json:"scanType"`
	ElapsedTime int    `xml:"elapsedTime,attr" json:"elapsedTime"`
}
