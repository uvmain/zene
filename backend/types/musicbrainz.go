package types

type CoverArtResponse struct {
	Images  []Image `json:"images"`
	Release string  `json:"release"`
}

type Image struct {
	Approved   bool              `json:"approved"`
	Back       bool              `json:"back"`
	Comment    string            `json:"comment"`
	Edit       int64             `json:"edit"`
	Front      bool              `json:"front"`
	ID         string            `json:"id"`
	Image      string            `json:"image"`
	Thumbnails map[string]string `json:"thumbnails"`
	Types      []string          `json:"types"`
}
