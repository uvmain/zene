<script setup lang="ts">
import { audioElement, createContextOnPlay, currentlyPlayingTrack, handleNextTrack, isPlaying, playcountPosted, trackUrl, updateProgress } from '~/logic/playbackQueue'

const audioRef = useTemplateRef('audioRef')

watch(currentlyPlayingTrack, () => {
  playcountPosted.value = false
})

watch(trackUrl, (newTrack, oldTrack) => {
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
})

onMounted(() => {
  const audio = audioElement.value = audioRef.value
  if (!audio) {
    return
  }
  audio.addEventListener('play', createContextOnPlay)
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
