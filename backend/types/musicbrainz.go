package types

type MbCoverArtResponse struct {
	Images  []MbImage `json:"images"`
	Release string    `json:"release"`
}

type MbImage struct {
	Approved   bool              `json:"approved"`
	Back       bool              `json:"back"`
	Comment    string            `json:"comment"`
	Edit       int64             `json:"edit"`
	Front      bool              `json:"front"`
	ID         interface{}       `json:"id"`
	Image      string            `json:"image"`
	Thumbnails map[string]string `json:"thumbnails"`
	Types      []string          `json:"types"`
}

type MbRelease struct {
	StatusID           string `json:"status-id"`
	Barcode            string `json:"barcode"`
	Title              string `json:"title"`
	PackagingID        string `json:"packaging-id"`
	Date               string `json:"date"`
	ID                 string `json:"id"`
	TextRepresentation struct {
		Script   string `json:"script"`
		Language string `json:"language"`
	} `json:"text-representation"`
	ReleaseEvents []struct {
		Date string `json:"date"`
		Area struct {
			Disambiguation string   `json:"disambiguation"`
			Name           string   `json:"name"`
			SortName       string   `json:"sort-name"`
			TypeID         *string  `json:"type-id"`
			ID             string   `json:"id"`
			Type           *string  `json:"type"`
			ISO31661Codes  []string `json:"iso-3166-1-codes"`
		} `json:"area"`
	} `json:"release-events"`
	CoverArtArchive struct {
		Back     bool `json:"back"`
		Front    bool `json:"front"`
		Darkened bool `json:"darkened"`
		Count    int  `json:"count"`
		Artwork  bool `json:"artwork"`
	} `json:"cover-art-archive"`
	Status         string  `json:"status"`
	ASIN           *string `json:"asin"`
	Country        string  `json:"country"`
	Packaging      string  `json:"packaging"`
	Disambiguation string  `json:"disambiguation"`
	Quality        string  `json:"quality"`
}
