import { debug } from 'node:console'
import { seek as elementSeek } from '~/logic/audioElement'
import { debugLog } from '~/logic/logger'
import { currentTime as uiCurrentTime } from '~/logic/playbackQueue'
import { currentVolume } from '~/logic/volume'

export interface CastEventValue<T> { value?: T }

export const currentMediaUrl = ref<string>('')
export const connected = ref<boolean>(false)
export const playing = ref<boolean>(false)
export const seeking = ref<boolean>(false)
export const duration = ref<number>(0)
export const currentTime = ref<number>(0)
export const volume = ref<number>(0.7)
export const savedVolume = ref<number>(0.7)
export const muted = ref<boolean>(false)
export const isChromecastReady = ref<boolean>(false)

const isChrome = ref<boolean>(false)

export function isBrowserChrome(): boolean {
  return /Chrome/.test(navigator.userAgent)
    && /Google Inc/.test(navigator.vendor)
}

export async function initialiseChromecast() {
  isChrome.value = isBrowserChrome()
  if (!isChrome.value) {
    debugLog('[chromecast] - initialiseChromecast called but browser is not Chrome, exiting')
    return
  }
  window.__onGCastApiAvailable = (isAvailable: boolean) => {
    if (isAvailable) {
      setTimeout(() => {
        initializeCastApi()
      }, 50)
    }
  }
  await cast.framework.CastContext.getInstance().requestSession()
  connected.value = cast.framework.CastContext.getInstance().getCurrentSession() !== null
}

export function initializeCastApi() {
  if (isChromecastReady.value) {
    return
  }

  debugLog(`[chromecast] - Initializing Cast API: ${chrome.cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID}`)
  cast.framework.setLoggerLevel(cast.framework.LoggerLevel.DEBUG)
  cast.framework.CastContext.getInstance().setOptions({
    receiverApplicationId: chrome.cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID,
    autoJoinPolicy: chrome.cast.AutoJoinPolicy.ORIGIN_SCOPED,
  })
  isChromecastReady.value = true
  setPlayerEvents()
  setVolume(currentVolume.value)
}

export async function connect() {
  if (!Object.hasOwn(window, 'cast')) {
    debugLog('[chromecast:connect] - Cast SDK not available yet')
    return
  }

  if (!isChromecastReady.value) {
    initializeCastApi()
  }

  if (!isChromecastReady.value) {
    debugLog('[chromecast:connect] - Cast options are not ready yet')
    return
  }

  await cast.framework.CastContext.getInstance().requestSession()

  connected.value = cast.framework.CastContext.getInstance().getCurrentSession() !== null
}

export function setPlayerEvents() {
  const player = new cast.framework.RemotePlayer()
  const playerController = new cast.framework.RemotePlayerController(player)

  playerController.addEventListener(
    cast.framework.RemotePlayerEventType.IS_CONNECTED_CHANGED,
    onIsConnectedChanged,
  )

  playerController.addEventListener(
    cast.framework.RemotePlayerEventType.IS_MEDIA_LOADED_CHANGED,
    (event: unknown) => {
      debugLog(`[chromecast:onIsMediaLoadedChanged] - ${JSON.stringify(event)}`)
    },
  )

  playerController.addEventListener(
    cast.framework.RemotePlayerEventType.CURRENT_TIME_CHANGED,
    onCurrentTimeChanged,
  )

  playerController.addEventListener(
    cast.framework.RemotePlayerEventType.DURATION_CHANGED,
    (event: unknown) => {
      debugLog(`[chromecast:onDurationChanged] - ${JSON.stringify(event)}`)
    },
  )

  playerController.addEventListener(
    cast.framework.RemotePlayerEventType.MEDIA_INFO_CHANGED,
    onMediaInfoChanged,
  )

  playerController.addEventListener(
    cast.framework.RemotePlayerEventType.PLAYER_STATE_CHANGED,
    onPlayerStateChanged,
  )
}

export async function loadMedia(mediaUrl: string): Promise<void> {
  const contentType = 'audio/aac'
  const castSession = cast.framework.CastContext.getInstance().getCurrentSession()
  if (!castSession) {
    return
  }

  const url = mediaUrl ?? currentMediaUrl.value
  if (!url) {
    debugLog('[chromecast:loadMedia] - No media URL provided')
    return
  }

  const mediaInfo = new window.chrome.cast.media.MediaInfo(url, contentType)
  const request = new window.chrome.cast.media.LoadRequest(mediaInfo)

  castSession.loadMedia(request).then(
    () => {
      debugLog('[chromecast] - Load succeeded')
    },
    (err: unknown) => {
      debugLog(`[chromecast] - Error:${err instanceof Error ? err.message : JSON.stringify(err)}`)
    },
  )
}

export async function play() {
  debugLog('[chromecast:play]')

  const castSession = cast.framework.CastContext.getInstance().getCurrentSession()
  if (!castSession) {
    return
  }

  const media = castSession.getMediaSession()

  if (!media) {
    debugLog('[chromecast:play] - No media session available')
    return
  }

  void castSession.sendMessage('urn:x-cast:com.google.cast.media', {
    type: 'PLAY',
    requestId: 1,
    mediaSessionId: media.mediaSessionId,
  })
}

export function pause() {
  debugLog('[chromecast:pause]')

  const castSession = cast.framework.CastContext.getInstance().getCurrentSession()
  if (!castSession) {
    return
  }

  const media = castSession.getMediaSession()
  if (!media) {
    return
  }

  void castSession.sendMessage('urn:x-cast:com.google.cast.media', {
    type: 'PAUSE',
    requestId: 1,
    mediaSessionId: media.mediaSessionId,
  })
}

export function stop() {
  debugLog('[chromecast:stop]')

  const castSession = cast.framework.CastContext.getInstance().getCurrentSession()
  if (!castSession) {
    return
  }

  const media = castSession.getMediaSession()
  if (!media) {
    return
  }

  void castSession.sendMessage('urn:x-cast:com.google.cast.media', {
    type: 'STOP',
    requestId: 1,
    mediaSessionId: media.mediaSessionId,
  })
}

export function seekTo(value: number) {
  debugLog(`[chromecast:seekTo] - ${value}`)

  seeking.value = true

  const castSession = cast.framework.CastContext.getInstance().getCurrentSession()
  if (!castSession) {
    return
  }

  const media = castSession.getMediaSession()
  if (!media) {
    return
  }

  void castSession.sendMessage('urn:x-cast:com.google.cast.media', {
    type: 'SEEK',
    requestId: 1,
    mediaSessionId: media.mediaSessionId,
    currentTime: value,
  })
  void play()
  seeking.value = false
}

export function setVolume(value: number) {
  debugLog(`[chromecast:setVolume] - ${value}`)
  const castSession = cast.framework.CastContext.getInstance().getCurrentSession()
  if (!castSession) {
    return
  }

  void castSession.setVolume(value)
  volume.value = value

  if (volume.value > 0) {
    muted.value = false
  }
}

export function toggleMute() {
  debugLog(`[chromecast:setMute] - ${!muted.value}`)
  muted.value = !muted.value
  if (muted.value) {
    savedVolume.value = volume.value
    volume.value = 0
  }
  else {
    volume.value = savedVolume.value
  }
  const castSession = cast.framework.CastContext.getInstance().getCurrentSession()
  if (!castSession) {
    return
  }

  void castSession.setMute(muted.value)
}

export function onSeekChange(event: Event) {
  const target = event.target as HTMLInputElement | null
  if (!target) {
    return
  }

  debugLog(`[chromecast:onSeekChange] - ${target.value}`)

  if (target.value) {
    seekTo(Number.parseFloat(target.value))
  }
  elementSeek(Number.parseFloat(target.value))
}

export function onVolumeChange(event: Event) {
  const target = event.target as HTMLInputElement | null
  if (!target) {
    return
  }

  const nextVolume = Number.parseFloat(target.value)
  debugLog(`[chromecast:onVolumeChange] - ${nextVolume}`)
  volume.value = nextVolume
  if (!Number.isNaN(nextVolume)) {
    setVolume(nextVolume)
    currentVolume.value = nextVolume
  }
}

export function onMediaInfoChanged(event: CastEventValue<{ duration?: number }>) {
  debugLog(`[chromecast:onMediaInfoChanged] - ${JSON.stringify(event)}`)
  duration.value = event.value?.duration ?? 0
}

export function onCurrentTimeChanged(event: CastEventValue<number>) {
  debugLog(`[chromecast:onCurrentTimeChanged] - ${JSON.stringify(event)}`)
  if (!seeking.value) {
    currentTime.value = event.value ?? 0
    uiCurrentTime.value = currentTime.value
  }
}

export function onIsConnectedChanged(event: CastEventValue<boolean>) {
  debugLog(`[chromecast:onIsConnectedChanged] - ${JSON.stringify(event)}`)
  connected.value = event.value ?? false
}

export function onPlayerStateChanged(event: CastEventValue<string>) {
  debugLog(`[chromecast:onPlayerStateChanged] - ${JSON.stringify(event)}`)
  playing.value = event.value === 'PLAYING' || event.value === 'BUFFERING'
  if (event.value === 'PLAYING') {
    seeking.value = false
  }
}
