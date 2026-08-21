package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/net"
	"zene/core/types"
)

func HandleShareTranscodeDecision(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotPost(w, r) {
		return
	}

	token := r.PathValue("share_token")
	ctx := r.Context()
	form := net.NormalisedForm(r, w)
	mediaId := form["mediaid"]

	if mediaId == "" {
		http.Error(w, "mediaId parameter is required", http.StatusBadRequest)
		return
	}

	_, tokenIsValid := database.ValidateTokenAndMediaId(ctx, token, mediaId)
	if !tokenIsValid {
		logger.Printf("Error validating share token %s", token)
		http.Error(w, "invalid token and media_id combo", http.StatusBadRequest)
		return
	}

	var clientInfo types.ClientInfo

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&clientInfo); err != nil {
		if err == io.EOF {
			net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "request body is required", "")
			return
		}
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	mediaType := "song"

	metadata, err := GetMediaTranscodeMetadata(ctx, mediaId, mediaType)
	if err != nil {
		logger.Printf("Error fetching media metadata for %s: %v", mediaId, err)
		http.Error(w, "Error fetching media metadata", http.StatusInternalServerError)
		return
	}

	sourceStream := types.StreamDetails{
		Protocol:        "http",
		Container:       metadata.Container,
		Codec:           strings.ToLower(strings.TrimSpace(metadata.Codec)),
		AudioChannels:   metadata.AudioChannels,
		AudioBitrate:    metadata.AudioBitrate,
		AudioSamplerate: metadata.AudioSamplerate,
		AudioBitdepth:   metadata.AudioBitdepth,
	}

	decision, transcodeParams := logic.BuildTranscodeDecision(mediaId, mediaType, clientInfo, sourceStream)

	err = database.UpsertTranscodeParams(ctx, types.TranscodeParamsRow{
		ParamString:   decision.TranscodeParams,
		MediaId:       mediaId,
		MediaType:     mediaType,
		Container:     transcodeParams.Container,
		AudioCodec:    transcodeParams.AudioCodec,
		Protocol:      transcodeParams.Protocol,
		TargetFormat:  transcodeParams.TargetFormat,
		Bitrate:       transcodeParams.Bitrate,
		AudioChannels: transcodeParams.AudioChannels,
	})
	if err != nil {
		logger.Printf("Error upserting transcode params for %s: %v", mediaId, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(decision); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
