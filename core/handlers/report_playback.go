package handlers

import (
	"net/http"
	"strconv"
	"time"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleReportPlayback(w http.ResponseWriter, r *http.Request) {
	if net.MethodIsNotGetOrPost(w, r) {
		return
	}

	form := net.NormalisedForm(r, w)
	format := form["f"]
	mediaId := form["mediaid"]
	mediaType := form["mediatype"]
	positionMs := form["positionms"]
	state := form["state"]
	playbackRate := form["playbackrate"]
	ignoreScrobble := form["ignorescrobble"]
	playerName := form["c"]

	ctx := r.Context()

	if mediaId == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "mediaId is required", "")
		return
	}
	if mediaType == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "mediaType is required", "")
		return
	}
	if mediaType != "song" && mediaType != "podcast" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "mediaType must be either 'song' or 'podcast'", "")
		return
	}
	if positionMs == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "positionMs is required", "")
		return
	}
	positionMsInt, err := strconv.Atoi(positionMs)
	if err != nil || positionMsInt < 0 {
		logger.Printf("Error parsing positionMs for media %s: %v", mediaId, err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid positionMs parameter, must be a positive integer", "")
		return
	}
	if state == "" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "State is required", "")
		return
	}
	if state != "starting" && state != "playing" && state != "paused" && state != "stopped" {
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "State must be either 'starting', 'playing', 'paused', or 'stopped'", "")
		return
	}
	if playbackRate == "" {
		playbackRate = "1.0"
	}
	playbackRateFloat, err := strconv.ParseFloat(playbackRate, 64)
	if err != nil || playbackRateFloat <= 0 {
		logger.Printf("Error parsing playbackRate for media %s: %v", mediaId, err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid playbackRate parameter, must be a positive float", "")
		return
	}
	if ignoreScrobble == "" {
		ignoreScrobble = "false"
	}
	if ignoreScrobble != "true" && ignoreScrobble != "false" {
		logger.Printf("Error parsing ignoreScrobble for media %s: %v", mediaId, err)
		net.WriteSubsonicError(w, r, types.ErrorMissingParameter, "Invalid ignoreScrobble parameter, must be either 'true' or 'false'", "")
		return
	}

	isValidMediaId, metadataType, _ := database.IsValidMetadataId(ctx, mediaId)
	if !isValidMediaId {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Invalid MediaId passed", "")
		return
	}
	if metadataType != database.MetadataTrack && metadataType != database.MetadataPodcastEpisode {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Invalid MediaId passed", "")
		return
	}
	if (metadataType == database.MetadataTrack && mediaType != "song") || (metadataType == database.MetadataPodcastEpisode && mediaType != "podcast") {
		net.WriteSubsonicError(w, r, types.ErrorDataNotFound, "Invalid MediaId passed", "")
		return
	}

	user, err := database.GetUserByContext(ctx)
	if err != nil {
		logger.Printf("Error getting user by context: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorNotAuthorized, "You do not have permission to report playback", "")
		return
	}

	now := int(time.Now().UnixMilli())
	playedAt := now
	nowPlayingId, err := database.UpsertNowPlaying(ctx, user.Id, mediaId, playedAt, 0, playerName)
	if err != nil {
		logger.Printf("Error upserting now playing for user %d: %v", user.Id, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to upsert user now playing", "")
		return
	}

	previousState, hasPreviousState, err := database.GetPlaybackReportState(ctx, nowPlayingId)
	if err != nil {
		logger.Printf("Error getting previous playback state for user %d: %v", user.Id, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to read playback report state", "")
		return
	}

	err = database.UpsertPlaybackReport(ctx, nowPlayingId, state, positionMsInt, playbackRate, now)
	if err != nil {
		logger.Printf("Error upserting playback report for user %d: %v", user.Id, err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to upsert playback report", "")
		return
	}

	shouldIncrementPlayCount := ignoreScrobble != "true" && state == "stopped" && (!hasPreviousState || previousState != "stopped")
	if shouldIncrementPlayCount {
		err = database.UpsertPlayCount(ctx, user.Id, mediaId)
		if err != nil {
			logger.Printf("Error upserting play count for user %d: %v", user.Id, err)
			net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to upsert user play count", "")
			return
		}
	}

	response := subsonic.GetPopulatedSubsonicResponse(ctx)

	net.WriteSubsonicResponse(w, r, response, format)
}
