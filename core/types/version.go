package types

type Version struct {
	ServerVersion          string `xml:"server_version,attr" json:"server_version"`
	DatabaseVersion        string `xml:"database_version,attr" json:"database_version"`
	SubsonicApiVersion     string `xml:"subsonic_api_version,attr" json:"subsonic_api_version"`
	OpenSubsonicApiVersion string `xml:"open_subsonic_api_version,attr" json:"open_subsonic_api_version"`
	Timestamp              string `xml:"timestamp,attr" json:"timestamp"`
}

type Versions []Version
