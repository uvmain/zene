<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { currentlyPlayingTrack, resetCurrentlyPlayingTrack, setCurrentlyPlayingTrack } from '../composables/globalState'
import { formatTime } from '../composables/logic'
import { getRandomTrack } from '../composables/randomTrack'

const audioRef = ref<HTMLAudioElement | null>(null)
const isPlaying = ref(false)
const currentTime = ref(0)
const currentVolume = ref(1)
const isPlayPauseActive = ref(false)
const debug = ref()

const trackUrl = computed<string>(() => {
  return currentlyPlayingTrack.value?.file_id ? `/api/files/${currentlyPlayingTrack.value.file_id}/stream` : ''
})

function togglePlayback() {
  if (!audioRef.value) {
    return
  }
  if (isPlaying.value) {
    audioRef.value.pause()
  }
  else {
    audioRef.value.play()
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
  if (!audioRef.value)
    return
  if (audioRef.value.currentTime < 1) {
    resetCurrentlyPlayingTrack()
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
    debug.value = 'no audioref'
    return
  }
  debug.value = `updating currentTime to ${currentTime.value}`
  currentTime.value = audioRef.value.currentTime
}

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

async function randomTrack() {
  if (!audioRef.value)
    return

  audioRef.value.pause()
  audioRef.value.removeAttribute('src')
  currentTime.value = 0
  const randomTrack = await getRandomTrack()
  setCurrentlyPlayingTrack(randomTrack)
  audioRef.value.addEventListener(
    'canplay',
    () => {
      audioRef.value?.play()
    },
    { once: true },
  )
  audioRef.value.src = trackUrl.value
  audioRef.value.load()
}

onBeforeMount(() => {
  randomTrack()
})

watch(currentlyPlayingTrack, (newTrack, oldTrack) => {
  const audio = audioRef.value
  if (!audio) {
    return
  }
  if (newTrack && newTrack.file_id !== oldTrack?.file_id) {
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
    isPlaying.value = false
  })
})

async function playNext() {
  const audio = audioRef.value
  if (!audio)
    return

  audio.pause()
  audio.removeAttribute('src')
  currentTime.value = 0
  const randomTrack = await getRandomTrack()
  setCurrentlyPlayingTrack(randomTrack)
  audio.addEventListener(
    'canplay',
    () => {
      audio.play()
    },
    { once: true },
  )
  audio.src = trackUrl.value
  audio.load()
}

onUnmounted(() => {
  const audio = audioRef.value
  if (!audio)
    return

  audio.removeEventListener('play', updateIsPlaying)
  audio.removeEventListener('pause', updateIsPlaying)
  audio.removeEventListener('timeupdate', updateProgress)
  audio.removeEventListener('ended', () => {
    isPlaying.value = false
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
            <button id="back" class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none" @click="togglePlayback()">
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
            <button id="forward" class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none" @click="playNext()">
              <icon-tabler-player-skip-forward class="text-xl" />
            </button>
            <button id="repeat" class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none" @click="togglePlayback()">
              <icon-tabler-repeat class="text-xl" />
            </button>
            <button
              id="shuffle"
              class="h-12 w-12 flex cursor-pointer items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none"
              @click="randomTrack()"
            >
              <icon-tabler-dice-3 class="text-xl" />
            </button>
          </div>
        </div>
      </div>
      <div id="volume-range-input" class="flex flex-row items-center gap-2">
        <icon-tabler-volume v-if="currentVolume > 0.5" class="text-sm" />
        <icon-tabler-volume-2 v-else-if="currentVolume > 0" class="text-sm" />
        <icon-tabler-volume-3 v-else class="text-sm" />
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
</style>
