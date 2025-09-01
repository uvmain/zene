package types

type ApiKey struct {
	Id          int    `xml:"id,omitempty" json:"id,omitempty"`
	UserId      int    `xml:"user_id,omitempty" json:"user_id,omitempty"`
	ApiKey      string `xml:"api_key,omitempty" json:"api_key,omitempty"`
	DateCreated string `xml:"date_created,omitempty" json:"date_created,omitempty"`
	LastUsed    string `xml:"last_used,omitempty" json:"last_used,omitempty"`
}

type ApiKeys struct {
	ApiKeys []ApiKey `xml:"apiKey" json:"apiKey"`
}
