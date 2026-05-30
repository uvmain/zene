<script setup lang="ts">
import { audioElement, createContextOnPlay } from '~/logic/audioElement'
import { handleNextTrack, isPlaying, updateProgress } from '~/logic/playbackQueue'

const audioRef1 = useTemplateRef('audioRef1') as Ref<HTMLAudioElement>

function isActiveAudioTarget(event: Event) {
  return event.target === audioElement.value
}

function onPlay(event: Event) {
  if (!isActiveAudioTarget(event))
    return
  isPlaying.value = true
}

function onPause(event: Event) {
  if (!isActiveAudioTarget(event))
    return
  isPlaying.value = false
}

function onTimeUpdate(event: Event) {
  if (!isActiveAudioTarget(event))
    return
  updateProgress()
}

function onEnded(event: Event) {
  if (!isActiveAudioTarget(event))
    return
  void handleNextTrack()
}

function registerAudioEvents(audio: HTMLAudioElement) {
  audio.addEventListener('play', createContextOnPlay)
  audio.addEventListener('play', onPlay)
  audio.addEventListener('pause', onPause)
  audio.addEventListener('timeupdate', onTimeUpdate)
  audio.addEventListener('ended', onEnded)
}

function removeAudioEvents(audio: HTMLAudioElement) {
  audio.removeEventListener('play', createContextOnPlay)
  audio.removeEventListener('play', onPlay)
  audio.removeEventListener('pause', onPause)
  audio.removeEventListener('timeupdate', onTimeUpdate)
  audio.removeEventListener('ended', onEnded)
}

onMounted(() => {
  audioElement.value = audioRef1.value
  registerAudioEvents(audioRef1.value)
})

onUnmounted(() => {
  audioElement.value = null
  removeAudioEvents(audioRef1.value)
})
</script>

<template>
  <div class="hidden">
    <audio ref="audioRef1" preload="auto" crossorigin="anonymous" />
  </div>
</template>
