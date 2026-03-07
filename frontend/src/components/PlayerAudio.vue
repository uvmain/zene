<script setup lang="ts">
import { debugLog } from '~/logic/logger'
import { audioContext, audioElement, audioNode, currentlyPlayingTrack, handleNextTrack, isPlaying, playcountPosted, trackUrl, updateProgress } from '~/logic/playbackQueue'

const audioRef = useTemplateRef('audioRef')
const contextCreated = ref(false)

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

function getAudioContext() {
  if (typeof window !== 'undefined') {
    const AudioCtx = (window as any).AudioContext || (window as any).webkitAudioContext
    return new AudioCtx()
  }
  return null
}

// One-time play event to create AudioContext after user interaction
function createContextOnPlay() {
  const audio = audioRef.value
  if (!audio) {
    return
  }
  if (!contextCreated.value) {
    audioContext.value = getAudioContext()
    if (audioContext.value) {
      audioNode.value = audioContext.value.createMediaElementSource(audio)
      audioNode.value.connect(audioContext.value.destination)
      contextCreated.value = true
      debugLog('Audio context created')
    }
    else {
      debugLog('Failed to create audio context')
    }
  }
  audio.removeEventListener('play', createContextOnPlay)
}

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
  track url:{{ trackUrl }}
  <audio ref="audioRef" :src="trackUrl" preload="metadata" />
</template>
