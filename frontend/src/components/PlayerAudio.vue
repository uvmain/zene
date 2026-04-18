<script setup lang="ts">
import { audioElement, createContextOnPlay } from '~/logic/audioElement'
import { handleNextTrack, isPlaying, trackUrl, updateProgress } from '~/logic/playbackQueue'

const audioRef = useTemplateRef('audioRef')

function onPlay() {
  isPlaying.value = true
}

function onPause() {
  isPlaying.value = false
}

function onTimeUpdate() {
  updateProgress()
}

function onEnded() {
  handleNextTrack()
}

onMounted(() => {
  if (!audioRef.value)
    return
  audioElement.value = audioRef.value
  const audio = audioRef.value
  audio.addEventListener('play', createContextOnPlay)
  audio.addEventListener('play', onPlay)
  audio.addEventListener('pause', onPause)
  audio.addEventListener('timeupdate', onTimeUpdate)
  audio.addEventListener('ended', onEnded)
})

onUnmounted(() => {
  if (!audioRef.value)
    return
  const audio = audioRef.value
  audio.removeEventListener('play', createContextOnPlay)
  audio.removeEventListener('play', onPlay)
  audio.removeEventListener('pause', onPause)
  audio.removeEventListener('timeupdate', onTimeUpdate)
  audio.removeEventListener('ended', onEnded)
})
</script>

<template>
  <audio ref="audioRef" :src="trackUrl" preload="metadata" crossorigin="anonymous" />
</template>
