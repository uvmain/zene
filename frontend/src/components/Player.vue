<script setup lang="ts">
import type { TrackMetadata } from '../types'
import { onMounted, onUnmounted, ref } from 'vue'

const props = defineProps({
  track: { type: Object as PropType<TrackMetadata>, required: false },
})

const trackUrl = computed<string>(() => {
  return props.track ? `api/files/${props.track?.file_id}/download` : undefined
})

const audioRef = ref<HTMLAudioElement | null>(null)
const isPlaying = ref(false)
const currentTime = ref(0)

const imageUrl = computed<string>(() => {
  return props.track ? `api/albums/${props.track.musicbrainz_album_id}/art` : 'default-square.png'
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
  updateIsPlaying()
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

function formatTime(time: number): string {
  const minutes = Math.floor(time / 60)
  const seconds = Math.floor(time % 60)
  return `${minutes}:${seconds.toString().padStart(2, '0')}`
}

function seek(event: Event) {
  if (!audioRef.value)
    return
  const target = event.target as HTMLInputElement
  const seekTime = Number.parseFloat(target.value)
  audioRef.value.currentTime = seekTime
}

function onImageError(event: Event) {
  const target = event.target as HTMLImageElement
  target.onerror = null
  target.src = '/default-square.png'
}

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
  <div class="m-4 rounded-xl bg-zene-800/60 p-4">
    <audio ref="audioRef" :src="trackUrl" preload="metadata" />

    <!-- album art -->
    <div class="flex flex-col items-center gap-2 text-center">
      <img :src="imageUrl" class="size-50 rounded-lg object-cover" @error="onImageError">
      <div class="font-semibold">
        {{ track?.title }}
      </div>
      <div class="text-sm text-gray-400">
        {{ track?.artist }} · {{ track?.album }} · {{ track?.track_number }} / {{ track?.total_tracks }}
      </div>
    </div>

    <!-- Progress Bar -->
    <div v-if="audioRef" class="mt-4 flex flex-col items-center">
      <input
        type="range"
        class="h-1 w-full cursor-pointer appearance-none rounded-lg bg-gray-600"
        :max="track ? track.duration : 0"
        :value="currentTime"
        @input="seek"
      />
      <!-- Time Display -->
      <div class="mt-2 w-full flex justify-between text-xs text-gray-400">
        <span>{{ formatTime(currentTime) }}</span>
        <span>{{ formatTime(track ? track.duration : 0) }}</span>
      </div>
    </div>
    <!-- buttons -->
    <div class="mt-2 flex justify-center space-x-4">
      <button id="repeat" class="h-12 w-12 flex items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none" @click="togglePlayback()">
        <icon-tabler-repeat class="text-xl" />
      </button>
      <button id="back" class="h-12 w-12 flex items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none" @click="togglePlayback()">
        <icon-tabler-player-skip-back class="text-xl" />
      </button>
      <button id="play-pause" class="h-12 w-12 flex items-center justify-center rounded-md border-none bg-zene-400 text-white font-semibold outline-none" @click="togglePlayback()">
        <icon-tabler-player-play v-if="!isPlaying" class="text-3xl" />
        <icon-tabler-player-pause v-else class="text-3xl" />
      </button>
      <button id="forward" class="h-12 w-12 flex items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none" @click="togglePlayback()">
        <icon-tabler-player-skip-forward class="text-xl" />
      </button>
      <button id="shuffle" class="h-12 w-12 flex items-center justify-center rounded-full border-none bg-zene-400/0 text-white font-semibold outline-none" @click="togglePlayback()">
        <icon-tabler-arrows-shuffle class="text-xl" />
      </button>
    </div>
    <button class="mt-2 text-xs underline">
      LYRICS
    </button>
  </div>
</template>
