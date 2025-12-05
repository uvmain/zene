<script setup lang="ts">
import { getAuthenticatedTrackUrl } from '~/composables/logic'
import { useDebug } from '~/composables/useDebug'
import { usePlaybackQueue } from '~/composables/usePlaybackQueue'
import { usePlaycounts } from '~/composables/usePlaycounts'
import { useRouteTracks } from '~/composables/useRouteTracks'

const { debugLog } = useDebug()
const { clearQueue, currentlyPlayingTrack, currentlyPlayingPodcastEpisode, resetCurrentlyPlayingTrack, getNextTrack, getPreviousTrack, getRandomTracks, currentQueue, setCurrentQueue } = usePlaybackQueue()
const { routeTracks } = useRouteTracks()
const { postPlaycount } = usePlaycounts()
const router = useRouter()

const audioPlayer = useTemplateRef('playerAudio')
const isPlaying = ref(false)
const playcountPosted = ref(false)
const currentTime = ref(0)
const previousVolume = ref(1)
const currentVolume = ref(1)

const trackUrl = computed(() => {
  if (currentlyPlayingTrack.value) {
    return getAuthenticatedTrackUrl(currentlyPlayingTrack.value?.musicBrainzId)
  }
  else if (currentlyPlayingPodcastEpisode.value) {
    return getAuthenticatedTrackUrl(currentlyPlayingPodcastEpisode.value?.streamId, true)
  }
  return ''
})

async function togglePlayback() {
  if (!audioPlayer.value?.audioRef) {
    console.error('Audio element not found')
    return
  }
  else if (!currentQueue.value?.tracks?.length && routeTracks.value.length) {
    setCurrentQueue(routeTracks.value)
  }
  else if (currentQueue.value?.tracks?.length && !currentlyPlayingTrack.value) {
    setCurrentQueue(currentQueue.value?.tracks)
  }

  if (isPlaying.value) {
    audioPlayer.value?.audioRef.pause()
    isPlaying.value = false
  }
  else {
    audioPlayer.value?.audioRef.play()
    isPlaying.value = true
  }
}

function toggleMute() {
  if (audioPlayer.value?.audioRef) {
    debugLog('Changing volume')
    if (audioPlayer.value.audioRef.volume !== 0) {
      previousVolume.value = audioPlayer.value.audioRef.volume
      audioPlayer.value.audioRef.volume = 0
      currentVolume.value = 0
    }
    else {
      audioPlayer.value.audioRef.volume = previousVolume.value
      currentVolume.value = previousVolume.value
    }
  }
}

async function stopPlayback() {
  if (audioPlayer.value?.audioRef) {
    if (audioPlayer.value.audioRef.currentTime < 1) {
      resetCurrentlyPlayingTrack()
      clearQueue()
    }
    audioPlayer.value.audioRef.pause()
    audioPlayer.value.audioRef.load()
  }

  isPlaying.value = false
}

function updateProgress() {
  if (!audioPlayer.value?.audioRef) {
    return
  }
  currentTime.value = audioPlayer.value.audioRef.currentTime

  if (currentlyPlayingTrack.value && !playcountPosted.value) {
    const halfwayPoint = currentlyPlayingTrack.value.duration / 2
    if (currentTime.value >= halfwayPoint) {
      postPlaycount(currentlyPlayingTrack.value.musicBrainzId)
      playcountPosted.value = true
    }
  }
  else if (currentlyPlayingPodcastEpisode.value && !playcountPosted.value) {
    const halfwayPoint = Number(currentlyPlayingPodcastEpisode.value.duration) / 2
    if (currentTime.value >= halfwayPoint) {
      postPlaycount(currentlyPlayingPodcastEpisode.value.streamId)
      playcountPosted.value = true
    }
  }
}

watch(currentlyPlayingTrack, (newTrack, oldTrack) => {
  if (newTrack !== oldTrack) {
    playcountPosted.value = false
  }
})

function seek(seekSeconds: number) {
  if (audioPlayer.value?.audioRef) {
    audioPlayer.value.audioRef.currentTime = seekSeconds
  }
}

function volumeInput(volumeString: string) {
  const volume = Number.parseFloat(volumeString)

  if (audioPlayer.value?.audioRef) {
    audioPlayer.value.audioRef.volume = volume
  }

  currentVolume.value = volume
}

async function handleNextTrack() {
  if (currentQueue.value && currentQueue.value.tracks.length > 0) {
    await getNextTrack()
  }
  else {
    isPlaying.value = false
  }
}

async function handlePreviousTrack() {
  await getPreviousTrack()
}

async function handleGetRandomTracks() {
  await getRandomTracks(500)
  router.push('/queue')
}
</script>

<template>
  <footer
    class="sticky bottom-0 mt-auto w-full border-0 border-t-1 border-zshade-400 border-zshade-600 border-solid background-2"
  >
    <div class="flex flex-col items-center px-2 md:flex-row space-y-2 md:px-4 md:space-x-2 md:space-y-0">
      <div
        class="h-full w-full flex flex-grow flex-col items-center justify-center py-2 space-y-2 md:py-2 md:space-y-2"
      >
        <PlayerAudio
          ref="playerAudio"
          :track-url="trackUrl"
          @play="() => { isPlaying = true }"
          @pause="() => { isPlaying = false }"
          @time-update="updateProgress"
          @ended="handleNextTrack()"
        />
        <div>
          <PlayerProgressBar
            :current-time-in-seconds="currentTime"
            :currently-playing-track="currentlyPlayingTrack"
            :currently-playing-podcast-episode="currentlyPlayingPodcastEpisode"
            @seek="seek"
          />

          <PlayerMediaControls
            :is-playing="isPlaying"
            @stop-playback="stopPlayback()"
            @toggle-playback="togglePlayback()"
            @next-track="handleNextTrack()"
            @previous-track="handlePreviousTrack()"
            @get-random-tracks="handleGetRandomTracks()"
          />
        </div>
      </div>

      <div class="flex flex-row items-center gap-x-3 md:gap-x-4">
        <!-- <PlayerCastButton /> -->
        <PlayerLyricsButton
          :current-time="currentTime"
          :currently-playing-track="currentlyPlayingTrack"
        />
        <PlayerPlaylistButton />
        <PlayerVolumeSlider
          :audio-ref="audioPlayer?.audioRef"
          :model-value="currentVolume"
          @toggle-mute="toggleMute()"
          @update:model-value="volumeInput"
        />
      </div>
    </div>
  </footer>
</template>
