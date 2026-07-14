package logic

import (
	"encoding/base64"
	"fmt"
	"slices"
	"strings"
	"zene/core/config"
	"zene/core/types"
)

func BuildTranscodeDecision(mediaId string, mediaType string, clientInfo types.ClientInfo, sourceStream types.StreamDetails) (types.TranscodeDecision, types.TranscodeParamsStruct) {
	clientInfo.MaxAudioBitrate = int(float32(clientInfo.MaxAudioBitrate) * float32(1.1)) // allow a small overhead for bitrate
	reasons := cannotDirectPlayReasons(clientInfo.DirectPlayProfiles, clientInfo, sourceStream)
	canDirectPlay := len(reasons) == 0

	decision := types.TranscodeDecision{
		CanDirectPlay:   canDirectPlay,
		TranscodeReason: reasons,
		SourceStream:    &sourceStream,
	}

	if canDirectPlay {
		return decision, types.TranscodeParamsStruct{}
	}

	transcodeProfile, targetFormat := chooseTranscodingProfile(clientInfo.TranscodingProfiles, sourceStream)
	if transcodeProfile == nil {
		decision.ErrorReason = "NoTranscodingProfileMatched"
		return decision, types.TranscodeParamsStruct{}
	}

	decision.CanTranscode = true

	transcodeBitrate := clientInfo.MaxTranscodingAudioBitrate
	if transcodeBitrate <= 0 {
		transcodeBitrate = clientInfo.MaxAudioBitrate
	}
	if transcodeBitrate <= 0 {
		transcodeBitrate = config.DefaultBitRate * 1000
	}
	if sourceStream.AudioBitrate > 0 && transcodeBitrate > sourceStream.AudioBitrate {
		transcodeBitrate = sourceStream.AudioBitrate
	}

	transcodeChannels := sourceStream.AudioChannels
	if transcodeProfile.MaxAudioChannels > 0 && (transcodeChannels == 0 || transcodeChannels > transcodeProfile.MaxAudioChannels) {
		transcodeChannels = transcodeProfile.MaxAudioChannels
	}

	decision.TranscodeStream = &types.StreamDetails{
		Protocol:        transcodeProfile.Protocol,
		Container:       transcodeProfile.Container,
		Codec:           transcodeProfile.AudioCodec,
		AudioChannels:   transcodeChannels,
		AudioBitrate:    transcodeBitrate,
		AudioSamplerate: sourceStream.AudioSamplerate,
		AudioBitdepth:   sourceStream.AudioBitdepth,
	}

	transcodeParams := types.TranscodeParamsStruct{
		MediaID:       mediaId,
		MediaType:     mediaType,
		Container:     transcodeProfile.Container,
		AudioCodec:    transcodeProfile.AudioCodec,
		Protocol:      transcodeProfile.Protocol,
		TargetFormat:  targetFormat,
		Bitrate:       transcodeBitrate,
		AudioChannels: transcodeChannels,
	}

	decision.TranscodeParams = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v%s", transcodeParams, GetCurrentTimeFormatted())))

	return decision, transcodeParams
}

func cannotDirectPlayReasons(profiles []types.DirectPlayProfile, clientInfo types.ClientInfo, sourceStream types.StreamDetails) []string {
	if len(profiles) == 0 {
		return []string{"NoDirectPlayProfiles"}
	}

	reasons := make([]string, 0, len(profiles))
	for _, profile := range profiles {
		if canDirectPlay(profile, clientInfo, sourceStream) {
			return nil
		}
		reasons = append(reasons, directPlayProfileReason(profile, clientInfo, sourceStream))
	}

	return slices.Compact(reasons)
}

func canDirectPlay(profile types.DirectPlayProfile, clientInfo types.ClientInfo, sourceStream types.StreamDetails) bool {
	if clientInfo.MaxAudioBitrate < 1 || sourceStream.AudioBitrate < 1 {
		return false
	}
	if sourceStream.AudioBitrate > clientInfo.MaxAudioBitrate {
		return false
	}

	if profile.MaxAudioChannels > 0 && sourceStream.AudioChannels > 0 && sourceStream.AudioChannels > profile.MaxAudioChannels {
		return false
	}

	if !slices.Contains(profile.Protocols, sourceStream.Protocol) {
		return false
	}

	if !slices.Contains(profile.Containers, sourceStream.Container) {
		return false
	}

	if !slices.Contains(profile.AudioCodecs, sourceStream.Codec) {
		return false
	}

	return true
}

func directPlayProfileReason(profile types.DirectPlayProfile, clientInfo types.ClientInfo, sourceStream types.StreamDetails) string {
	if clientInfo.MaxAudioBitrate > 0 && sourceStream.AudioBitrate > 0 && sourceStream.AudioBitrate > clientInfo.MaxAudioBitrate {
		return "AudioBitrateNotSupported"
	}
	if profile.MaxAudioChannels > 0 && sourceStream.AudioChannels > 0 && sourceStream.AudioChannels > profile.MaxAudioChannels {
		return "AudioChannelsNotSupported"
	}
	if !slices.Contains(profile.Protocols, sourceStream.Protocol) {
		return "ProtocolNotSupported"
	}
	if !slices.Contains(profile.Containers, sourceStream.Container) {
		return "ContainerNotSupported"
	}
	if !slices.Contains(profile.AudioCodecs, sourceStream.Codec) {
		return "AudioCodecNotSupported"
	}
	return "UnsupportedProfile"
}

func chooseTranscodingProfile(profiles []types.TranscodingProfile, sourceStream types.StreamDetails) (*types.TranscodingProfile, string) {
	for i := range profiles {
		profile := profiles[i]
		if profile.Container == "" || profile.AudioCodec == "" || profile.Protocol == "" {
			continue
		}
		if profile.MaxAudioChannels > 0 && sourceStream.AudioChannels > 0 && sourceStream.AudioChannels > profile.MaxAudioChannels {
			continue
		}

		targetFormat := targetFormatForProfile(profile)

		if targetFormat == "" {
			continue
		}

		return &profile, targetFormat
	}

	return nil, ""
}

func targetFormatForProfile(profile types.TranscodingProfile) string {
	container := strings.ToLower(strings.TrimSpace(profile.Container))
	codec := strings.ToLower(strings.TrimSpace(profile.AudioCodec))

	switch {
	case codec == "mp3" || container == "mp3":
		return "mp3"
	case codec == "aac" || container == "aac" || container == "m4a" || container == "mp4":
		return "aac"
	case codec == "opus" || container == "opus":
		return "opus"
	case codec == "flac" || container == "flac":
		return "flac"
	case codec == "vorbis" || container == "ogg" || container == "vorbis":
		return "vorbis"
	case codec == "wav" || container == "wav":
		return "wav"
	case codec == "alac" || container == "alac":
		return "alac"
	case codec == "wma" || container == "wma" || container == "asf":
		return "wma"
	case codec == "aac_latm":
		return "aac_latm"
	default:
		return ""
	}
}
