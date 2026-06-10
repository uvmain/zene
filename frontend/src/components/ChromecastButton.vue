<script setup lang="ts">
import { isBrowserChrome } from '~/logic/browser'

interface CastEventValue<T> { value?: T }

const mediaUrl = ref('')
const connected = ref(false)
const debugEnabled = ref(true)
const playing = ref(false)
const seeking = ref(false)
const duration = ref(0)
const currentTime = ref(0)
const volume = ref(0.7)
const savedVolume = ref(0.7)
const chromecastMuted = ref(false)
const isChromecastReady = ref(false)

function debugLog(logMessage: string) {
  console.log(`[DEBUG] ${logMessage}`)
}

function initialiseChromecast() {
  window.__onGCastApiAvailable = (isAvailable: boolean) => {
    if (isAvailable) {
      setTimeout(() => {
        initializeCastApi()
      }, 50)
    }
  }
  connected.value = cast.framework.CastContext.getInstance().getCurrentSession() !== null
}

function initializeCastApi() {
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
}

function connect() {
  if (!window.cast) {
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

  cast.framework.CastContext.getInstance().requestSession()

  connected.value = cast.framework.CastContext.getInstance().getCurrentSession() !== null
}

function loadMedia() {
  const contentType = 'audio/aac'
  const castSession = cast.framework.CastContext.getInstance().getCurrentSession()
  if (!castSession) {
    return
  }

  const mediaInfo = new window.chrome.cast.media.MediaInfo(mediaUrl.value, contentType)
  const request = new window.chrome.cast.media.LoadRequest(mediaInfo)

  castSession.loadMedia(request).then(
    () => {
      debugLog('[chromecast] - Load succeeded')
    },
    (err: unknown) => {
      debugLog(`[chromecast] - Error:${err}`)
    },
  )
}

function setPlayerEvents() {
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

function play() {
  debugLog('[chromecast:play]')

  const castSession = cast.framework.CastContext.getInstance().getCurrentSession()
  if (!castSession) {
    return
  }

  const media = castSession.getMediaSession()

  if (!media) {
    loadMedia()
    return
  }

  castSession.sendMessage('urn:x-cast:com.google.cast.media', {
    type: 'PLAY',
    requestId: 1,
    mediaSessionId: media.mediaSessionId,
  })
}

function pause() {
  debugLog('[chromecast:pause]')

  const castSession = cast.framework.CastContext.getInstance().getCurrentSession()
  if (!castSession) {
    return
  }

  const media = castSession.getMediaSession()
  if (!media) {
    return
  }

  castSession.sendMessage('urn:x-cast:com.google.cast.media', {
    type: 'PAUSE',
    requestId: 1,
    mediaSessionId: media.mediaSessionId,
  })
}

function stopChromecast() {
  debugLog('[chromecast:stop]')

  const castSession = cast.framework.CastContext.getInstance().getCurrentSession()
  if (!castSession) {
    return
  }

  const media = castSession.getMediaSession()
  if (!media) {
    return
  }

  castSession.sendMessage('urn:x-cast:com.google.cast.media', {
    type: 'STOP',
    requestId: 1,
    mediaSessionId: media.mediaSessionId,
  })
}

function seekTo(value: number) {
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

  castSession.sendMessage('urn:x-cast:com.google.cast.media', {
    type: 'SEEK',
    requestId: 1,
    mediaSessionId: media.mediaSessionId,
    currentTime: value,
  })
  play()
  seeking.value = false
}

function setVolume(value: number) {
  debugLog(`[chromecast:setVolume] - ${value}`)
  const castSession = cast.framework.CastContext.getInstance().getCurrentSession()
  if (!castSession) {
    return
  }

  castSession.setVolume(value)
  volume.value = value

  if (volume.value > 0) {
    chromecastMuted.value = false
  }
}

function setMute() {
  debugLog(`[chromecast:setMute] - ${!chromecastMuted.value}`)
  chromecastMuted.value = !chromecastMuted.value
  if (chromecastMuted.value) {
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

  castSession.setMute(chromecastMuted.value)
}

function onDebugChange() {
  debugLog(`[chromecast:setDebugPanel] - ${debugEnabled.value}`)
}

function onSeekChange(event: Event) {
  const target = event.target as HTMLInputElement | null
  if (!target) {
    return
  }

  debugLog(`[chromecast:onSeekChange] - ${target.value}`)
  seeking.value = true
  if (target.value) {
    seekTo(Number.parseFloat(target.value))
  }
}

function onVolumeChange(event: Event) {
  const target = event.target as HTMLInputElement | null
  if (!target) {
    return
  }

  const nextVolume = Number.parseFloat(target.value)
  debugLog(`[chromecast:onVolumeChange] - ${nextVolume}`)
  volume.value = nextVolume
  if (!Number.isNaN(nextVolume)) {
    setVolume(nextVolume)
  }
}

function onMediaInfoChanged(event: CastEventValue<{ duration?: number }>) {
  debugLog(`[chromecast:onMediaInfoChanged] - ${JSON.stringify(event)}`)
  duration.value = event.value?.duration ?? 0
}

function onCurrentTimeChanged(event: CastEventValue<number>) {
  debugLog(`[chromecast:onCurrentTimeChanged] - ${JSON.stringify(event)}`)
  if (!seeking.value) {
    currentTime.value = event.value ?? 0
  }
}

function onIsConnectedChanged(event: CastEventValue<boolean>) {
  debugLog(`[chromecast:onIsConnectedChanged] - ${JSON.stringify(event)}`)
  connected.value = event.value ?? false
}

function onPlayerStateChanged(event: CastEventValue<string>) {
  debugLog(`[chromecast:onPlayerStateChanged] - ${JSON.stringify(event)}`)
  playing.value = event.value === 'PLAYING' || event.value === 'BUFFERING'
  if (event.value === 'PLAYING') {
    seeking.value = false
  }
}

onMounted(() => {
  initialiseChromecast()
})
</script>

<template>
  <div v-if="isBrowserChrome" class="mb-2.5 inline-block">
    <button v-if="connected" class="mc-button-active" @click="connect">
      <icon-nrk-media-chromecast-active />
    </button>
    <button v-else class="mc-button-primary" @click="connect">
      <icon-nrk-media-chromecast />
    </button>
    <span v-if="!isBrowserChrome" class="ml-2.5">Google Chrome required!</span>
    <button v-if="connected" class="mc-button" @click="loadMedia">
      Load Media
    </button>
    <button v-if="connected" class="mc-button" @click="stopChromecast">
      Stop
    </button>
    <label v-if="connected" class="mx-2.5 inline" for="checkbox">
      <input id="checkbox" v-model="debugEnabled" type="checkbox" @change="onDebugChange">
      <span>Debug Panel</span>
    </label>
  </div>

  <div v-if="connected" class="mc-player-controls">
    <button v-if="playing" class="mc-icon-button mr-[7px]" @click="pause">
      pause_arrow
    </button>
    <button v-else class="mc-icon-button mr-[7px]" @click="play">
      play_arrow
    </button>
    <input class="mc-seek" type="range" step="any" min="0" :max="duration" :value="currentTime" @change="onSeekChange">
    <button class="mc-icon-button mr-[7px] hidden">
      fast_rewind
    </button>
    <button class="mc-icon-button mr-[7px] hidden">
      fast_forward
    </button>
    <button v-if="chromecastMuted" class="mc-icon-button mr-[7px]" @click="setMute">
      volume_mute
    </button>
    <button v-else class="mc-icon-button mr-[7px]" @click="setMute">
      volume_up
    </button>
    <input
      class="mc-volume"
      type="range"
      step="any"
      min="0"
      max="1"
      :value="volume"
      :style="{ background: `linear-gradient(to right, rgb(204, 204, 204) ${volume * 100}%, rgb(0, 0, 0) ${volume * 100}%, rgb(0, 0, 0) 100%)` }"
      @change="onVolumeChange"
    >
  </div>
</template>
