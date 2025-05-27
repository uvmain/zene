package types

type FfprobeOutput struct {
	Format struct {
		Filename   string            `json:"filename"`
		FormatName string            `json:"format_name"`
		Tags       map[string]string `json:"tags"`
		Duration   string            `json:"duration"`
		Size       string            `json:"size"`
		Bitrate    string            `json:"bit_rate"`
	} `json:"format"`
}

type FfprobeOpusOutput struct {
	Streams []Stream `json:"streams"`
	Format  struct {
		Filename   string `json:"filename"`
		FormatName string `json:"format_name"`
		Duration   string `json:"duration"`
		Size       string `json:"size"`
		Bitrate    string `json:"bit_rate"`
	} `json:"format"`
}

type Stream struct {
	Tags map[string]string `json:"tags"`
}
