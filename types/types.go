package types

type ScanRow = struct {
	Id           string
	ScanDate     string
	FileCount    string
	DateModified string
}

type FilesRow = struct {
	Id           string
	DirPath      string
	Filename     string
	DateAdded    string
	DateModified string
}
