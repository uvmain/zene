<script setup lang="ts">
import type { TokenResponse } from '../types/auth'
import { onKeyStroke } from '@vueuse/core'
import { formatTime } from '../composables/logic'
import { useBackendFetch } from '../composables/useBackendFetch'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'
import { usePlaycounts } from '../composables/usePlaycounts'
import { useRouteTracks } from '../composables/useRouteTracks'
import { useSettings } from '../composables/useSettings'

const { getTemporaryToken } = useBackendFetch()

const { clearQueue, currentlyPlayingTrack, resetCurrentlyPlayingTrack, getNextTrack, getPreviousTrack, refreshRandomSeed, getRandomTracks, currentQueue, setCurrentQueue } = usePlaybackQueue()
const { streamQuality } = useSettings()
const { routeTracks } = useRouteTracks()
const { postPlaycount, updatePlaycount } = usePlaycounts()
const router = useRouter()
const route = useRoute()

const audioRef = ref<HTMLAudioElement | null>(null)
const isPlaying = ref(false)
const playcountPosted = ref(false)
const currentTime = ref(0)
const previousVolume = ref(1)
const currentVolume = ref(1)
const isPlayPauseActive = ref(false)
const session = ref<cast.framework.CastSession | null>(null)
const temporaryToken = ref<TokenResponse | null>(null)

const currentRoute = computed(() => {
  return route.path
})

const trackUrl = computed<string>(() => {
  return currentlyPlayingTrack.value?.musicbrainz_track_id ? `/api/tracks/${currentlyPlayingTrack.value.musicbrainz_track_id}/stream?quality=${streamQuality.value}` : ''
})

async function togglePlayback() {
  if (!audioRef.value) {
    return
  }
  if (!currentQueue.value?.tracks?.length && routeTracks.value.length) {
    setCurrentQueue(routeTracks.value)
  }
  else if (currentQueue.value?.tracks?.length && !currentlyPlayingTrack.value) {
    setCurrentQueue(currentQueue.value?.tracks)
  }
  if (isPlaying.value) {
    audioRef.value.pause()
  }
  else {
    const context = cast.framework.CastContext.getInstance()
    session.value = context.getCurrentSession()
    if (session.value) {
      console.log('Casting audio')
      await castAudio()
      return
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

  updateIsPlaying()
}

function toggleMute() {
  if (!audioRef.value) {
    return
  }
  console.log('Changing volume')
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

async function stopPlayback() {
  if (!audioRef.value)
    return
  if (audioRef.value.currentTime < 1) {
    resetCurrentlyPlayingTrack()
    clearQueue()
  }
  audioRef.value.pause()
  audioRef.value.load()
  isPlaying.value = false
}

function updateIsPlaying() {
  if (!audioRef.value)
    return
  isPlaying.value = !audioRef.value.paused
}

function updateProgress() {
  if (!audioRef.value) {
    return
  }

  currentTime.value = audioRef.value.currentTime

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
  if (!audioRef.value)
    return
  const target = event.target as HTMLInputElement
  const seekTime = Number.parseFloat(target.value)
  audioRef.value.currentTime = seekTime
}

function volumeInput(event: Event) {
  if (!audioRef.value)
    return
  const target = event.target as HTMLInputElement
  const volume = Number.parseFloat(target.value)
  audioRef.value.volume = volume
  currentVolume.value = volume
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
  getPreviousTrack()
})

onKeyStroke('MediaTrackNext', (e) => {
  e.preventDefault()
  getNextTrack()
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

async function getMimeType(): Promise<string> {
  let requestUrl = trackUrl.value
  if (trackUrl.value.includes('?')) {
    requestUrl = `${trackUrl.value}&token=${temporaryToken.value?.token}`
  }
  else {
    requestUrl = `${trackUrl.value}?token=${temporaryToken.value?.token}`
  }
  await fetch(requestUrl, { method: 'HEAD' })
    .then((response) => {
      const contentType = response.headers.get('content-type') ?? ''
      console.log(contentType)
      return contentType
    })
    .catch((err) => {
      console.log(`fetch failed: ${err}`)
    })
  return ''
}

async function castAudio() {
  const context = cast.framework.CastContext.getInstance()
  session.value = context.getCurrentSession()

  if (!session.value) {
    return
  }

  if (!trackUrl.value) {
    return
  }

  const contentType = await (getMimeType())

  let requestUrl: string
  if (trackUrl.value.includes('?')) {
    requestUrl = `${trackUrl.value}&token=${temporaryToken.value?.token}`
  }
  else {
    requestUrl = `${trackUrl.value}?token=${temporaryToken.value?.token}`
  }
  // prefix base url to requestUrl
  if (window) {
    const protocol = window.location.protocol
    const host = window.location.host
    requestUrl = `${protocol}://${host}${requestUrl}`
  }
  const mediaInfo = new chrome.cast.media.MediaInfo(requestUrl, contentType)
  const request = new chrome.cast.media.LoadRequest(mediaInfo)

  session.value
    .loadMedia(request)
    .then(() => console.log('Media loaded to cast device'))
    .catch(err => console.error('Error loading media:', err))
}

function initializeCast() {
  const context = cast.framework.CastContext.getInstance()
  context.setOptions({
    receiverApplicationId: chrome.cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID,
    autoJoinPolicy: chrome.cast.AutoJoinPolicy.ORIGIN_SCOPED,
  })
  console.log('CastContext initialized')
}

onMounted(async () => {
  temporaryToken.value = await getTemporaryToken()
  console.log('Waiting for Cast SDK...')

  // Hook for when the SDK becomes available (in case it's not already)
  window.__onGCastApiAvailable = (isAvailable: boolean) => {
    console.log('Cast API available (async):', isAvailable)
    if (isAvailable)
      initializeCast()
    else console.warn('Cast API not available')
  }

  // 🔥 If the SDK already loaded and called __onGCastApiAvailable BEFORE this script ran
  // we have to check manually and initialize right now
  if ((window.cast && window.cast.isAvailable) || (window.chrome?.cast && window.chrome.cast.isAvailable)) {
    console.log('Cast API already available (sync)')
    initializeCast()
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
      getNextTrack()
    }
    else {
      isPlaying.value = false
    }
  })
})

onUnmounted(() => {
  const audio = audioRef.value
  if (!audio)
    return

  audio.removeEventListener('play', updateIsPlaying)
  audio.removeEventListener('pause', updateIsPlaying)
  audio.removeEventListener('timeupdate', updateProgress)
  audio.removeEventListener('ended', () => {
    if (currentQueue.value && currentQueue.value.tracks.length > 0) {
      getNextTrack()
    }
    else {
      isPlaying.value = false
    }
  })

  audio.pause()
  audio.removeAttribute('src')
  audio.load()
})
</script>

<template>
  <footer
    class="sticky bottom-0 mt-auto w-full bg-zene-700 bg-cover bg-center"
    :class="{ 'animate-pulse-bg': currentlyPlayingTrack && isPlaying }"
    :style="{ backgroundImage: `url(${currentlyPlayingTrack?.image_url})` }"
  >
    <div class="flex flex-row items-center border-0 border-t-1 border-white/20 border-solid px-4 backdrop-blur-2xl backdrop-contrast-30 space-x-2">
      <div
        class="h-full w-full flex flex-grow flex-col items-center justify-center py-2 space-y-2"
      >
        <audio ref="audioRef" :src="trackUrl" preload="metadata" class="hidden" />
        <div class="">
          <!-- Progress Bar -->
          <div v-if="audioRef" class="max-w-200 flex flex-row items-center gap-2">
            <span id="currentTime" class="w-12 text-right text-sm text-gray-2">
              {{ formatTime(currentTime) }}
            </span>
            <input
              type="range"
              class="h-1 w-full cursor-pointer bg-white/60 accent-zene-200"
              :max="currentlyPlayingTrack ? currentlyPlayingTrack.duration : 0"
              :value="currentTime"
              @input="seek"
            />
            <span id="duration" class="w-12 text-sm text-gray-2">
              {{ formatTime(currentlyPlayingTrack ? Number.parseFloat(currentlyPlayingTrack.duration) : 0) }}
            </span>
          </div>

          <!-- Buttons -->
          <div class="mt-2 flex flex-row items-center justify-center gap-x-4">
            <button id="repeat" class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none" @click="stopPlayback()">
              <icon-tabler-player-stop class="text-xl" />
            </button>
            <button id="shuffle" class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none" @click="togglePlayback()">
              <icon-tabler-arrows-shuffle class="text-xl" />
            </button>
            <button id="back" class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none" @click="getPreviousTrack()">
              <icon-tabler-player-skip-back class="text-xl" />
            </button>
            <button
              id="play-pause"
              class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-md border-none text-white font-semibold outline-none"
              :class="isPlayPauseActive ? 'bg-zene-200' : 'bg-zene-400 transition-colors duration-200'"
              @click="togglePlayback()"
            >
              <icon-tabler-player-play v-if="!isPlaying" class="text-3xl" />
              <icon-tabler-player-pause v-else class="text-3xl" />
            </button>
            <button id="forward" class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none" @click="getNextTrack()">
              <icon-tabler-player-skip-forward class="text-xl" />
            </button>
            <button id="repeat" class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none" @click="togglePlayback()">
              <icon-tabler-repeat class="text-xl" />
            </button>
            <button
              id="shuffle"
              class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none"
              @click="handleGetRandomTracks()"
            >
              <icon-tabler-dice-3 class="text-xl" />
            </button>
          </div>
        </div>
      </div>

      <!-- Cast button -->
      <div class="inline-block size-24px flex cursor-pointer items-center">
        <google-cast-launcher />
      </div>

      <div>
        <RouterLink
          to="/queue"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/' }"
        >
          <icon-tabler-playlist />
        </RouterLink>
      </div>
      <div v-if="audioRef" id="volume-range-input" class="flex flex-row cursor-pointer items-center gap-2">
        <div @click="toggleMute()">
          <icon-tabler-volume v-if="audioRef.volume > 0.5" class="text-sm" />
          <icon-tabler-volume-2 v-else-if="audioRef.volume > 0" class="text-sm" />
          <icon-tabler-volume-3 v-else class="text-sm" />
        </div>
        <input
          type="range"
          class="h-1 w-30 cursor-pointer bg-white/60 accent-zene-200"
          max="1"
          step="0.01"
          :value="currentVolume"
          @input="volumeInput"
        />
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
