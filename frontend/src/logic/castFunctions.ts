import { seek } from '~/logic/audioElement'
import { isPlaying, updateProgress } from '~/logic/playbackQueue'
import { currentVolume } from '~/logic/volume'
import { castContext, castPlayer, castPlayerController, castSession, chromecastAvailable, chromecastConnected, savedLocalPosition } from './castRefs'
import { debugLog } from './logger'

export function initializeCast() {
  castContext.value = cast.framework.CastContext.getInstance()
  castContext.value.setOptions({
    receiverApplicationId: chrome.cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID,
    autoJoinPolicy: chrome.cast.AutoJoinPolicy.ORIGIN_SCOPED,
  })

  castContext.value.addEventListener(cast.framework.CastContextEventType.CAST_STATE_CHANGED, onCastStateChanged)
  castContext.value.addEventListener(cast.framework.CastContextEventType.SESSION_STATE_CHANGED, onSessionStateChanged)

  castPlayer.value = new cast.framework.RemotePlayer()
  castPlayerController.value = new cast.framework.RemotePlayerController(castPlayer.value)
  castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.IS_PAUSED_CHANGED, onCastPlayerStateChanged)
  castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.CURRENT_TIME_CHANGED, onCastTimeChanged)
  castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.VOLUME_LEVEL_CHANGED, onCastVolumeChanged)
  castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.IS_MUTED_CHANGED, onCastMuteChanged)
  castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.IS_CONNECTED_CHANGED, onCastConnectionChanged)
  castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.MEDIA_INFO_CHANGED, onCastMediaInfoChanged)

  castPlayerController.value.setVolumeLevel(currentVolume.value)

  debugLog('CastContext initialized')
}

function onCastStateChanged(event: cast.framework.CastStateEventData) {
  debugLog(`Cast state changed: ${event.castState}`)
  chromecastAvailable.value = event.castState !== cast.framework.CastState.NO_DEVICES_AVAILABLE
  if (event.castState === cast.framework.CastState.NOT_CONNECTED) {
    debugLog('Cast device disconnected')
    chromecastConnected.value = false
  }
  else if (event.castState === cast.framework.CastState.CONNECTED) {
    debugLog('Cast device connected')
    chromecastConnected.value = true
  }
}

function onSessionStateChanged(event: cast.framework.SessionStateEventData) {
  debugLog(`Session state changed: ${event.sessionState}`)
  if (event.sessionState === cast.framework.SessionState.SESSION_ENDED
    || event.sessionState === cast.framework.SessionState.SESSION_ENDING) {
    debugLog('Cast session ending/ended')

    if (castPlayer.value && castPlayer.value.isConnected) {
      savedLocalPosition.value = castPlayer.value.currentTime
      debugLog(`Final cast position captured: ${savedLocalPosition.value}s`)
    }
  }
  else if (event.sessionState === cast.framework.SessionState.SESSION_STARTED) {
    debugLog('Cast session started')
    castSession.value = castContext.value?.getCurrentSession() as cast.framework.CastSession | null
  }
}

function onCastPlayerStateChanged() {
  debugLog(`Cast player state changed: paused=${castPlayer.value?.isPaused}`)
  if (castPlayer.value) {
    isPlaying.value = !castPlayer.value.isPaused
  }
}

function onCastTimeChanged() {
  debugLog(`Cast player time changed: currentTime=${castPlayer.value?.currentTime}`)
  if (castPlayer.value) {
    seek(castPlayer.value.currentTime)
    updateProgress()
  }
}

function onCastVolumeChanged() {
  debugLog(`Cast player volume changed: volumeLevel=${castPlayer.value?.volumeLevel}`)
  if (castPlayer.value) {
    currentVolume.value = castPlayer.value.isMuted ? 0 : castPlayer.value.volumeLevel
    if (castPlayer.value.isMuted) {
      currentVolume.value = 0
    }
  }
}

function onCastMuteChanged() {
  onCastVolumeChanged()
}

function onCastConnectionChanged() {
  debugLog(`Cast player connection changed: isConnected=${castPlayer.value?.isConnected}`)
  if (castPlayer.value && !castPlayer.value.isConnected) {
    debugLog('Cast device disconnected unexpectedly')
    if (castPlayer.value.currentTime > 0) {
      savedLocalPosition.value = castPlayer.value.currentTime
      debugLog(`Position captured on disconnect: ${savedLocalPosition.value}s`)
    }
  }
}

function onCastMediaInfoChanged() {
  debugLog('Cast player media info changed')
  if (!castPlayer.value || !castPlayer.value.isConnected || !castPlayer.value.mediaInfo) {
    return
  }

  debugLog('Cast media info changed, checking for track changes')

  const remoteMediaUrl = castPlayer.value.mediaInfo.contentId
  debugLog(`Remote media URL: ${remoteMediaUrl}`)
}

export function cleanupCastPlayer() {
  if (castPlayerController.value) {
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.IS_PAUSED_CHANGED, onCastPlayerStateChanged)
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.CURRENT_TIME_CHANGED, onCastTimeChanged)
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.VOLUME_LEVEL_CHANGED, onCastVolumeChanged)
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.IS_MUTED_CHANGED, onCastMuteChanged)
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.IS_CONNECTED_CHANGED, onCastConnectionChanged)
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.MEDIA_INFO_CHANGED, onCastMediaInfoChanged)
  }

  if (castContext.value) {
    castContext.value.removeEventListener(cast.framework.CastContextEventType.CAST_STATE_CHANGED, onCastStateChanged)
    castContext.value.removeEventListener(cast.framework.CastContextEventType.SESSION_STATE_CHANGED, onSessionStateChanged)
  }

  castPlayer.value = null
  castPlayerController.value = null
}
