<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { currentlyPlayingTrack, resetCurrentlyPlayingTrack } from '../composables/globalState'
import { formatTime } from '../composables/logic'
import { getRandomTrack } from '../composables/randomTrack'

const audioRef = ref<HTMLAudioElement | null>(null)
const isPlaying = ref(false)
const currentTime = ref(0)
const isPlayPauseActive = ref(false)

const trackUrl = computed<string>(() => {
  return currentlyPlayingTrack.value ? `api/files/${currentlyPlayingTrack.value.file_id}/download` : ''
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
  if (!audioRef.value)
    return
  currentTime.value = audioRef.value.currentTime
}

function seek(event: Event) {
  if (!audioRef.value)
    return
  const target = event.target as HTMLInputElement
  const seekTime = Number.parseFloat(target.value)
  audioRef.value.currentTime = seekTime
}

async function randomTrack() {
  if (!audioRef.value)
    return

  audioRef.value.pause()
  audioRef.value.removeAttribute('src')
  currentTime.value = 0
  currentlyPlayingTrack.value = await getRandomTrack()
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

onMounted(() => {
  const audio = audioRef.value
  if (!audio)
    return

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
  currentlyPlayingTrack.value = await getRandomTrack()
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
  <footer class="sticky bottom-0 h-40 w-full border-0 border-t-1 border-white/20 border-solid bg-zene-700">
    <div class="w-full flex flex-grow justify-center">
      <div class="m-2 rounded-xl p-4">
        <audio ref="audioRef" :src="trackUrl" preload="metadata" />

        <!-- Progress Bar -->
        <div v-if="audioRef" class="mt-4 flex flex-col items-center">
          <input
            type="range"
            class="h-1 w-full cursor-pointer appearance-none rounded-lg bg-gray-600 accent-zene-200"
            :max="currentlyPlayingTrack ? currentlyPlayingTrack.duration : 0"
            :value="currentTime"
            @input="seek"
          />
          <!-- Time Display -->
          <div class="mt-2 w-full flex justify-between text-xs text-gray-400">
            <span>{{ formatTime(currentTime) }}</span>
            <span>{{ formatTime(currentlyPlayingTrack ? Number.parseFloat(currentlyPlayingTrack.duration) : 0) }}</span>
          </div>
        </div>
        <!-- Buttons -->
        <div class="mt-2 flex justify-center space-x-4">
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
  </footer>
</template>
