<script setup lang="ts">
import type { TokenResponse } from '../types/auth'
import { onKeyStroke } from '@vueuse/core'
import { formatTime } from '../composables/logic'
import { useBackendFetch } from '../composables/useBackendFetch'
import { useDebug } from '../composables/useDebug'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'
import { usePlaycounts } from '../composables/usePlaycounts'
import { useRouteTracks } from '../composables/useRouteTracks'
import { useSettings } from '../composables/useSettings'

const { getMimeType, getTemporaryToken } = useBackendFetch()
const { debugLog } = useDebug()
const { clearQueue, currentlyPlayingTrack, resetCurrentlyPlayingTrack, getNextTrack, getPreviousTrack, refreshRandomSeed, getRandomTracks, currentQueue, setCurrentQueue } = usePlaybackQueue()
const { streamQuality } = useSettings()
const { routeTracks } = useRouteTracks()
const { postPlaycount, updatePlaycount } = usePlaycounts()
const router = useRouter()

const audioRef = ref<HTMLAudioElement | null>(null)
const isPlaying = ref(false)
const playcountPosted = ref(false)
const currentTime = ref(0)
const previousVolume = ref(1)
const currentVolume = ref(1)
const isPlayPauseActive = ref(false)
const session = ref<cast.framework.CastSession | null>(null)
const castPlayer = ref<cast.framework.RemotePlayer | null>(null)
const castPlayerController = ref<cast.framework.RemotePlayerController | null>(null)
const isCasting = ref(false)
const castProgressInterval = ref<NodeJS.Timeout | null>(null)
const castUrl = ref<string | null>()
const temporaryToken = ref<TokenResponse | null>(null)

const trackUrl = computed<string>(() => {
  return currentlyPlayingTrack.value?.musicbrainz_track_id ? `/api/tracks/${currentlyPlayingTrack.value.musicbrainz_track_id}/stream?quality=${streamQuality.value}` : ''
})

async function togglePlayback() {
  if (!audioRef.value) {
    console.error('Audio element not found')
    return
  }

  if (!currentQueue.value?.tracks?.length && routeTracks.value.length) {
    setCurrentQueue(routeTracks.value)
  }
  else if (currentQueue.value?.tracks?.length && !currentlyPlayingTrack.value) {
    setCurrentQueue(currentQueue.value?.tracks)
  }

  // Handle casting
  const context = cast.framework.CastContext.getInstance()
  session.value = context.getCurrentSession()

  if (session.value && isCasting.value) {
    // Control cast playback
    if (castPlayerController.value) {
      castPlayerController.value.playOrPause()
    }
  }
  else if (session.value && !isCasting.value) {
    // Start casting
    debugLog('Starting cast session')
    await castAudio()
    return
  }
  else {
    // Control local playback
    if (isPlaying.value) {
      audioRef.value.pause()
    }
    else {
      audioRef.value.play()
    }
  }

  if (!isPlayPauseActive.value) {
    isPlayPauseActive.value = true
    setTimeout(() => {
      isPlayPauseActive.value = false
    }, 200)
  }

  if (!isCasting.value) {
    updateIsPlaying()
  }
}

function toggleMute() {
  if (isCasting.value && castPlayerController.value) {
    // Control cast device mute
    castPlayerController.value.muteOrUnmute()
  }
  else if (audioRef.value) {
    // Control local audio mute
    debugLog('Changing volume')
    if (audioRef.value.volume !== 0) {
      previousVolume.value = audioRef.value.volume
      audioRef.value.volume = 0
      currentVolume.value = 0
    }
    else {
      audioRef.value.volume = previousVolume.value
      currentVolume.value = previousVolume.value
    }
  }
}

async function stopPlayback() {
  if (isCasting.value && castPlayerController.value) {
    // Stop cast playback
    castPlayerController.value.stop()
  }
  else if (audioRef.value) {
    // Stop local playback
    if (audioRef.value.currentTime < 1) {
      resetCurrentlyPlayingTrack()
      clearQueue()
    }
    audioRef.value.pause()
    audioRef.value.load()
  }

  isPlaying.value = false
}

function updateIsPlaying() {
  if (!audioRef.value)
    return
  isPlaying.value = !audioRef.value.paused
}

function updateProgress() {
  // Get progress from cast player if casting, otherwise from local audio
  if (isCasting.value && castPlayer.value) {
    currentTime.value = castPlayer.value.currentTime
  }
  else if (audioRef.value) {
    currentTime.value = audioRef.value.currentTime
  }
  else {
    return
  }

  if (currentlyPlayingTrack.value && !playcountPosted.value) {
    const halfwayPoint = Number.parseFloat(currentlyPlayingTrack.value.duration) / 2
    if (currentTime.value >= halfwayPoint) {
      postPlaycount(currentlyPlayingTrack.value.musicbrainz_track_id)
      updatePlaycount(currentlyPlayingTrack.value.musicbrainz_track_id)
      playcountPosted.value = true
    }
  }
}

watch(currentlyPlayingTrack, (newTrack, oldTrack) => {
  if (newTrack !== oldTrack) {
    playcountPosted.value = false
  }
})

function seek(event: Event) {
  const target = event.target as HTMLInputElement
  const seekTime = Number.parseFloat(target.value)

  if (isCasting.value && castPlayer.value && castPlayerController.value) {
    // Seek cast playback
    castPlayer.value.currentTime = seekTime
    castPlayerController.value.seek()
  }
  else if (audioRef.value) {
    // Seek local playback
    audioRef.value.currentTime = seekTime
  }
}

function volumeInput(event: Event) {
  const target = event.target as HTMLInputElement
  const volume = Number.parseFloat(target.value)

  if (isCasting.value && castPlayer.value && castPlayerController.value) {
    // Control cast device volume
    castPlayer.value.volumeLevel = volume
    castPlayerController.value.setVolumeLevel()
  }
  else if (audioRef.value) {
    // Control local volume
    audioRef.value.volume = volume
  }

  currentVolume.value = volume
}

async function handleNextTrack() {
  if (isCasting.value && castPlayerController.value) {
    // For cast, try to use the native next track functionality
    castPlayerController.value.nextTrack()
  }

  // Always update the local queue for track management
  await getNextTrack()

  // If casting, load the new track to the cast device
  if (isCasting.value && currentlyPlayingTrack.value) {
    await castAudio()
  }
}

async function handlePreviousTrack() {
  if (isCasting.value && castPlayerController.value) {
    // For cast, try to use the native previous track functionality
    castPlayerController.value.previousTrack()
  }

  // Always update the local queue for track management
  await getPreviousTrack()

  // If casting, load the new track to the cast device
  if (isCasting.value && currentlyPlayingTrack.value) {
    await castAudio()
  }
}

async function handleGetRandomTracks() {
  refreshRandomSeed()
  await getRandomTracks()
  router.push('/queue')
}

onKeyStroke('MediaPlayPause', (e) => {
  e.preventDefault()
  togglePlayback()
})

onKeyStroke('MediaTrackPrevious', (e) => {
  e.preventDefault()
  handlePreviousTrack()
})

onKeyStroke('MediaTrackNext', (e) => {
  e.preventDefault()
  handleNextTrack()
})

onKeyStroke('MediaStop', (e) => {
  e.preventDefault()
  stopPlayback()
})

watch(currentlyPlayingTrack, (newTrack, oldTrack) => {
  const audio = audioRef.value
  if (!audio) {
    return
  }
  if (newTrack && newTrack.musicbrainz_track_id !== oldTrack?.musicbrainz_track_id) {
    audio.pause()
    audio.load()
    audio.addEventListener(
      'canplaythrough',
      () => {
        audio?.play()
      },
      { once: true },
    )
  }
  else if (!newTrack) {
    audio.pause()
    audio.removeAttribute('src')
    audio.load()
    currentTime.value = 0
    isPlaying.value = false
  }
})

async function castAudio() {
  const context = cast.framework.CastContext.getInstance()
  session.value = context.getCurrentSession()

  if (!session.value) {
    console.error('No active cast session found')
    return
  }

  if (!trackUrl.value) {
    console.error('No track URL available for casting')
    return
  }

  if (trackUrl.value.includes('?')) {
    castUrl.value = `${trackUrl.value}&token=${temporaryToken.value?.token}`
  }
  else {
    castUrl.value = `${trackUrl.value}?token=${temporaryToken.value?.token}`
  }
  // prefix base url to requestUrl
  if (window) {
    const protocol = window.location.protocol
    const host = window.location.host
    castUrl.value = `${protocol}//${host}${castUrl.value}`
  }

  const contentType = await getMimeType(castUrl.value)
  if (!contentType) {
    console.error('Could not determine content type for casting')
    return
  }

  debugLog(`Casting URL: ${castUrl.value} with content type: ${contentType}`)
  const mediaInfo = new chrome.cast.media.MediaInfo(castUrl.value, contentType)

  // Add metadata for better cast experience
  if (currentlyPlayingTrack.value) {
    mediaInfo.metadata = {
      title: currentlyPlayingTrack.value.title,
      artist: currentlyPlayingTrack.value.artist || currentlyPlayingTrack.value.album_artist || '',
      images: currentlyPlayingTrack.value.image_url ? [{ url: currentlyPlayingTrack.value.image_url }] : [],
    }
  }

  const request = new chrome.cast.media.LoadRequest(mediaInfo)

  try {
    await session.value.loadMedia(request)
    debugLog('Media loaded to cast device')
    isCasting.value = true

    // Pause local audio when casting starts
    if (audioRef.value) {
      audioRef.value.pause()
    }

    // Update cast state
    updateCastState()
  }
  catch (err) {
    console.error('Error loading media:', err)
  }
}

function initializeCast() {
  const context = cast.framework.CastContext.getInstance()
  context.setOptions({
    receiverApplicationId: chrome.cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID,
    autoJoinPolicy: chrome.cast.AutoJoinPolicy.ORIGIN_SCOPED,
  })

  debugLog('CastContext initialized')
}

function onCastStateChanged(event: any) {
  debugLog(`Cast state changed:', ${event.castState}`)
  updateCastState()
}

function onSessionStateChanged(event: any) {
  debugLog(`Session state changed:', ${event.sessionState}`)
  updateCastState()
}

function updateCastState() {
  const context = cast.framework.CastContext.getInstance()
  session.value = context.getCurrentSession()

  if (session.value) {
    isCasting.value = true
    setupCastPlayer()
  }
  else {
    isCasting.value = false
    cleanupCastPlayer()
  }
}

function setupCastPlayer() {
  if (!castPlayer.value) {
    castPlayer.value = new cast.framework.RemotePlayer()
    castPlayerController.value = new cast.framework.RemotePlayerController(castPlayer.value)

    castPlayerController.value.addEventListener(cast.framework.CastContextEventType.CAST_STATE_CHANGED, onCastStateChanged)
    castPlayerController.value.addEventListener(cast.framework.CastContextEventType.SESSION_STATE_CHANGED, onSessionStateChanged)

    // Listen for remote player events
    castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.IS_PAUSED_CHANGED, onCastPlayerStateChanged)
    castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.CURRENT_TIME_CHANGED, onCastTimeChanged)
    castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.VOLUME_LEVEL_CHANGED, onCastVolumeChanged)
    castPlayerController.value.addEventListener(cast.framework.RemotePlayerEventType.IS_MUTED_CHANGED, onCastMuteChanged)

    // Start progress tracking for cast
    startCastProgressTracking()
  }
}

function cleanupCastPlayer() {
  if (castPlayerController.value) {
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.IS_PAUSED_CHANGED, onCastPlayerStateChanged)
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.CURRENT_TIME_CHANGED, onCastTimeChanged)
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.VOLUME_LEVEL_CHANGED, onCastVolumeChanged)
    castPlayerController.value.removeEventListener(cast.framework.RemotePlayerEventType.IS_MUTED_CHANGED, onCastMuteChanged)
  }

  stopCastProgressTracking()
  castPlayer.value = null
  castPlayerController.value = null
}

function startCastProgressTracking() {
  if (castProgressInterval.value) {
    clearInterval(castProgressInterval.value)
  }

  castProgressInterval.value = setInterval(() => {
    if (castPlayer.value && castPlayer.value.isConnected) {
      updateProgress()
    }
  }, 1000) // Update every second
}

function stopCastProgressTracking() {
  if (castProgressInterval.value) {
    clearInterval(castProgressInterval.value)
    castProgressInterval.value = null
  }
}

function onCastPlayerStateChanged() {
  if (castPlayer.value) {
    isPlaying.value = !castPlayer.value.isPaused
  }
}

function onCastTimeChanged() {
  if (castPlayer.value) {
    currentTime.value = castPlayer.value.currentTime
  }
}

function onCastVolumeChanged() {
  if (castPlayer.value) {
    currentVolume.value = castPlayer.value.volumeLevel
  }
}

function onCastMuteChanged() {
  if (castPlayer.value) {
    currentVolume.value = castPlayer.value.isMuted ? 0 : castPlayer.value.volumeLevel
  }
}

onMounted(async () => {
  temporaryToken.value = await getTemporaryToken()
  debugLog('Waiting for Cast SDK...')

  // Hook for when the SDK becomes available (in case it's not already)
  window.__onGCastApiAvailable = (isAvailable: boolean) => {
    debugLog(`Cast API available (async): ${isAvailable}`)
    if (isAvailable) {
      initializeCast()
      updateCastState() // Initialize cast state
    }
    else {
      console.warn('Cast API not available')
    }
  }

  // 🔥 If the SDK already loaded and called __onGCastApiAvailable BEFORE this script ran
  // we have to check manually and initialize right now
  if ((window.cast && window.cast.isAvailable) || (window.chrome?.cast && window.chrome.cast.isAvailable)) {
    debugLog('Cast API already available (sync)')
    initializeCast()
    updateCastState() // Initialize cast state
  }
})

onMounted(() => {
  const audio = audioRef.value
  if (!audio) {
    return
  }
  audio.addEventListener('play', updateIsPlaying)
  audio.addEventListener('pause', updateIsPlaying)
  audio.addEventListener('timeupdate', updateProgress)
  audio.addEventListener('ended', () => {
    if (currentQueue.value && currentQueue.value.tracks.length > 0) {
      handleNextTrack()
    }
    else {
      isPlaying.value = false
    }
  })
})

onUnmounted(() => {
  const audio = audioRef.value
  if (audio) {
    audio.removeEventListener('play', updateIsPlaying)
    audio.removeEventListener('pause', updateIsPlaying)
    audio.removeEventListener('timeupdate', updateProgress)
    audio.removeEventListener('ended', () => {
      if (currentQueue.value && currentQueue.value.tracks.length > 0) {
        handleNextTrack()
      }
      else {
        isPlaying.value = false
      }
    })

    audio.pause()
    audio.removeAttribute('src')
    audio.load()
  }

  // Clean up cast resources
  cleanupCastPlayer()

  // Remove cast context listeners
  const context = cast.framework.CastContext.getInstance()
  if (context) {
    context.removeEventListener(cast.framework.CastContextEventType.CAST_STATE_CHANGED, onCastStateChanged)
    context.removeEventListener(cast.framework.CastContextEventType.SESSION_STATE_CHANGED, onSessionStateChanged)
  }
})
</script>

<template>
  <footer
    class="sticky bottom-0 mt-auto w-full bg-zene-700 bg-cover bg-center"
    :class="{ 'animate-pulse-bg': currentlyPlayingTrack && isPlaying }"
    :style="{ backgroundImage: `url(${currentlyPlayingTrack?.image_url})` }"
  >
    <div>
      {{ castUrl }}
    </div>
    <div class="flex flex-col items-center border-0 border-t-1 border-white/20 border-solid px-2 backdrop-blur-2xl backdrop-contrast-30 md:flex-row space-y-2 md:px-4 md:space-x-2 md:space-y-0">
      <div
        class="h-full w-full flex flex-grow flex-col items-center justify-center py-2 space-y-2 md:py-2 md:space-y-2"
      >
        <audio ref="audioRef" :src="trackUrl" preload="metadata" class="hidden" />
        <div class="">
          <!-- Progress Bar -->
          <div v-if="audioRef" class="max-w-xs w-full flex flex-row items-center gap-2 lg:max-w-200 md:max-w-lg sm:max-w-md md:gap-2">
            <span id="currentTime" class="w-10 text-right text-sm text-gray-2 md:w-12 sm:w-10 sm:text-sm">
              {{ formatTime(currentTime) }}
            </span>
            <input
              type="range"
              class="h-2 w-full cursor-pointer bg-white/60 accent-zene-200 md:h-1"
              :max="currentlyPlayingTrack ? currentlyPlayingTrack.duration : 0"
              :value="currentTime"
              @input="seek"
            />
            <span id="duration" class="w-10 text-sm text-gray-2 md:w-12 sm:w-10 sm:text-sm">
              {{ formatTime(currentlyPlayingTrack ? Number.parseFloat(currentlyPlayingTrack.duration) : 0) }}
            </span>
          </div>

          <!-- Buttons -->
          <div class="mt-2 flex flex-row items-center justify-center gap-x-2 md:mt-2 md:gap-x-4 sm:gap-x-2">
            <button id="repeat" class="h-10 w-10 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none md:h-12 md:w-12 sm:h-10 sm:w-10" @click="stopPlayback()">
              <icon-tabler-player-stop class="text-lg md:text-xl sm:text-lg" />
            </button>
            <button id="shuffle" class="h-10 w-10 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none md:h-12 md:w-12 sm:h-10 sm:w-10" @click="togglePlayback()">
              <icon-tabler-arrows-shuffle class="text-lg md:text-xl sm:text-lg" />
            </button>
            <button id="back" class="h-10 w-10 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none md:h-12 md:w-12 sm:h-10 sm:w-10" @click="handlePreviousTrack()">
              <icon-tabler-player-skip-back class="text-lg md:text-xl sm:text-lg" />
            </button>
            <button
              id="play-pause"
              class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-md border-none text-white font-semibold outline-none md:h-12 md:w-12 sm:h-12 sm:w-12"
              :class="isPlayPauseActive ? 'bg-zene-200' : 'bg-zene-400 transition-colors duration-200'"
              @click="togglePlayback()"
            >
              <icon-tabler-player-play v-if="!isPlaying" class="text-2xl md:text-3xl sm:text-2xl" />
              <icon-tabler-player-pause v-else class="text-2xl md:text-3xl sm:text-2xl" />
            </button>
            <button id="forward" class="h-10 w-10 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none md:h-12 md:w-12 sm:h-10 sm:w-10" @click="handleNextTrack()">
              <icon-tabler-player-skip-forward class="text-lg md:text-xl sm:text-lg" />
            </button>
            <button id="repeat" class="h-10 w-10 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none md:h-12 md:w-12 sm:h-10 sm:w-10" @click="togglePlayback()">
              <icon-tabler-repeat class="text-lg md:text-xl sm:text-lg" />
            </button>
            <button
              id="shuffle"
              class="h-10 w-10 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none md:h-12 md:w-12 sm:h-10 sm:w-10"
              @click="handleGetRandomTracks()"
            >
              <icon-tabler-dice-3 class="text-lg md:text-xl sm:text-lg" />
            </button>
          </div>
        </div>
      </div>

      <!-- Cast button, Playlist button, and Volume controls in a row -->
      <div class="flex flex-row items-center gap-x-3 md:gap-x-4">
        <!-- Cast button -->
        <div class="inline-block size-22px flex cursor-pointer items-center sm:size-24px">
          <google-cast-launcher />
        </div>

        <!-- Playlist button -->
        <div>
          <RouterLink
            to="/queue"
            class="block flex gap-x-1 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200 sm:gap-x-2 sm:px-3 sm:py-2"
          >
            <icon-tabler-playlist class="text-xl sm:text-xl" />
          </RouterLink>
        </div>

        <!-- Volume controls -->
        <div v-if="audioRef" id="volume-range-input" class="flex flex-row cursor-pointer items-center gap-2 md:gap-2">
          <div @click="toggleMute()">
            <icon-tabler-volume v-if="audioRef.volume > 0.5" class="text-sm sm:text-sm" />
            <icon-tabler-volume-2 v-else-if="audioRef.volume > 0" class="text-sm sm:text-sm" />
            <icon-tabler-volume-3 v-else class="text-sm sm:text-sm" />
          </div>
          <input
            type="range"
            class="h-2 w-20 cursor-pointer bg-white/60 accent-zene-200 md:w-30 sm:w-24"
            max="1"
            step="0.01"
            :value="currentVolume"
            @input="volumeInput"
          />
        </div>
      </div>
    </div>
  </footer>
</template>

<style scoped>
@keyframes pulse-bg {
  0% {
    background-position: top;
  }
  50% {
    background-position: bottom;
  }
  100% {
    background-position: top;
  }
}

.animate-pulse-bg {
  animation: pulse-bg 120s infinite ease-in-out;
}

:global(google-cast-launcher .cast_caf_state_c) {
  fill: rgb(101, 199, 117);
  opacity: 100;
}
:global(google-cast-launcher .cast_caf_state_d) {
  fill: rgb(250, 250, 250);
  opacity: 100;
}
:global(google-cast-launcher .cast_caf_state_h) {
  fill: rgb(250, 250, 250);
  opacity: 50;
}
</style>
