<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { TrackMetadata } from '../types'

const props = defineProps({
  track: { type: Object as PropType<TrackMetadata>, required: false },
})

const trackUrl = computed<string>(() => {
  return `api/files/${props.track?.file_id}/download`
})

const audioRef = ref<HTMLAudioElement | null>(null)
const isPlaying = ref(false)

const imageUrl = computed<string>(() => {
  return `api/albums/${props.track?.musicbrainz_album_id}/art`
})

const togglePlayback = () => {
  if (!audioRef.value) return

  if (isPlaying.value) {
    audioRef.value.pause()
  } else {
    audioRef.value.play()
  }
}

const updateIsPlaying = () => {
  if (!audioRef.value) return
  isPlaying.value = !audioRef.value.paused
}

function onImageError(event: Event) {
  const target = event.target as HTMLImageElement
  target.onerror = null
  target.src = '/default-square.png'
}

onMounted(() => {
  const audio = audioRef.value
  if (!audio) return

  audio.addEventListener('play', updateIsPlaying)
  audio.addEventListener('pause', updateIsPlaying)
  audio.addEventListener('ended', () => {
    isPlaying.value = false
  })
})

onUnmounted(() => {
  const audio = audioRef.value
  if (!audio) return

  audio.pause()
  audio.removeAttribute('src')
  audio.load()
})
</script>

<template>
  <div class="m-4 rounded-xl bg-zene-800/60 p-4">
    <audio v-if="track" ref="audioRef" :src="trackUrl" preload="metadata" />
    <div class="flex flex-col items-center gap-2 text-center">
      <img :src="imageUrl" class="size-50 rounded-lg object-cover" @error="onImageError">
      <div class="font-semibold">
        {{ track?.title }}
      </div>
      <div class="text-sm text-gray-400">
        {{ track?.artist }} · {{ track?.album }}
      </div>
    </div>
    <div class="mt-4 flex justify-between text-xs">
      <span>2:34</span><span>3:21</span>
    </div>
    <div class="my-1 h-1 rounded bg-gray-600"></div>
    <div class="mt-2 flex justify-center space-x-4">
      <button>
        <icon-tabler-player-skip-back-filled />
      </button>
      <button @click="togglePlayback">
        <icon-tabler-player-play-filled class="text-3xl" />
        {{ isPlaying ? 'Pause' : 'Play' }}
      </button>
      <button>
        <icon-tabler-player-skip-forward-filled />
      </button>
    </div>
    <button class="mt-2 text-xs underline">
      LYRICS
    </button>
  </div>
</template>
