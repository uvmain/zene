package types

type ScanStatus struct {
	Scanning      bool   `xml:"scanning,attr" json:"scanning"`
	Count         int    `xml:"count,attr" json:"count"`
	FolderCount   int    `xml:"folderCount,attr" json:"folderCount"`
	StartedDate   string `xml:"startedDate,attr" json:"startedDate"`
	Type          string `xml:"type,attr" json:"type"`
	CompletedDate string `xml:"completedDate,attr" json:"completedDate"`
}
