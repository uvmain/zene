package types

type ScanRow = struct {
	Id           int
	ScanDate     string
	FileCount    string
	DateModified string
}

type FilesRow = struct {
	Id           int
	DirPath      string
	Filename     string
	DateAdded    string
	DateModified string
}
