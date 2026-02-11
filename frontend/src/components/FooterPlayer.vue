<script setup lang="ts">
import { getAuthenticatedTrackUrl } from '~/logic/common'
import { currentlyPlayingPodcastEpisode, currentlyPlayingTrack, getRandomTracks } from '~/logic/playbackQueue'
import { currentTime, currentVolume, handleNextTrack, handlePreviousTrack, isPlaying, playcountPosted, seek, shuffleEnabled, stopPlayback, toggleMute, togglePlayback, trackUrl, updateProgress, volumeInput } from '~/logic/playerActions'
import { episodeIsStored, getStoredEpisode } from '~/stores/usePodcastStore'

const router = useRouter()
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

async function handleGetRandomTracks() {
  await getRandomTracks(500, shuffleEnabled.value)
  router.push('/queue')
}
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
          @play="() => { isPlaying = true }"
          @pause="() => { isPlaying = false }"
          @time-update="updateProgress(audioPlayer?.audioRef)"
          @ended="handleNextTrack(audioPlayer?.audioRef)"
        />
        <div>
          <PlayerProgressBar
            :current-time-in-seconds="currentTime"
            :currently-playing-track="currentlyPlayingTrack"
            :currently-playing-podcast-episode="currentlyPlayingPodcastEpisode"
            @seek="seek(audioPlayer?.audioRef, $event)"
          />
          <PlayerMediaControls
            :is-playing="isPlaying"
            @stop-playback="stopPlayback(audioPlayer?.audioRef)"
            @toggle-playback="togglePlayback(audioPlayer?.audioRef)"
            @next-track="handleNextTrack(audioPlayer?.audioRef)"
            @previous-track="handlePreviousTrack()"
            @get-random-tracks="handleGetRandomTracks()"
          />
        </div>
      </div>

      <div class="flex flex-row items-center gap-x-3 lg:gap-x-4">
        <!-- <PlayerCastButton /> -->
        <PlayerLyricsButton
          :current-time="currentTime"
          :currently-playing-track="currentlyPlayingTrack"
        />
        <PlayerPlaylistButton />
        <PlayerVolumeSlider
          :audio-ref="audioPlayer?.audioRef"
          :model-value="currentVolume"
          @toggle-mute="toggleMute(audioPlayer?.audioRef)"
          @update:model-value="volumeInput(audioPlayer?.audioRef, $event)"
        />
      </div>
    </div>
  </footer>
</template>
