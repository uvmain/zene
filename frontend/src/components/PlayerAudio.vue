<script setup lang="ts">
import { audioElement, createContextOnPlay } from '~/logic/audioElement'
import { handleNextTrack, isPlaying, trackUrl, updateProgress } from '~/logic/playbackQueue'

const audioRef = useTemplateRef('audioRef')

onMounted(() => {
  const audio = audioElement.value = audioRef.value
  if (!audio) {
    return
  }
  audio.addEventListener('play', () => createContextOnPlay())
  audio.addEventListener('play', () => isPlaying.value = true)
  audio.addEventListener('pause', () => isPlaying.value = false)
  audio.addEventListener('timeupdate', () => updateProgress())
  audio.addEventListener('ended', () => handleNextTrack())
})

onUnmounted(() => {
  const audio = audioRef.value
  if (audio) {
    audio.replaceWith(audio.cloneNode(true))
    audio.pause()
    audio.removeAttribute('src')
    audio.load()
  }
})
</script>

<template>
  <audio ref="audioRef" :src="trackUrl" preload="metadata" />
</template>
