<script setup lang="ts">
import { onKeyStroke } from '@vueuse/core'
import { formatTime } from '../composables/logic'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'
import { usePlaycounts } from '../composables/usePlaycounts'
import { useRouteTracks } from '../composables/useRouteTracks'
import { useSettings } from '../composables/useSettings'

const { clearQueue, currentlyPlayingTrack, resetCurrentlyPlayingTrack, getNextTrack, getPreviousTrack, refreshRandomSeed, getRandomTracks, currentQueue, setCurrentQueue } = usePlaybackQueue()
const { streamQuality } = useSettings()
const { routeTracks } = useRouteTracks()
const { postPlaycount, updatePlaycount } = usePlaycounts()

const audioRef = ref<HTMLAudioElement | null>(null)
const isPlaying = ref(false)
const playcountPosted = ref(false)
const currentTime = ref(0)
const previousVolume = ref(1)
const currentVolume = ref(1)
const isPlayPauseActive = ref(false)
const isCasting = ref(false) // TODO: Implement actual casting detection

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
  
  if (isCasting.value) {
    // TODO: Implement casting playback control
    console.log('Casting playback toggle not yet implemented')
  } else {
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

  updateIsPlaying()
}

async function stopPlayback() {
  if (isCasting.value) {
    // TODO: Implement casting stop control
    console.log('Casting stop not yet implemented')
  } else {
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
}

async function previousTrack() {
  if (isCasting.value) {
    // TODO: Implement casting previous control
    console.log('Casting previous not yet implemented')
  }
  getPreviousTrack()
}

async function nextTrack() {
  if (isCasting.value) {
    // TODO: Implement casting next control
    console.log('Casting next not yet implemented')
  }
  getNextTrack()
}

function updateIsPlaying() {
  if (isCasting.value) {
    // TODO: Get casting playback state
    console.log('Casting state not yet implemented')
  } else {
    if (!audioRef.value)
      return
    isPlaying.value = !audioRef.value.paused
  }
}

function updateProgress() {
  if (isCasting.value) {
    // TODO: Get casting progress
    console.log('Casting progress not yet implemented')
    return
  }

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
  const target = event.target as HTMLInputElement
  const seekTime = Number.parseFloat(target.value)
  
  if (isCasting.value) {
    // TODO: Implement casting seek
    console.log(`Casting seek to ${seekTime} not yet implemented`)
  } else {
    if (!audioRef.value)
      return
    audioRef.value.currentTime = seekTime
  }
}

function toggleCasting() {
  isCasting.value = !isCasting.value
  console.log(`Casting ${isCasting.value ? 'enabled' : 'disabled'} (placeholder)`)
}

// Keyboard controls
onKeyStroke('MediaPlayPause', (e) => {
  e.preventDefault()
  togglePlayback()
})

onKeyStroke('MediaTrackPrevious', (e) => {
  e.preventDefault()
  previousTrack()
})

onKeyStroke('MediaTrackNext', (e) => {
  e.preventDefault()
  nextTrack()
})

onKeyStroke('MediaStop', (e) => {
  e.preventDefault()
  stopPlayback()
})

watch(currentlyPlayingTrack, (newTrack, oldTrack) => {
  const audio = audioRef.value
  if (!audio || isCasting.value) {
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
  <div class="max-w-4xl mx-auto p-6 space-y-8">
    <h1 class="text-3xl font-bold text-white mb-8">Debug Player Controls</h1>
    
    <!-- Audio element without controls -->
    <audio ref="audioRef" :src="trackUrl" preload="metadata" class="hidden" />
    
    <!-- Current Track Info -->
    <div v-if="currentlyPlayingTrack" class="bg-zene-600 rounded-lg p-4">
      <h2 class="text-xl font-semibold text-white mb-2">Currently Playing</h2>
      <div class="flex items-center space-x-4">
        <img v-if="currentlyPlayingTrack.image_url" :src="currentlyPlayingTrack.image_url" class="w-16 h-16 rounded-md" alt="Album art" />
        <div>
          <p class="text-white font-medium">{{ currentlyPlayingTrack.title }}</p>
          <p class="text-gray-300">{{ currentlyPlayingTrack.artist }}</p>
          <p class="text-gray-400 text-sm">{{ currentlyPlayingTrack.album }}</p>
        </div>
      </div>
    </div>

    <!-- Casting Status -->
    <div class="bg-zene-600 rounded-lg p-4">
      <h2 class="text-xl font-semibold text-white mb-2">Casting Status</h2>
      <div class="flex items-center space-x-4">
        <span class="text-white">Mode: {{ isCasting ? 'Casting' : 'Local Audio' }}</span>
        <button
          class="px-4 py-2 bg-zene-400 hover:bg-zene-300 text-white rounded-md transition-colors"
          @click="toggleCasting()"
        >
          {{ isCasting ? 'Disable Casting' : 'Enable Casting' }}
        </button>
      </div>
    </div>
    
    <!-- Custom Seek Bar -->
    <div class="bg-zene-600 rounded-lg p-4">
      <h2 class="text-xl font-semibold text-white mb-4">Seek Control</h2>
      <div v-if="audioRef || isCasting" class="flex flex-row items-center gap-4">
        <span class="w-16 text-right text-sm text-gray-300">
          {{ formatTime(currentTime) }}
        </span>
        <input
          type="range"
          class="flex-1 h-2 cursor-pointer bg-white/20 accent-zene-200 rounded-lg"
          :max="currentlyPlayingTrack ? currentlyPlayingTrack.duration : 0"
          :value="currentTime"
          @input="seek"
        />
        <span class="w-16 text-sm text-gray-300">
          {{ formatTime(currentlyPlayingTrack ? Number.parseFloat(currentlyPlayingTrack.duration) : 0) }}
        </span>
      </div>
      <div v-else class="text-gray-400">
        No audio loaded
      </div>
    </div>

    <!-- Media Control Buttons -->
    <div class="bg-zene-600 rounded-lg p-6">
      <h2 class="text-xl font-semibold text-white mb-4">Media Controls</h2>
      <div class="flex flex-row items-center justify-center gap-x-6">
        <!-- Previous Button -->
        <button
          class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-md border-none bg-zene-400 hover:bg-zene-300 text-white font-semibold outline-none transition-colors duration-200"
          @click="previousTrack()"
          title="Previous Track"
        >
          <icon-tabler-player-skip-back class="text-2xl" />
        </button>

        <!-- Play/Pause Button -->
        <button
          class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-md border-none text-white font-semibold outline-none"
          :class="isPlayPauseActive ? 'bg-zene-200' : 'bg-zene-400 hover:bg-zene-300 transition-colors duration-200'"
          @click="togglePlayback()"
          title="Play/Pause"
        >
          <icon-tabler-player-play v-if="!isPlaying" class="text-3xl" />
          <icon-tabler-player-pause v-else class="text-3xl" />
        </button>

        <!-- Stop Button -->
        <button
          class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-md border-none bg-zene-400 hover:bg-zene-300 text-white font-semibold outline-none transition-colors duration-200"
          @click="stopPlayback()"
          title="Stop"
        >
          <icon-tabler-player-stop class="text-2xl" />
        </button>

        <!-- Next Button -->
        <button
          class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-md border-none bg-zene-400 hover:bg-zene-300 text-white font-semibold outline-none transition-colors duration-200"
          @click="nextTrack()"
          title="Next Track"
        >
          <icon-tabler-player-skip-forward class="text-2xl" />
        </button>
      </div>
    </div>

    <!-- Debug Information -->
    <div class="bg-zene-600 rounded-lg p-4">
      <h2 class="text-xl font-semibold text-white mb-4">Debug Information</h2>
      <div class="space-y-2 text-sm text-gray-300">
        <p><strong>Is Playing:</strong> {{ isPlaying }}</p>
        <p><strong>Current Time:</strong> {{ currentTime.toFixed(2) }}s</p>
        <p><strong>Duration:</strong> {{ currentlyPlayingTrack ? currentlyPlayingTrack.duration : 'N/A' }}s</p>
        <p><strong>Track URL:</strong> {{ trackUrl || 'N/A' }}</p>
        <p><strong>Queue Length:</strong> {{ currentQueue?.tracks?.length || 0 }}</p>
        <p><strong>Casting Mode:</strong> {{ isCasting ? 'Active' : 'Inactive' }}</p>
        <p><strong>Playcount Posted:</strong> {{ playcountPosted }}</p>
      </div>
    </div>
  </div>
</template>