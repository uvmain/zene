<script setup lang="ts">
import { handleNextTrack, isPlaying, updateProgress } from '~/logic/playbackQueue'

const props = defineProps({
  trackUrl: { type: String, required: false, default: '' },
})

const audioRef = useTemplateRef('audioRef')

defineExpose({ audioRef })

watch(
  () => props.trackUrl,
  (newTrack, oldTrack) => {
    const audio = audioRef.value
    if (!audio || newTrack === oldTrack) {
      return
    }
    if (!newTrack) {
      audio.pause()
      audio.removeAttribute('src')
      audio.load()
      audio.currentTime = 0
      return
    }
    if (audio) {
      audio.addEventListener(
        'canplaythrough',
        () => {
          audio?.play()
        },
        { once: true },
      )
      audio.pause()
      audio.load()
    }
  },
)

onMounted(() => {
  const audio = audioRef.value
  if (!audio) {
    return
  }
  audio.addEventListener('play', () => isPlaying.value = true)
  audio.addEventListener('pause', () => isPlaying.value = false)
  audio.addEventListener('timeupdate', () => updateProgress())
  audio.addEventListener('ended', () => handleNextTrack())
})

onUnmounted(() => {
  const audio = audioRef.value
  if (audio) {
    audio.removeEventListener('play', () => isPlaying.value = true)
    audio.removeEventListener('pause', () => isPlaying.value = false)
    audio.removeEventListener('timeupdate', () => updateProgress())
    audio.removeEventListener('ended', () => handleNextTrack())

    audio.pause()
    audio.removeAttribute('src')
    audio.load()
  }
})
</script>

<template>
  <audio ref="audioRef" :src="trackUrl" preload="metadata" class="hidden" />
</template>
