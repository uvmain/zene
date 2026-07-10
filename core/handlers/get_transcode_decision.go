package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	stdio "io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"zene/core/config"
	"zene/core/database"
	"zene/core/ffprobe"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetTranscodeDecision(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Unsupported method: "+r.Method, "")
		return
	}

	form := net.NormalisedForm(r, w)
	if form == nil {
		return
	}

	format := form["f"]
	mediaID := firstNonEmpty(form["mediaid"], form["id"])
	mediaType := strings.ToLower(firstNonEmpty(form["mediatype"], form["type"]))

	if mediaID == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "mediaId parameter is required", "")
		return
	}

	if mediaType == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "mediaType parameter is required", "")
		return
	}

	clientInfo, err := decodeClientInfo(r)
	if err != nil {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, err.Error(), "")
		return
	}

	ctx := r.Context()
	mediaPath, err := resolveMediaPath(ctx, mediaID, mediaType)
	if err != nil {
		logger.Printf("Error resolving media path for %s (%s): %v", mediaID, mediaType, err)
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "File not available to stream.", "")
		return
	}

	if mediaPath == "" {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "File not available to stream.", "")
		return
	}

	metadata, err := ffprobe.GetMetadataFromFile(ctx, mediaPath)
	if err != nil {
		logger.Printf("Error probing media file %s: %v", mediaPath, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Unable to inspect media file.", "")
		return
	}

	sourceStream := buildSourceStream(mediaPath, metadata)
	decision := buildTranscodeDecision(mediaID, mediaType, clientInfo, sourceStream)

	response := subsonic.GetPopulatedSubsonicResponse(ctx)
	response.SubsonicResponse.TranscodeDecision = &decision

	net.WriteSubsonicResponse(w, r, response, format)
}

func decodeClientInfo(r *http.Request) (types.ClientInfo, error) {
	var clientInfo types.ClientInfo

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&clientInfo); err != nil {
		if err == stdio.EOF {
			return types.ClientInfo{}, fmt.Errorf("request body is required")
		}
		return types.ClientInfo{}, fmt.Errorf("request body must be a JSON ClientInfo object")
	}

	return clientInfo, nil
}

func resolveMediaPath(ctx context.Context, mediaID string, mediaType string) (string, error) {
	switch mediaType {
	case "song":
		mediaPath, err := database.GetMediaFilePath(ctx, mediaID)
		if err != nil {
			return "", err
		}
		return mediaPath, nil
	case "podcast":
		episode, err := database.GetPodcastEpisodeByGuid(ctx, mediaID)
		if err != nil {
			return "", err
		}
		if episode.Path != "" {
			return episode.Path, nil
		}
		if episode.SourceUrl != "" {
			return episode.SourceUrl, nil
		}
		return "", nil
	default:
		return "", fmt.Errorf("unsupported mediaType: %s", mediaType)
	}
}

func buildSourceStream(mediaPath string, metadata types.FfprobeStandard) types.StreamDetails {
	return types.StreamDetails{
		Protocol:        "http",
		Container:       deriveContainer(mediaPath, metadata),
		Codec:           strings.ToLower(strings.TrimSpace(metadata.Codec)),
		AudioChannels:   metadata.Channels,
		AudioBitrate:    parsePositiveInt(metadata.Bitrate),
		AudioSamplerate: metadata.SampleRate,
		AudioBitdepth:   metadata.BitDepth,
	}
}

func buildTranscodeDecision(mediaID string, mediaType string, clientInfo types.ClientInfo, sourceStream types.StreamDetails) types.TranscodeDecision {
	reasons := directPlayReasons(clientInfo.DirectPlayProfiles, clientInfo, sourceStream)
	canDirectPlay := len(reasons) == 0

	decision := types.TranscodeDecision{
		CanDirectPlay:   canDirectPlay,
		TranscodeReason: reasons,
		SourceStream:    &sourceStream,
	}

	if canDirectPlay {
		return decision
	}

	transcodeProfile, targetFormat := chooseTranscodingProfile(clientInfo.TranscodingProfiles, sourceStream)
	if transcodeProfile == nil {
		decision.ErrorReason = "NoTranscodingProfileMatched"
		return decision
	}

	decision.CanTranscode = true
	transcodeBitrate := chooseTranscodeBitrate(clientInfo, sourceStream)
	transcodeChannels := sourceStream.AudioChannels
	if transcodeProfile.MaxAudioChannels > 0 && (transcodeChannels == 0 || transcodeChannels > transcodeProfile.MaxAudioChannels) {
		transcodeChannels = transcodeProfile.MaxAudioChannels
	}

	decision.TranscodeParams = buildTranscodeParams(mediaID, mediaType, *transcodeProfile, targetFormat, transcodeBitrate, transcodeChannels)
	decision.TranscodeStream = &types.StreamDetails{
		Protocol:        transcodeProfile.Protocol,
		Container:       transcodeProfile.Container,
		Codec:           transcodeProfile.AudioCodec,
		AudioChannels:   transcodeChannels,
		AudioBitrate:    transcodeBitrate,
		AudioSamplerate: sourceStream.AudioSamplerate,
		AudioBitdepth:   sourceStream.AudioBitdepth,
	}

	return decision
}

func directPlayReasons(profiles []types.DirectPlayProfile, clientInfo types.ClientInfo, sourceStream types.StreamDetails) []string {
	if len(profiles) == 0 {
		return []string{"NoDirectPlayProfiles"}
	}

	reasons := make([]string, 0, len(profiles))
	for _, profile := range profiles {
		if directPlayProfileMatches(profile, clientInfo, sourceStream) {
			return nil
		}
		reasons = append(reasons, directPlayProfileReason(profile, clientInfo, sourceStream))
	}

	return dedupeStrings(reasons)
}

func directPlayProfileMatches(profile types.DirectPlayProfile, clientInfo types.ClientInfo, sourceStream types.StreamDetails) bool {
	if clientInfo.MaxAudioBitrate > 0 && sourceStream.AudioBitrate > 0 && sourceStream.AudioBitrate > clientInfo.MaxAudioBitrate {
		return false
	}

	if profile.MaxAudioChannels > 0 && sourceStream.AudioChannels > 0 && sourceStream.AudioChannels > profile.MaxAudioChannels {
		return false
	}

	if !matchesAny(profile.Protocols, sourceStream.Protocol) {
		return false
	}

	if !matchesAny(profile.Containers, sourceStream.Container) {
		return false
	}

	if !matchesAny(profile.AudioCodecs, sourceStream.Codec) {
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
	if !matchesAny(profile.Protocols, sourceStream.Protocol) {
		return "ProtocolNotSupported"
	}
	if !matchesAny(profile.Containers, sourceStream.Container) {
		return "ContainerNotSupported"
	}
	if !matchesAny(profile.AudioCodecs, sourceStream.Codec) {
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
		if !isSupportedTranscodeFormat(profile) {
			continue
		}
		targetFormat := targetFormatForProfile(profile)
		return &profile, targetFormat
	}

	return nil, ""
}

func chooseTranscodeBitrate(clientInfo types.ClientInfo, sourceStream types.StreamDetails) int {
	bitrate := clientInfo.MaxTranscodingAudioBitrate
	if bitrate <= 0 {
		bitrate = clientInfo.MaxAudioBitrate
	}
	if bitrate <= 0 {
		bitrate = config.DefaultBitRate * 1000
	}
	if sourceStream.AudioBitrate > 0 && bitrate > sourceStream.AudioBitrate {
		bitrate = sourceStream.AudioBitrate
	}
	return bitrate
}

func buildTranscodeParams(mediaID string, mediaType string, profile types.TranscodingProfile, targetFormat string, bitrate int, channels int) string {
	payload := map[string]interface{}{
		"mediaId":       mediaID,
		"mediaType":     mediaType,
		"container":     profile.Container,
		"audioCodec":    profile.AudioCodec,
		"protocol":      profile.Protocol,
		"targetFormat":  targetFormat,
		"bitrate":       bitrate,
		"audioChannels": channels,
	}
	encoded, _ := json.Marshal(payload)
	return base64.RawURLEncoding.EncodeToString(encoded)
}

func isSupportedTranscodeFormat(profile types.TranscodingProfile) bool {
	return targetFormatForProfile(profile) != ""
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

func deriveContainer(mediaPath string, metadata types.FfprobeStandard) string {
	if parsedURL, err := url.Parse(mediaPath); err == nil && parsedURL.Path != "" {
		if ext := filepath.Ext(parsedURL.Path); ext != "" {
			return strings.TrimPrefix(strings.ToLower(ext), ".")
		}
	}

	if ext := filepath.Ext(mediaPath); ext != "" {
		return strings.TrimPrefix(strings.ToLower(ext), ".")
	}

	formatName := strings.ToLower(strings.TrimSpace(metadata.FormatName))
	if formatName != "" {
		parts := strings.Split(formatName, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	return ""
}

func parsePositiveInt(value string) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed < 0 {
		return 0
	}
	return parsed
}

func matchesAny(values []string, candidate string) bool {
	if len(values) == 0 {
		return true
	}
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value), strings.TrimSpace(candidate)) {
			return true
		}
	}
	return false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func dedupeStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
