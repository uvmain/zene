package types

type InternetRadio struct {
	Id          string `xml:"id,attr" json:"id"`
	Name        string `xml:"name,attr" json:"name"`
	StreamUrl   string `xml:"streamUrl,attr" json:"streamUrl"`
	HomepageUrl string `xml:"homepageUrl,attr,omitempty" json:"homepageUrl,omitempty"`
}

type InternetRadioStations struct {
	InternetRadio []InternetRadio `xml:"internetRadioStation" json:"internetRadioStation"`
}
