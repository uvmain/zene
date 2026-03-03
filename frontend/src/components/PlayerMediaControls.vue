<script setup lang="ts">
import { onKeyStroke } from '@vueuse/core'
import { getRandomTracks, handleNextTrack, handlePreviousTrack, isPlaying, stopPlayback, togglePlayback, toggleRepeat, toggleShuffle } from '~/logic/playbackQueue'
import { repeatStatus, shuffleEnabled } from '~/logic/store'

defineProps({
  compact: { type: Boolean, default: false },
})

const router = useRouter()
const route = useRoute()

async function handleGetRandomTracks() {
  await getRandomTracks(500)
  if (route.path !== '/queue' && route.path !== '/visualizer') {
    router.push('/queue')
  }
}

const repeatAbbreviation = computed(() => {
  switch (repeatStatus.value) {
    case 'off':
      return 'repeat 1'
    case '1':
      return 'repeat all'
    case 'all':
      return 'repeat off'
    default:
      return 'repeat off'
  }
})

onKeyStroke('MediaPlayPause', (e) => {
  e.preventDefault()
  togglePlayback()
})

onKeyStroke('MediaTrackPrevious', (e) => {
  e.preventDefault()
  handlePreviousTrack()
})

onKeyStroke('MediaTrackNext', (e) => {
  e.preventDefault()
  handleNextTrack()
})

onKeyStroke('MediaStop', (e) => {
  e.preventDefault()
  stopPlayback()
})
</script>

<template>
  <div
    class="flex flex-row items-center justify-center"
    :class="{
      'gap-x-2 lg:gap-x-4': !compact,
      'gap-x-1': compact,
    }"
  >
    <button id="repeat" title="Stop playback" class="media-control-button" @click="stopPlayback()">
      <icon-nrk-media-stop class="footer-icon" />
    </button>
    <button id="shuffle" title="Shuffle" class="media-control-button" @click="toggleShuffle()">
      <icon-ion-shuffle-sharp
        :class="{
          'footer-icon': !shuffleEnabled,
          'footer-icon-on': shuffleEnabled,
        }"
      />
    </button>
    <button id="back" title="Previous track" class="media-control-button" @click="handlePreviousTrack()">
      <icon-nrk-media-previous class="footer-icon" />
    </button>
    <ZButton
      v-if="!compact"
      id="play-pause"
      class="group/button"
      :primary="true"
      :size12="true"
      hover-text="Play/Pause"
      @click="togglePlayback()"
    >
      <icon-nrk-media-play v-if="!isPlaying" class="footer-icon" />
      <icon-nrk-media-pause v-else class="footer-icon" />
    </ZButton>
    <button v-else id="play-pause-compact" title="Play/Pause" class="media-control-button" @click="togglePlayback()">
      <icon-nrk-media-play v-if="!isPlaying" class="footer-icon" />
      <icon-nrk-media-pause v-else class="footer-icon" />
    </button>
    <button id="forward" title="Next track" class="media-control-button" @click="handleNextTrack()">
      <icon-nrk-media-next class="footer-icon" />
    </button>
    <button id="repeat" :title="repeatAbbreviation" class="relative media-control-button" @click="toggleRepeat">
      <icon-nrk-media-jumpto
        :class="{
          'footer-icon': repeatStatus === 'off',
          'footer-icon-on': repeatStatus !== 'off',
        }"
      />
      <span
        v-if="repeatStatus !== 'off'"
        class="absolute top-0 w-4 text-left text-xs text-primary2 -right-1"
      >
        {{ repeatStatus }}
      </span>
    </button>
    <button
      id="shuffle"
      title="Play random tracks"
      class="media-control-button"
      @click="handleGetRandomTracks()"
    >
      <icon-nrk-dice-3-active class="footer-icon" />
    </button>
  </div>
</template>
