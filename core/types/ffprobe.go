package types

type FfprobeStandard struct {
	Filename   string            `json:"file_name"`
	FormatName string            `json:"format_name"`
	Tags       map[string]string `json:"tags"`
	Duration   string            `json:"duration"`
	Size       string            `json:"size"`
	Bitrate    string            `json:"bit_rate"`
	Codec      string            `json:"codec_name"`
	BitDepth   int               `json:"bits_per_sample"`
	SampleRate int               `json:"sample_rate"`
	Channels   int               `json:"channels"`
}

type FfprobeStandardOutput struct {
	Format struct {
		Filename   string            `json:"file_name"`
		FormatName string            `json:"format_name"`
		Tags       map[string]string `json:"tags"`
		Duration   string            `json:"duration"`
		Size       string            `json:"size"`
		Bitrate    string            `json:"bit_rate"`
	} `json:"format"`
	Streams []struct {
		Codec      string `json:"codec_name"`
		BitDepth   int    `json:"bits_per_sample"`
		SampleRate string `json:"sample_rate"`
		Channels   int    `json:"channels"`
	} `json:"streams"`
}

type FfprobeOpusOutput struct {
	Streams []struct {
		Codec      string            `json:"codec_name"`
		Tags       map[string]string `json:"tags"`
		BitDepth   int               `json:"bits_per_sample"`
		SampleRate string            `json:"sample_rate"`
		Channels   int               `json:"channels"`
	} `json:"streams"`
	Format struct {
		Filename   string `json:"file_name"`
		FormatName string `json:"format_name"`
		Duration   string `json:"duration"`
		Size       string `json:"size"`
		Bitrate    string `json:"bit_rate"`
	} `json:"format"`
}
