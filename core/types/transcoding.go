package types

type ClientInfo struct {
	Name                       string               `json:"name"`
	Platform                   string               `json:"platform"`
	MaxAudioBitrate            int                  `json:"maxAudioBitrate"`
	MaxTranscodingAudioBitrate int                  `json:"maxTranscodingAudioBitrate"`
	DirectPlayProfiles         []DirectPlayProfile  `json:"directPlayProfiles"`
	TranscodingProfiles        []TranscodingProfile `json:"transcodingProfiles"`
	CodecProfiles              []CodecProfile       `json:"codecProfiles"`
}

type DirectPlayProfile struct {
	Containers       []string `json:"containers"`
	AudioCodecs      []string `json:"audioCodecs"`
	Protocols        []string `json:"protocols"`
	MaxAudioChannels int      `json:"maxAudioChannels"`
}

type TranscodingProfile struct {
	Container        string `json:"container"`
	AudioCodec       string `json:"audioCodec"`
	Protocol         string `json:"protocol"`
	MaxAudioChannels int    `json:"maxAudioChannels"`
}

type CodecProfile struct {
	Type        string       `json:"type"`
	Name        string       `json:"name"`
	Limitations []Limitation `json:"limitations"`
}

type Limitation struct {
	Name       string   `json:"name"`
	Comparison string   `json:"comparison"`
	Values     []string `json:"values"`
	Required   bool     `json:"required"`
}

type StreamDetails struct {
	Protocol        string `xml:"protocol,attr" json:"protocol"`
	Container       string `xml:"container,attr" json:"container"`
	Codec           string `xml:"codec,attr" json:"codec"`
	AudioChannels   int    `xml:"audioChannels,attr,omitempty" json:"audioChannels,omitempty"`
	AudioBitrate    int    `xml:"audioBitrate,attr,omitempty" json:"audioBitrate,omitempty"`
	AudioProfile    string `xml:"audioProfile,attr,omitempty" json:"audioProfile,omitempty"`
	AudioSamplerate int    `xml:"audioSamplerate,attr,omitempty" json:"audioSamplerate,omitempty"`
	AudioBitdepth   int    `xml:"audioBitdepth,attr,omitempty" json:"audioBitdepth,omitempty"`
}

type TranscodeDecision struct {
	CanDirectPlay   bool           `xml:"canDirectPlay,attr" json:"canDirectPlay"`
	CanTranscode    bool           `xml:"canTranscode,attr" json:"canTranscode"`
	TranscodeReason []string       `xml:"transcodeReason,omitempty" json:"transcodeReason,omitempty"`
	ErrorReason     string         `xml:"errorReason,attr,omitempty" json:"errorReason,omitempty"`
	TranscodeParams string         `xml:"transcodeParams,attr,omitempty" json:"transcodeParams,omitempty"`
	SourceStream    *StreamDetails `xml:"sourceStream,omitempty" json:"sourceStream,omitempty"`
	TranscodeStream *StreamDetails `xml:"transcodeStream,omitempty" json:"transcodeStream,omitempty"`
}
