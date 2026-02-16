<script setup lang="ts">
import { getAuthenticatedTrackUrl } from '~/logic/common'
import { audioElement, currentlyPlayingPodcastEpisode, currentlyPlayingTrack, handleNextTrack, isPlaying, playcountPosted, trackUrl, updateProgress } from '~/logic/playbackQueue'
import { episodeIsStored, getStoredEpisode } from '~/stores/usePodcastStore'

const audioRef = useTemplateRef('audioRef')

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

watch(
  trackUrl,
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

watch(audioRef, () => {
  if (audioRef.value) {
    audioElement.value = audioRef.value
  }
})

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
  <audio ref="audioRef" :src="trackUrl" preload="metadata" />
</template>
