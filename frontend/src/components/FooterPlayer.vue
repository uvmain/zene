<script setup lang="ts">
import { getAuthenticatedTrackUrl } from '~/logic/common'
import { audioElement, currentlyPlayingPodcastEpisode, currentlyPlayingTrack, playcountPosted, trackUrl } from '~/logic/playbackQueue'
import { episodeIsStored, getStoredEpisode } from '~/stores/usePodcastStore'

const audioPlayer = useTemplateRef('audioPlayerElement')

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

watch(audioPlayer, () => {
  if (audioPlayer.value) {
    audioElement.value = audioPlayer.value.audioRef
  }
})
</script>

<template>
  <footer
    class="border-zshade-400 sticky bottom-0 mt-auto w-full border-0 border-t-1 border-zshade-600 border-solid background-2"
  >
    <div class="flex flex-col items-center px-2 lg:flex-row space-y-2 lg:px-4 lg:space-x-2 lg:space-y-0">
      <div
        class="h-full w-full flex flex-grow flex-col items-center justify-center py-2 space-y-2 lg:py-2 lg:space-y-2"
      >
        <PlayerAudio
          ref="audioPlayerElement"
          :track-url="trackUrl"
        />
        <div>
          <PlayerProgressBar />
          <PlayerMediaControls />
        </div>
      </div>
      <div class="flex flex-row items-center gap-x-3 lg:gap-x-4">
        <!-- <PlayerCastButton /> -->
        <PlayerLyricsButton />
        <PlayerPlaylistButton />
        <PlayerVolumeSlider />
      </div>
    </div>
  </footer>
</template>
