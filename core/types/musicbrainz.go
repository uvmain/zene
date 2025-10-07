package types

type MbCoverArtResponse struct {
	Images  []MbImage `json:"images"`
	Release string    `json:"release"`
}

type MbImage struct {
	Approved   bool         `json:"approved"`
	Back       bool         `json:"back"`
	Comment    string       `json:"comment"`
	Edit       int          `json:"edit"`
	Front      bool         `json:"front"`
	ID         interface{}  `json:"id"`
	Image      string       `json:"image"`
	Thumbnails MbThumbnails `json:"thumbnails"`
	Types      []string     `json:"types"`
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
	Media          []struct {
		Format   string `json:"format"`
		Position int    `json:"position"`
		Pregap   struct {
			ID        string `json:"id"`
			Title     string `json:"title"`
			Length    int    `json:"length"`
			Recording struct {
				Disambiguation   string `json:"disambiguation"`
				FirstReleaseDate string `json:"first-release-date"`
				Title            string `json:"title"`
				ID               string `json:"id"`
				Length           int    `json:"length"`
				Video            bool   `json:"video"`
			} `json:"recording"`
			Number   string `json:"number"`
			Position int    `json:"position"`
		} `json:"pregap"`
		TrackCount  int    `json:"track-count"`
		TrackOffset int    `json:"track-offset"`
		ID          string `json:"id"`
		FormatID    string `json:"format-id"`
		Title       string `json:"title"`
		Tracks      []struct {
			Position  int    `json:"position"`
			ID        string `json:"id"`
			Length    int    `json:"length"`
			Recording struct {
				Disambiguation   string `json:"disambiguation"`
				FirstReleaseDate string `json:"first-release-date"`
				ID               string `json:"id"`
				Length           int    `json:"length"`
				Title            string `json:"title"`
				Video            bool   `json:"video"`
			} `json:"recording"`
			Title  string `json:"title"`
			Number string `json:"number"`
		} `json:"tracks"`
	} `json:"media"`
}

/* cspell: disable */
type MbArtist struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	SortName       string       `json:"sort-name"`
	Type           string       `json:"type"`
	TypeID         string       `json:"type-id"`
	Country        string       `json:"country"`
	GenderID       *string      `json:"gender-id"`
	Disambiguation string       `json:"disambiguation"`
	IPIs           []string     `json:"ipis"`
	ISNIs          []string     `json:"isnis"`
	LifeSpan       MbLifeSpan   `json:"life-span"`
	Area           MbArea       `json:"area"`
	BeginArea      *MbArea      `json:"begin-area"`
	EndArea        *MbArea      `json:"end-area"`
	Relations      []MbRelation `json:"relations"`
}

/* cspell: enable */

type MbLifeSpan struct {
	Begin string  `json:"begin"`
	End   *string `json:"end"`
	Ended bool    `json:"ended"`
}

type MbArea struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	SortName       string   `json:"sort-name"`
	ISO3166_1Codes []string `json:"iso-3166-1-codes,omitempty"`
	ISO3166_2Codes []string `json:"iso-3166-2-codes,omitempty"`
	Type           *string  `json:"type"`
	TypeID         *string  `json:"type-id"`
	Disambiguation string   `json:"disambiguation"`
}

type MbRelation struct {
	Type            string        `json:"type"`
	TypeID          string        `json:"type-id"`
	Direction       string        `json:"direction"`
	Ended           bool          `json:"ended"`
	Begin           *string       `json:"begin"`
	End             *string       `json:"end"`
	URL             MbRelationURL `json:"url"`
	TargetType      string        `json:"target-type"`
	TargetCredit    string        `json:"target-credit"`
	SourceCredit    string        `json:"source-credit"`
	Attributes      []interface{} `json:"attributes"`
	AttributeIDs    interface{}   `json:"attribute-ids"`
	AttributeValues interface{}   `json:"attribute-values"`
}

type MbRelationURL struct {
	ID       string `json:"id"`
	Resource string `json:"resource"`
}

type MbThumbnails struct {
	Large string `json:"large"`
	Small string `json:"small"`
}
