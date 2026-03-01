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
    <abbr title="Stop playback">
      <button id="repeat" class="media-control-button" @click="stopPlayback()">
        <icon-nrk-media-stop class="footer-icon" />
      </button>
    </abbr>
    <abbr title="Shuffle">
      <button id="shuffle" class="media-control-button" @click="toggleShuffle()">
        <icon-ion-shuffle-sharp
          :class="{
            'footer-icon': !shuffleEnabled,
            'footer-icon-on': shuffleEnabled,
          }"
        />
      </button>
    </abbr>
    <abbr title="Previous track">
      <button id="back" class="media-control-button" @click="handlePreviousTrack()">
        <icon-nrk-media-previous class="footer-icon" />
      </button>
    </abbr>
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
    <abbr v-else title="Play/Pause">
      <button id="play-pause-compact" class="media-control-button" @click="togglePlayback()">
        <icon-nrk-media-play v-if="!isPlaying" class="footer-icon" />
        <icon-nrk-media-pause v-else class="footer-icon" />
      </button>
    </abbr>
    <abbr title="Next track">
      <button id="forward" class="media-control-button" @click="handleNextTrack()">
        <icon-nrk-media-next class="footer-icon" />
      </button>
    </abbr>
    <abbr :title="repeatAbbreviation">
      <button id="repeat" class="relative media-control-button" @click="toggleRepeat">
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
    </abbr>
    <abbr title="Play random tracks">
      <button
        id="shuffle"
        class="media-control-button"
        @click="handleGetRandomTracks()"
      >
        <icon-nrk-dice-3-active class="footer-icon" />
      </button>
    </abbr>
  </div>
</template>
