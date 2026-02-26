<script setup lang="ts">
import { getAuthenticatedTrackUrl } from '~/logic/common'
import { debugLog } from '~/logic/logger'
import { audioContext, audioElement, audioNode, currentlyPlayingPodcastEpisode, currentlyPlayingTrack, handleNextTrack, isPlaying, playcountPosted, trackUrl, updateProgress } from '~/logic/playbackQueue'
import { episodeIsStored, getStoredEpisode } from '~/stores/usePodcastStore'

const audioRef = useTemplateRef('audioRef')
const contextCreated = ref(false)

watch(currentlyPlayingTrack, (newTrack, oldTrack) => {
  if (newTrack !== oldTrack) {
    trackUrl.value = newTrack ? getAuthenticatedTrackUrl(newTrack.musicBrainzId) : ''
    playcountPosted.value = false
  }
  else {
    trackUrl.value = newTrack ? getAuthenticatedTrackUrl(newTrack.musicBrainzId) : ''
  }
})

watch(currentlyPlayingPodcastEpisode, (newEpisode, oldEpisode) => {
  if (newEpisode !== oldEpisode) {
    if (newEpisode) {
      episodeIsStored(newEpisode.streamId).then((stored) => {
        if (!stored) {
          trackUrl.value = getAuthenticatedTrackUrl(newEpisode.streamId, true)
        }
        else {
          getStoredEpisode(newEpisode.streamId).then((blob) => {
            if (blob) {
              const objectUrl = URL.createObjectURL(blob)
              trackUrl.value = objectUrl
            }
          })
        }
      })
    }
    else {
      trackUrl.value = ''
    }
  }
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

onMounted(() => {
  audioElement.value = audioRef.value

  const audio = audioRef.value
  if (!audio) {
    return
  }
  // One-time play event to create AudioContext after user interaction
  const createContextOnPlay = () => {
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
